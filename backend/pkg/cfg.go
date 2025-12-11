package pkg

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadEnv() error {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load("../.env"); err != nil {
			slog.Info("No .env file found, using system environment variables")
		}
	}
	return nil
}

func LoadConfig() error {
	viper.AutomaticEnv()
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	slog.Info("Loaded config", "file", viper.ConfigFileUsed())

	config.Load()

	return nil
}

func ConnectDB() (*sqlx.DB, error) {
	dbConfig := repository.Config{
		Host:     config.DBHost,
		Port:     config.DBPort,
		Username: config.DBUser,
		DBName:   config.DBName,
		SSLMode:  config.DBSSLMode,
		Password: config.DBPassword,
	}

	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Database connection established")
	return db, nil
}
