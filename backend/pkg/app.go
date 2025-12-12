package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/data/migrations"
	"github.com/dimqueue/darts/pkg/data/seeds"
	"github.com/dimqueue/darts/pkg/handler"
	"github.com/dimqueue/darts/pkg/logger"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/dimqueue/darts/pkg/swagger"
	"github.com/dimqueue/darts/pkg/validation"
	"github.com/dimqueue/darts/pkg/worker"
	"github.com/jmoiron/sqlx"
)

type command string

const (
	cmdRunServer   command = "run-server"
	cmdMigrateUp   command = "migrates-up"
	cmdMigrateDown command = "migrates-down"
	cmdSeed        command = "seed"
)

func ParseCommand() (command, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("usage: %s <command>\n\nAvailable commands:\n  run-server    - Start the HTTP server\n  migrates-up   - Run database migrations\n  migrates-down - Rollback migrations\n  seed          - Load seed data (dev only)", os.Args[0])
	}

	cmd := command(os.Args[1])

	switch cmd {
	case cmdRunServer, cmdMigrateUp, cmdMigrateDown, cmdSeed:
		return cmd, nil
	default:
		return "", fmt.Errorf("unknown command: %s\n\nAvailable commands:\n  run-server\n  migrates-up\n  migrates-down\n  seed", os.Args[1])
	}
}

func ExecuteCommand(cmd command, db *sqlx.DB) error {
	switch cmd {
	case cmdMigrateUp:
		return runMigrateUp(db)
	case cmdMigrateDown:
		return runMigrateDown(db)
	case cmdSeed:
		return runSeed(db)
	case cmdRunServer:
		return runServer(db)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func InitLogger() error {
	return logger.Init(logger.Config{
		Level:        config.LogLevel,
		Format:       config.LogFormat,
		Output:       config.LogOutput,
		ReportCaller: config.LogReportCaller,
	})
}

func runMigrateUp(db *sqlx.DB) error {
	slog.Info("Running database migrations...")

	if err := migrations.Up(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	slog.Info("Migrations completed successfully")
	return nil
}

func runSeed(db *sqlx.DB) error {
	if err := seeds.Run(db); err != nil {
		return fmt.Errorf("seed data failed: %w", err)
	}
	return nil
}

func runMigrateDown(db *sqlx.DB) error {
	slog.Info("Rolling back database migrations...")

	if err := migrations.Down(db); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	slog.Info("Migrations rolled back successfully")
	return nil
}

func runServer(db *sqlx.DB) error {

	repos := repository.NewRepository(db)

	computeClientConfig := connections.Config{
		Type:    config.WordServiceType,
		BaseURL: config.WordServiceURL,
		Timeout: config.WordServiceTimeout,
	}

	computeClient, err := connections.NewComputeClient(computeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to create compute client: %w", err)
	}

	defer computeClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if resp, err := computeClient.HealthCheck(ctx); err != nil || len(resp.LoadedLanguages) == 0 {
		slog.Warn("Compute service not available or failed to load languages", logger.FieldError, err)
	} else {
		slog.Info("Compute service ready", "languages", resp.LoadedLanguages)
	}

	services := service.NewService(repos, computeClient)

	validator := validation.New()

	handlers := handler.NewHandler(services, validator)

	router := handlers.InitRoutes()

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	expiryWorker := worker.NewGameExpiryWorker(repos.Game, services.Stats, repos.TxManager, time.Minute)
	go expiryWorker.Start(workerCtx)

	swagger.SetupSwagger(router)

	srv := new(Server)

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("Starting server", "port", config.Port)
		serverErrors <- srv.Run(config.Port, router)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("Received signal, starting graceful shutdown", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		slog.Info("Server stopped gracefully")
	}

	return nil
}
