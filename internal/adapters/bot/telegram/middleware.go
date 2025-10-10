package telegram

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Middleware функция middleware
type Middleware func(next Handler) Handler

// Handler обработчик обновлений
type Handler func(ctx context.Context, update tgbotapi.Update) error

// LoggingMiddleware логирует все обновления
func LoggingMiddleware(next Handler) Handler {
	return func(ctx context.Context, update tgbotapi.Update) error {
		start := time.Now()

		var userID int64
		if update.Message != nil && update.Message.From != nil {
			userID = update.Message.From.ID
		} else if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
			userID = update.CallbackQuery.From.ID
		}

		log.Printf("[Middleware] Processing update from user %d", userID)

		err := next(ctx, update)

		duration := time.Since(start)
		if err != nil {
			log.Printf("[Middleware] Error processing update: %v (took %v)", err, duration)
		} else {
			log.Printf("[Middleware] Successfully processed update (took %v)", duration)
		}

		return err
	}
}

// RecoveryMiddleware восстанавливается после паники
func RecoveryMiddleware(next Handler) Handler {
	return func(ctx context.Context, update tgbotapi.Update) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Recovery] Panic recovered: %v", r)
				err = nil // Не роняем бота
			}
		}()

		return next(ctx, update)
	}
}

// Chain объединяет middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
