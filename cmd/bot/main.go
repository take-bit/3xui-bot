package main

import (
	"context"
	"flag"
	"log"
	"os"

	"3xui-bot/internal/app"
)

func main() {
	ctx := context.Background()

	// Получаем путь к конфигу с приоритетом: флаг > env > default
	configPath := getConfigPath()
	log.Printf("Using config file: %s", configPath)

	// Создаем контейнер зависимостей
	container, err := app.NewContainer(ctx, configPath)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer container.Close()

	// Запускаем приложение
	if err := app.Run(ctx, container); err != nil {
		log.Fatalf("App error: %v", err)
	}
}

// getConfigPath возвращает путь к файлу конфигурации
// Приоритет: флаг -config > env CONFIG_PATH > default configs/config.json
func getConfigPath() string {
	// 1. Проверяем флаг командной строки (высший приоритет)
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "Path to config file")
	flag.Parse()

	if configFlag != "" {
		return configFlag
	}

	// 2. Проверяем переменную окружения
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	// 3. Используем значение по умолчанию
	return "configs/config.json"
}
