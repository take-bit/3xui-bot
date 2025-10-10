package ports

import "context"

// Notifier интерфейс для отправки уведомлений
type Notifier interface {
	// Send отправляет сообщение пользователю
	Send(ctx context.Context, chatID int64, text string, markup interface{}) error

	// SendWithParseMode отправляет сообщение с указанным режимом парсинга
	SendWithParseMode(ctx context.Context, chatID int64, text string, parseMode string, markup interface{}) error

	// EditMessage редактирует существующее сообщение
	EditMessage(ctx context.Context, chatID int64, messageID int, text string, markup interface{}) error

	// DeleteMessage удаляет сообщение
	DeleteMessage(ctx context.Context, chatID int64, messageID int) error
}
