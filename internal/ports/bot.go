package ports

import "context"

type BotPort interface {
	Send(ctx context.Context, chatID int64, text string, markup any) error
	Edit(ctx context.Context, chatID int64, messageID int, text string, markup any) error
	AnswerCallback(ctx context.Context, callbackQueryID string, text string, showAlert bool) error
}
