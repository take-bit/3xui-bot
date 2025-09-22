package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LoggingMiddleware представляет middleware для логирования
type LoggingMiddleware struct{}

// NewLoggingMiddleware создает новый middleware для логирования
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// Process обрабатывает обновление с логированием
func (m *LoggingMiddleware) Process(ctx context.Context, update tgbotapi.Update, next func(ctx context.Context, update tgbotapi.Update) error) error {
	start := time.Now()

	// Логируем входящее обновление
	m.logUpdate(update)

	// Обрабатываем обновление
	err := next(ctx, update)

	// Логируем результат
	duration := time.Since(start)
	m.logResult(update, err, duration)

	return err
}

// logUpdate логирует входящее обновление
func (m *LoggingMiddleware) logUpdate(update tgbotapi.Update) {
	var logMessage string

	if update.Message != nil {
		message := update.Message
		logMessage = fmt.Sprintf("Message from %s (%d): %s",
			message.From.UserName, message.From.ID, message.Text)
	} else if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		logMessage = fmt.Sprintf("Callback from %s (%d): %s",
			callback.From.UserName, callback.From.ID, callback.Data)
	} else {
		logMessage = "Unknown update type"
	}

	log.Printf("📨 %s", logMessage)
}

// logResult логирует результат обработки
func (m *LoggingMiddleware) logResult(update tgbotapi.Update, err error, duration time.Duration) {
	var logMessage string

	if err != nil {
		logMessage = fmt.Sprintf("❌ Error processing update: %v (took %v)", err, duration)
	} else {
		logMessage = fmt.Sprintf("✅ Successfully processed update (took %v)", duration)
	}

	log.Printf("📤 %s", logMessage)
}
