package pkg

import (
	"fmt"
	"os"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port          string
	DB            repository.Config
	ComputeClient connections.Config
}

func LoadEnv() error {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logrus.Info("No .env file found, using system environment variables")
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

	logrus.Infof("Loaded config from: %s", viper.ConfigFileUsed())
	return nil
}

func GetConfig() Config {
	return Config{
		Port: viper.GetString("port"),
		DB: repository.Config{
			Host:     viper.GetString("db.host"),
			Username: viper.GetString("db.username"),
			Port:     viper.GetString("db.port"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		ComputeClient: connections.Config{
			Type:    viper.GetString("word-service.type"),
			BaseURL: viper.GetString("word-service.baseURL"),
			Timeout: viper.GetInt("word-service.timeout"),
		},
	}
}

func ConnectDB(config repository.Config) (*sqlx.DB, error) {
	db, err := repository.NewPostgresDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Database connection established")
	return db, nil
}
