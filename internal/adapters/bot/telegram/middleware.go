package telegram

import (
	"context"

	"3xui-bot/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlerFunc тип функции-обработчика
type HandlerFunc func(ctx context.Context, upd tgbotapi.Update) error

// Middleware тип middleware функции
type Middleware func(next HandlerFunc) HandlerFunc

// EarlyAckMiddleware отправляет ранний ACK на callback query до основной обработки
func EarlyAckMiddleware(bot ports.BotPort) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, upd tgbotapi.Update) error {
			// Если это callback query, отправляем ACK сразу
			if cb := upd.CallbackQuery; cb != nil {
				_ = bot.AnswerCallback(ctx, cb.ID, "", false)
			}
			return next(ctx, upd)
		}
	}
}
