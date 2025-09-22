package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlerInterface представляет интерфейс для обработчиков
type HandlerInterface interface {
	Handle(ctx context.Context, update tgbotapi.Update) error
	Command() string
	Description() string
}

// BaseHandlerInterface представляет интерфейс для базового обработчика
type BaseHandlerInterface interface {
	GetUserID(update tgbotapi.Update) int64
	GetChatID(update tgbotapi.Update) int64
	GetMessageID(update tgbotapi.Update) int
	GetText(update tgbotapi.Update) string
	GetUsername(update tgbotapi.Update) string
	GetFirstName(update tgbotapi.Update) string
	GetLastName(update tgbotapi.Update) string
	GetLanguageCode(update tgbotapi.Update) string
	SendMessage(ctx context.Context, chatID int64, text string, replyMarkup interface{}) error
	SendPhoto(ctx context.Context, chatID int64, photo tgbotapi.FileBytes, caption string, replyMarkup interface{}) error
	AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string, showAlert bool) error
	EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error
	DeleteMessage(ctx context.Context, chatID int64, messageID int) error
	HandleError(ctx context.Context, chatID int64, err error, userMessage string)
	IsUserBlocked(ctx context.Context, userID int64) bool
	EnsureUserRegistered(ctx context.Context, update tgbotapi.Update) error
}
