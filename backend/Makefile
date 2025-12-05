.PHONY: build swagger up-db migrates-up migrates-down run run-only dev-run dev-run-only down clean rebuild

build:
	docker compose build backend

swagger:
	swag init -g cmd/app/main.go

up-db:
	docker compose up -d db

migrates-up: up-db
	docker compose run --rm backend migrates-up

migrates-down: up-db
	docker compose run --rm backend migrates-down

run: up-db build
	docker compose up backend run-server

run-only: up-db
	docker compose up backend run-server

dev-run:
	BUILD_MODE=dev $(MAKE) run

dev-run-only:
	BUILD_MODE=dev $(MAKE) run-only

down:
	docker compose down

clean:
	docker compose down -v
	docker system prune -f

rebuild: clean build run