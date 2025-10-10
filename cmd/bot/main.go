package main

import (
	"context"
	"log"

	"3xui-bot/internal/app"
)

func main() {
	ctx := context.Background()

	// Создаем контейнер зависимостей
	container, err := app.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer container.Close()

	// Запускаем приложение
	if err := app.Run(ctx, container); err != nil {
		log.Fatalf("App error: %v", err)
	}
}
