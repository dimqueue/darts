package main

import (
	"darts"
	"darts/pkg/handler"
	"log"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(darts.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
