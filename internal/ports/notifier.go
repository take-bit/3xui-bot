package ports

import (
	"context"
	"io"
)

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

	// SendPhoto отправляет фото с подписью
	SendPhoto(ctx context.Context, chatID int64, photoFileID string, caption string, markup interface{}) error

	// SendPhotoFromReader отправляет фото из io.Reader
	SendPhotoFromReader(ctx context.Context, chatID int64, photoReader io.Reader, caption string, markup interface{}) error

	// SendPhotoFromFile отправляет фото из файла
	SendPhotoFromFile(ctx context.Context, chatID int64, photoPath string, caption string, markup interface{}) error

	// SendPhotoFromFileWithParseMode отправляет фото из файла с указанным режимом парсинга
	SendPhotoFromFileWithParseMode(ctx context.Context, chatID int64, photoPath string, caption string, parseMode string, markup interface{}) error

	// EditMessagePhoto редактирует сообщение с фото
	EditMessagePhoto(ctx context.Context, chatID int64, messageID int, photoFileID string, caption string, markup interface{}) error
}
