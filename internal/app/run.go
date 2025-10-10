package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Run запускает приложение
func Run(ctx context.Context, container *Container) error {
	// Создаем контекст с отменой для graceful shutdown
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Обрабатываем сигналы для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		container.Logger.Info("Shutting down gracefully...")
		cancel()
	}()

	// Запускаем планировщик фоновых задач
	go container.Scheduler.Start(appCtx)
	container.Logger.Info("Scheduler started")

	// Настраиваем long polling для Telegram
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := container.Bot.GetUpdatesChan(u)

	container.Logger.Info("Bot started successfully. Press Ctrl+C to stop")

	// Обрабатываем обновления
	for {
		select {
		case <-appCtx.Done():
			container.Logger.Info("Bot stopped")
			return nil
		case update := <-updates:
			// Обрабатываем обновление через роутер
			if err := container.Router.HandleUpdate(appCtx, update); err != nil {
				log.Printf("Error processing update: %v", err)
			}
		}
	}
}
