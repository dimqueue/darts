.PHONY: up-dev up-prod up-dev-b up-prod-b down logs-dev local-run

up-dev:
	@echo "🐳 Starting development services..."
	docker compose --profile dev up -d backend-dev db

up-prod:
	@echo "🐳 Starting production services..."
	docker compose up -d backend db

up-dev-b:
	@echo "🐳 Starting development services -b"
	swag init -g cmd/app/main.go
	docker compose --profile dev build backend-dev db
	docker compose --profile dev up -d backend-dev db

up-prod-b:
	@echo "🐳 Starting production services -b"
	docker compose build backend db
	docker compose up -d backend db

down:
	@echo "🛑 Stopping all services..."
	docker compose --profile dev down

logs-dev:
	@echo "📋 Showing dev logs..."
	docker compose logs -f backend-dev

local-run:
	go run -tags=dev ./cmd/app run-server