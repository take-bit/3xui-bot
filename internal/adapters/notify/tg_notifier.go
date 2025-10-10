package notify

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
