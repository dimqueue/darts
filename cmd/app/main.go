package main

import (
	"context"
	"errors"
	"github.com/dimqueue/darts/pkg"
	"github.com/dimqueue/darts/pkg/data/migrations"
	"github.com/dimqueue/darts/pkg/handler"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/dimqueue/darts/pkg/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/siruspen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// @title           Darts API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logrus.Println("⚠️  No .env file found, using system vars")
		}
	}

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs %s", err.Error())
	}

	if len(os.Args[1:]) != 1 {
		err := errors.New("expected exactly one argument")
		logrus.Fatalf("wrong CLI arguments: %v", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Username: viper.GetString("db.username"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	switch os.Args[1] {
	case "migrates-up":
		if err := migrations.Up(db); err != nil {
			logrus.Fatalf("failed to migrates-up: %v", err)
		}
	case "migrates-down":
		if err := migrations.Down(db); err != nil {
			logrus.Fatalf("failed to migrates-down: %v", err)
		}
	case "run-server":
		var srv *pkg.Server
		go func() {
			if srv, err = RunServer(db); err != nil {
				logrus.Fatalf("failed to run-server: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		logrus.Print("application shutting down")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Fatalf("error occured on server shutting down: %s", err.Error())
		}

		if err := db.Close(); err != nil {
			logrus.Fatalf("error occured on db connection close: %s", err.Error())
		}
	}

}

func RunServer(db *sqlx.DB) (*pkg.Server, error) {
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	swagger.SetupSwagger(router)

	srv := new(pkg.Server)
	if err := srv.Run("8080", router); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

	return srv, nil
}

func initConfig() error {
	viper.AutomaticEnv()
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
