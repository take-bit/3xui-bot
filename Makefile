.PHONY: help build run test clean fmt lint migrate-up migrate-down db-create db-drop dev install

# Переменные
BINARY_NAME=bot
BINARY_PATH=./bin/$(BINARY_NAME)
CMD_PATH=./cmd/bot/main.go

# Цвета для вывода
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Показать помощь
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Собрать бинарник
	@echo "$(GREEN)Сборка проекта...$(NC)"
	@go build -o $(BINARY_PATH) $(CMD_PATH)
	@echo "$(GREEN)✓ Бинарник создан: $(BINARY_PATH)$(NC)"

run: build ## Запустить бота
	@echo "$(GREEN)Запуск бота...$(NC)"
	@$(BINARY_PATH)

dev: ## Запустить в режиме разработки (без сборки)
	@echo "$(GREEN)Запуск в dev режиме...$(NC)"
	@go run $(CMD_PATH)

test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	@go test -v -race -cover ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Отчет о покрытии: coverage.html$(NC)"

clean: ## Удалить бинарники и временные файлы
	@echo "$(GREEN)Очистка...$(NC)"
	@rm -rf $(BINARY_PATH)
	@rm -rf coverage.out coverage.html
	@echo "$(GREEN)✓ Очистка завершена$(NC)"

fmt: ## Форматировать код
	@echo "$(GREEN)Форматирование кода...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Код отформатирован$(NC)"

lint: ## Проверить линтером
	@echo "$(GREEN)Запуск линтера...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Проверка линтером завершена$(NC)"

install: ## Установить зависимости
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Зависимости установлены$(NC)"

db-create: ## Создать базу данных
	@echo "$(GREEN)Создание базы данных...$(NC)"
	@createdb $(DB_NAME) || true
	@echo "$(GREEN)✓ База данных создана$(NC)"

db-drop: ## Удалить базу данных
	@echo "$(GREEN)Удаление базы данных...$(NC)"
	@dropdb $(DB_NAME) || true
	@echo "$(GREEN)✓ База данных удалена$(NC)"

migrate-up: ## Применить миграции
	@echo "$(GREEN)Применение миграций...$(NC)"
	@psql -d $(DB_NAME) -f migrations/001_complete_schema.sql
	@echo "$(GREEN)✓ Миграции применены$(NC)"

migrate-down: db-drop db-create ## Откатить все миграции (пересоздать БД)
	@echo "$(GREEN)✓ База данных пересоздана$(NC)"

docker-up: ## Запустить PostgreSQL в Docker
	@echo "$(GREEN)Запуск PostgreSQL в Docker...$(NC)"
	@docker run -d \
		--name 3xui-bot-postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=3xui_bot \
		-p 5432:5432 \
		postgres:14-alpine
	@echo "$(GREEN)✓ PostgreSQL запущен$(NC)"

docker-down: ## Остановить PostgreSQL в Docker
	@echo "$(GREEN)Остановка PostgreSQL...$(NC)"
	@docker stop 3xui-bot-postgres || true
	@docker rm 3xui-bot-postgres || true
	@echo "$(GREEN)✓ PostgreSQL остановлен$(NC)"

setup: install db-create migrate-up ## Полная настройка проекта
	@echo "$(GREEN)✓ Проект настроен и готов к работе!$(NC)"
	@echo "$(GREEN)Скопируйте env.sample в .env и заполните переменные$(NC)"

# Переменные окружения по умолчанию для make команд
DB_NAME ?= 3xui_bot

