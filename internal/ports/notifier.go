package ports

import (
	"context"
	"io"
)

type Notifier interface {
	Send(ctx context.Context, chatID int64, text string, markup interface{}) error

	SendWithParseMode(ctx context.Context, chatID int64, text string, parseMode string, markup interface{}) error

	EditMessage(ctx context.Context, chatID int64, messageID int, text string, markup interface{}) error

	DeleteMessage(ctx context.Context, chatID int64, messageID int) error

	SendPhoto(ctx context.Context, chatID int64, photoFileID string, caption string, markup interface{}) error

	SendPhotoFromReader(ctx context.Context, chatID int64, photoReader io.Reader, caption string, markup interface{}) error

	SendPhotoFromFile(ctx context.Context, chatID int64, photoPath string, caption string, markup interface{}) error

	SendPhotoFromFileWithParseMode(ctx context.Context, chatID int64, photoPath string, caption string, parseMode string, markup interface{}) error

	EditMessagePhoto(ctx context.Context, chatID int64, messageID int, photoFileID string, caption string, markup interface{}) error
}
