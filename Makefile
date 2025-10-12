.PHONY: up-dev up-prod up-dev-b up-prod-b down logs

up-dev:
	@echo "🐳 Starting development services..."
	docker compose --profile dev up -d backend-dev db

up-prod:
	@echo "🐳 Starting production services..."
	docker compose up -d backend db

up-dev-b:
	@echo "🐳 Starting development services -b"
	docker compose --profile dev build backend-dev db
	docker compose --profile dev up -d backend-dev db

up-prod-b:
	@echo "🐳 Starting production services -b"
	docker compose build backend db
	docker compose up -d backend db

down:
	@echo "🛑 Stopping all services..."
	docker compose --profile dev down

logs:
	@echo "📋 Showing logs..."
	docker compose logs -f

logs-dev:
	@echo "📋 Showing dev logs..."
	docker compose logs -f backend-dev