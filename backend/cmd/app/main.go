package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dimqueue/darts/pkg"
	_ "github.com/lib/pq"
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
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := pkg.LoadEnv(); err != nil {
		return err
	}

	if err := pkg.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := pkg.InitLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cmd, err := pkg.ParseCommand()
	if err != nil {
		return err
	}

	db, err := pkg.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	return pkg.ExecuteCommand(cmd, db)
}
