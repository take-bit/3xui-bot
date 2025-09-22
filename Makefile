# 3X-UI Bot Makefile

.PHONY: help build run test clean deps lint format

# Переменные
BINARY_NAME=3xui-bot
BUILD_DIR=build
GO_VERSION=1.21

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Показать справку
	@echo "$(BLUE)3X-UI Bot - Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Установить зависимости
	@echo "$(BLUE)Установка зависимостей...$(NC)"
	go mod download
	go mod tidy

build: ## Собрать проект
	@echo "$(BLUE)Сборка проекта...$(NC)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bot

build-linux: ## Собрать для Linux
	@echo "$(BLUE)Сборка для Linux...$(NC)"
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/bot

build-windows: ## Собрать для Windows
	@echo "$(BLUE)Сборка для Windows...$(NC)"
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe ./cmd/bot

build-darwin: ## Собрать для macOS
	@echo "$(BLUE)Сборка для macOS...$(NC)"
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin ./cmd/bot

build-all: build-linux build-windows build-darwin ## Собрать для всех платформ

run: ## Запустить бота
	@echo "$(BLUE)Запуск бота...$(NC)"
	go run ./cmd/bot

run-dev: ## Запустить бота в режиме разработки
	@echo "$(BLUE)Запуск бота в режиме разработки...$(NC)"
	LOG_LEVEL=debug go run ./cmd/bot

test: ## Запустить тесты
	@echo "$(BLUE)Запуск тестов...$(NC)"
	go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(BLUE)Запуск тестов с покрытием...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии сохранен в coverage.html$(NC)"

lint: ## Запустить линтер
	@echo "$(BLUE)Запуск линтера...$(NC)"
	golangci-lint run

format: ## Форматировать код
	@echo "$(BLUE)Форматирование кода...$(NC)"
	go fmt ./...
	goimports -w .

clean: ## Очистить артефакты сборки
	@echo "$(BLUE)Очистка...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

install-tools: ## Установить инструменты разработки
	@echo "$(BLUE)Установка инструментов разработки...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

setup: deps install-tools ## Настроить проект
	@echo "$(GREEN)Проект настроен!$(NC)"
	@echo "$(YELLOW)Не забудьте:$(NC)"
	@echo "  1. Скопировать config.example.json в config.json"
	@echo "  2. Скопировать plans.example.json в plans.json"
	@echo "  3. Настроить переменные окружения"

docker-build: ## Собрать Docker образ
	@echo "$(BLUE)Сборка Docker образа...$(NC)"
	docker build -t 3xui-bot .

docker-run: ## Запустить в Docker
	@echo "$(BLUE)Запуск в Docker...$(NC)"
	docker run --rm -p 8080:8080 3xui-bot

docker-compose-up: ## Запустить с docker-compose
	@echo "$(BLUE)Запуск с docker-compose...$(NC)"
	docker-compose up -d

docker-compose-down: ## Остановить docker-compose
	@echo "$(BLUE)Остановка docker-compose...$(NC)"
	docker-compose down

migrate-up: ## Применить миграции
	@echo "$(BLUE)Применение миграций...$(NC)"
	# TODO: Добавить команду для миграций

migrate-down: ## Откатить миграции
	@echo "$(BLUE)Откат миграций...$(NC)"
	# TODO: Добавить команду для отката миграций

generate: ## Генерировать код
	@echo "$(BLUE)Генерация кода...$(NC)"
	go generate ./...

check: lint test ## Проверить код (линтер + тесты)

ci: deps check build ## CI pipeline

# Переменные окружения для разработки
dev-env:
	@echo "$(BLUE)Настройка переменных окружения для разработки...$(NC)"
	@echo "export BOT_TOKEN=your_bot_token_here"
	@echo "export BOT_URL=https://yourdomain.com/webhook"
	@echo "export LOG_LEVEL=debug"
	@echo "export DATABASE_URL=postgres://user:password@localhost:5432/3xui_bot?sslmode=disable"

# Показать версию Go
version:
	@echo "$(BLUE)Версия Go:$(NC)"
	@go version

# Показать информацию о модуле
info:
	@echo "$(BLUE)Информация о модуле:$(NC)"
	@go mod graph
	@echo ""
	@echo "$(BLUE)Зависимости:$(NC)"
	@go list -m all
