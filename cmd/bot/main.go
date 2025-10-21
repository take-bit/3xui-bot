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

	configPath := getConfigPath()
	log.Printf("Using config file: %s", configPath)

	container, err := app.NewContainer(ctx, configPath)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer container.Close()

	if err := app.Run(ctx, container); err != nil {
		log.Fatalf("App error: %v", err)
	}
}

func getConfigPath() string {
	var configFlag string
	flag.StringVar(&configFlag, "config", "", "Path to config file")
	flag.Parse()

	if configFlag != "" {

		return configFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {

		return envPath
	}

	return "configs/config.json"
}
