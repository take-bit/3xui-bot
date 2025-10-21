package telegram

import (
	"context"

	"3xui-bot/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx context.Context, upd tgbotapi.Update) error

type Middleware func(next HandlerFunc) HandlerFunc

func EarlyAckMiddleware(bot ports.BotPort) Middleware {

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, upd tgbotapi.Update) error {
			if cb := upd.CallbackQuery; cb != nil {
				_ = bot.AnswerCallback(ctx, cb.ID, "", false)
			}

			return next(ctx, upd)
		}
	}
}
