package notify

import (
	"context"
	"fmt"
	"io"

	"3xui-bot/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Гарантируем, что TelegramNotifier реализует оба интерфейса
var _ ports.Notifier = (*TelegramNotifier)(nil)
var _ ports.BotPort = (*TelegramNotifier)(nil)

// TelegramNotifier адаптер для отправки уведомлений через Telegram
type TelegramNotifier struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramNotifier создает новый Telegram notifier
func NewTelegramNotifier(bot *tgbotapi.BotAPI) *TelegramNotifier {
	return &TelegramNotifier{
		bot: bot,
	}
}

// Send отправляет сообщение пользователю
func (n *TelegramNotifier) Send(ctx context.Context, chatID int64, text string, markup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// SendWithParseMode отправляет сообщение с указанным режимом парсинга
func (n *TelegramNotifier) SendWithParseMode(ctx context.Context, chatID int64, text string, parseMode string, markup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// EditMessage редактирует существующее сообщение
func (n *TelegramNotifier) EditMessage(ctx context.Context, chatID int64, messageID int, text string, markup interface{}) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = &kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to edit message: %w", err)
	}

	return nil
}

// DeleteMessage удаляет сообщение
func (n *TelegramNotifier) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)

	if _, err := n.bot.Request(msg); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// SendPhoto отправляет фото с подписью
func (n *TelegramNotifier) SendPhoto(ctx context.Context, chatID int64, photoFileID string, caption string, markup interface{}) error {
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(photoFileID))
	msg.Caption = caption

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}

	return nil
}

// SendPhotoFromReader отправляет фото из io.Reader
func (n *TelegramNotifier) SendPhotoFromReader(ctx context.Context, chatID int64, photoReader io.Reader, caption string, markup interface{}) error {
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
		Name:   "photo.jpg",
		Reader: photoReader,
	})
	msg.Caption = caption

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send photo from reader: %w", err)
	}

	return nil
}

// SendPhotoFromFile отправляет фото из файла
func (n *TelegramNotifier) SendPhotoFromFile(ctx context.Context, chatID int64, photoPath string, caption string, markup interface{}) error {
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(photoPath))
	msg.Caption = caption

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send photo from file: %w", err)
	}

	return nil
}

// SendPhotoFromFileWithParseMode отправляет фото из файла с указанным режимом парсинга
func (n *TelegramNotifier) SendPhotoFromFileWithParseMode(ctx context.Context, chatID int64, photoPath string, caption string, parseMode string, markup interface{}) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(photoPath))
	photo.Caption = caption
	photo.ParseMode = parseMode

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			photo.ReplyMarkup = kb
		}
	}

	if _, err := n.bot.Send(photo); err != nil {
		return fmt.Errorf("failed to send photo from file with parse mode: %w", err)
	}

	return nil
}

// EditMessagePhoto редактирует сообщение с фото
func (n *TelegramNotifier) EditMessagePhoto(ctx context.Context, chatID int64, messageID int, photoFileID string, caption string, markup interface{}) error {
	// Для редактирования фото нужно использовать EditMessageText с новой подписью
	// или удалить старое сообщение и отправить новое
	msg := tgbotapi.NewEditMessageCaption(chatID, messageID, caption)

	if markup != nil {
		if kb, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = &kb
		}
	}

	if _, err := n.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to edit message photo caption: %w", err)
	}

	return nil
}

// ============================================================================
// РЕАЛИЗАЦИЯ ports.BotPort
// ============================================================================

// Edit редактирует существующее сообщение (реализация BotPort)
func (n *TelegramNotifier) Edit(ctx context.Context, chatID int64, messageID int, text string, markup any) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	if markup != nil {
		if m, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = &m
		}
	}

	_, err := n.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to edit message: %w", err)
	}
	return nil
}

// AnswerCallback отвечает на callback query (реализация BotPort)
func (n *TelegramNotifier) AnswerCallback(ctx context.Context, callbackQueryID string, text string, showAlert bool) error {
	ack := tgbotapi.NewCallback(callbackQueryID, text)
	ack.ShowAlert = showAlert

	_, err := n.bot.Request(ack)
	if err != nil {
		return fmt.Errorf("failed to answer callback: %w", err)
	}
	return nil
}
