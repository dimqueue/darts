package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimqueue/darts/pkg"
	"github.com/dimqueue/darts/pkg/data/migrations"
	"github.com/dimqueue/darts/pkg/handler"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/dimqueue/darts/pkg/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title           Darts API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := run(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func run() error {
	// Load environment and config
	if err := loadEnv(); err != nil {
		return err
	}

	if err := loadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse command
	cmd, err := parseCommand()
	if err != nil {
		return err
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Execute command
	return executeCommand(cmd, db)
}

// loadEnv loads .env file in non-production environments
func loadEnv() error {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logrus.Info("No .env file found, using system environment variables")
		}
	}
	return nil
}

func loadConfig() error {
	viper.AutomaticEnv()
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Set defaults
	viper.SetDefault("port", "8080")
	viper.SetDefault("db.sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	logrus.Infof("Loaded config from: %s", viper.ConfigFileUsed())
	return nil
}

func connectDB() (*sqlx.DB, error) {
	config := repository.Config{
		Host:     viper.GetString("db.host"),
		Username: viper.GetString("db.username"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	db, err := repository.NewPostgresDB(config)
	if err != nil {
		return nil, err
	}

	logrus.Info("Database connection established")
	return db, nil
}

type command string

const (
	cmdRunServer   command = "run-server"
	cmdMigrateUp   command = "migrates-up"
	cmdMigrateDown command = "migrates-down"
)

func parseCommand() (command, error) {
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

func executeCommand(cmd command, db *sqlx.DB) error {
	switch cmd {
	case cmdMigrateUp:
		return runMigrateUp(db)
	case cmdMigrateDown:
		return runMigrateDown(db)
	case cmdRunServer:
		return runServer(db)
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

func runServer(db *sqlx.DB) error {
	repos := repository.NewRepository(db)

	services := service.NewService(repos)

	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	swagger.SetupSwagger(router)

	srv := new(pkg.Server)

	serverErrors := make(chan error, 1)

	go func() {
		logrus.Infof("Starting server on port %s", viper.GetString("port"))
		serverErrors <- srv.Run(viper.GetString("port"), router)
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
