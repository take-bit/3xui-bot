package notify

import (
	"context"
	"fmt"
	"io"

	"3xui-bot/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ ports.Notifier = (*TelegramNotifier)(nil)
var _ ports.BotPort = (*TelegramNotifier)(nil)

type TelegramNotifier struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramNotifier(bot *tgbotapi.BotAPI) *TelegramNotifier {

	return &TelegramNotifier{
		bot: bot,
	}
}

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

func (n *TelegramNotifier) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)

	if _, err := n.bot.Request(msg); err != nil {

		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

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

func (n *TelegramNotifier) EditMessagePhoto(ctx context.Context, chatID int64, messageID int, photoFileID string, caption string, markup interface{}) error {
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

func (n *TelegramNotifier) AnswerCallback(ctx context.Context, callbackQueryID string, text string, showAlert bool) error {
	ack := tgbotapi.NewCallback(callbackQueryID, text)
	ack.ShowAlert = showAlert

	_, err := n.bot.Request(ack)
	if err != nil {

		return fmt.Errorf("failed to answer callback: %w", err)
	}

	return nil
}
