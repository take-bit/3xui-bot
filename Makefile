.PHONY: help build run dev migrate migrate-down migrate-fresh clean docker-up docker-down test setup

help:
	@echo "Available commands:"
	@echo "  make setup         - Full setup (docker + migrate)"
	@echo "  make build         - Build the bot binary"
	@echo "  make run           - Run the bot"
	@echo "  make dev           - Run bot in development mode"
	@echo "  make migrate       - Apply database migrations (drop + create)"
	@echo "  make migrate-down  - Drop all tables"
	@echo "  make migrate-fresh - Fresh migration (down + up)"
	@echo "  make docker-up     - Start PostgreSQL in Docker"
	@echo "  make docker-down   - Stop PostgreSQL"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"

setup: docker-up migrate
	@echo "✅ Setup complete! You can now run: make run"

build:
	@echo "Building bot..."
	go build -o bot cmd/bot/main.go
	@echo "✅ Build complete: ./bot"

run: build
	@echo "Starting bot..."
	./bot

dev:
	@echo "Running bot in development mode..."
	go run cmd/bot/main.go

migrate:
	@echo "🔄 Dropping all tables..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/000_drop_all.sql 2>/dev/null || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/000_drop_all.sql 2>/dev/null || true
	@echo "📋 Applying schema..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/001_complete_schema.sql || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/001_complete_schema.sql
	@echo "🌱 Seeding plans..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/002_seed_plans.sql || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/002_seed_plans.sql
	@echo "✅ Migrations applied successfully"

migrate-down:
	@echo "🔄 Rolling back migrations (dropping all tables)..."
	@docker exec -i 3xui_bot_db psql -U bot_user -d 3xui_bot < migrations/000_drop_all.sql || \
	psql -h localhost -U bot_user -d 3xui_bot -f migrations/000_drop_all.sql
	@echo "✅ All tables dropped"

migrate-fresh: migrate-down migrate
	@echo "✅ Fresh migration complete!"

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

