package pkg

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/data/migrations"
	"github.com/dimqueue/darts/pkg/handler"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/dimqueue/darts/pkg/swagger"
	"github.com/dimqueue/darts/pkg/validation"
	"github.com/dimqueue/darts/pkg/worker"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type command string

const (
	cmdRunServer   command = "run-server"
	cmdMigrateUp   command = "migrates-up"
	cmdMigrateDown command = "migrates-down"
)

func ParseCommand() (command, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("usage: %s <command>\n\nAvailable commands:\n  run-server    - Start the HTTP server\n  migrates-up   - Run database migrations\n  migrates-down - Rollback migrations", os.Args[0])
	}

	cmd := command(os.Args[1])

	switch cmd {
	case cmdRunServer, cmdMigrateUp, cmdMigrateDown:
		return cmd, nil
	default:
		return "", fmt.Errorf("unknown command: %s\n\nAvailable commands:\n  run-server\n  migrates-up\n  migrates-down", os.Args[1])
	}
}

func ExecuteCommand(cmd command, db *sqlx.DB, config Config) error {
	switch cmd {
	case cmdMigrateUp:
		return runMigrateUp(db)
	case cmdMigrateDown:
		return runMigrateDown(db)
	case cmdRunServer:
		return runServer(db, config)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func runMigrateUp(db *sqlx.DB) error {
	logrus.Info("Running database migrations...")

	if err := migrations.Up(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	logrus.Info("Migrations completed successfully")
	return nil
}

func runMigrateDown(db *sqlx.DB) error {
	logrus.Info("Rolling back database migrations...")

	if err := migrations.Down(db); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	logrus.Info("Migrations rolled back successfully")
	return nil
}

func runServer(db *sqlx.DB, config Config) error {

	repos := repository.NewRepository(db)

	computeClient, err := connections.NewComputeClient(config.ComputeClient)
	if err != nil {
		return fmt.Errorf("failed to create compute client: %w", err)
	}

	defer computeClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if resp, err := computeClient.HealthCheck(ctx); err != nil || len(resp.LoadedLanguages) == 0 {
		logrus.Warnf("Compute service not available or failed to load languages: %v", err)
	} else {
		logrus.Infof("Compute service ready with languages: %v", resp.LoadedLanguages)
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
		logrus.Infof("Starting server on port %s", config.Port)
		serverErrors <- srv.Run(config.Port, router)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logrus.Infof("Received signal: %v. Starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		logrus.Info("Server stopped gracefully")
	}

	return nil
}
