.PHONY: help build run migrate migrate-down clean docker-up docker-down test

help:
	@echo "Available commands:"
	@echo "  make build        - Build the bot binary"
	@echo "  make run          - Run the bot"
	@echo "  make migrate      - Apply database migrations"
	@echo "  make migrate-down - Rollback last migration"
	@echo "  make docker-up    - Start PostgreSQL in Docker"
	@echo "  make docker-down  - Stop PostgreSQL"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"

build:
	@echo "Building bot..."
	go build -o bot cmd/bot/main.go
	@echo "✅ Build complete: ./bot"

run: build
	@echo "Starting bot..."
	./bot

migrate:
	@echo "Applying migrations..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/001_complete_schema.sql || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/001_complete_schema.sql
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/002_seed_plans.sql || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/002_seed_plans.sql
	@echo "✅ Migrations applied"

migrate-down:
	@echo "Rolling back migrations..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" || \
	psql -h localhost -U bot_user -d 3xui_bot -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "✅ Database reset"

docker-up:
	@echo "Starting PostgreSQL..."
	docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@echo "✅ PostgreSQL is running"

docker-down:
	@echo "Stopping PostgreSQL..."
	docker-compose down
	@echo "✅ PostgreSQL stopped"

clean:
	@echo "Cleaning..."
	rm -f bot marzbanTest
	@echo "✅ Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

