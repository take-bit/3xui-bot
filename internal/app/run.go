package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(ctx context.Context, container *Container) error {
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		container.Logger.Info("Shutting down gracefully...")
		cancel()
	}()

	go container.Scheduler.Start(appCtx)
	container.Logger.Info("Scheduler started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := container.Bot.GetUpdatesChan(u)

	container.Logger.Info("Bot started successfully. Press Ctrl+C to stop")

	for {
		select {
		case <-appCtx.Done():
			container.Logger.Info("Bot stopped")

			return nil
		case update := <-updates:
			if err := container.Router.HandleUpdate(appCtx, update); err != nil {
				slog.Error("Error processing update", "error", err)
			}
		}
	}
}
