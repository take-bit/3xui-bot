package service

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageService encapsulates all Telegram message operations
type MessageService struct {
	bot *tgbotapi.BotAPI
}

// NewMessageService creates a new MessageService
func NewMessageService(bot *tgbotapi.BotAPI) *MessageService {
	return &MessageService{
		bot: bot,
	}
}

// SendMessage sends a simple text message
func (s *MessageService) SendMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	return err
}

// SendMessageWithKeyboard sends a text message with inline keyboard
func (s *MessageService) SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}
	_, err := s.bot.Send(msg)
	return err
}

// SendMessageWithMarkdownV2 sends a text message with MarkdownV2 formatting
func (s *MessageService) SendMessageWithMarkdownV2(ctx context.Context, chatID int64, text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	if keyboard != nil {
		if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = kb
		}
	}
	_, err := s.bot.Send(msg)
	return err
}

// EditMessageText edits an existing message
func (s *MessageService) EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if replyMarkup != nil {
		if kb, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
			editMsg.ReplyMarkup = &kb
		}
	}
	_, err := s.bot.Send(editMsg)
	return err
}

// EditMessageWithMarkdownV2 edits an existing message with MarkdownV2 formatting
func (s *MessageService) EditMessageWithMarkdownV2(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "MarkdownV2"
	if replyMarkup != nil {
		if kb, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
			editMsg.ReplyMarkup = &kb
		}
	}
	_, err := s.bot.Send(editMsg)
	return err
}

// DeleteMessage deletes a message
func (s *MessageService) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := s.bot.Request(msg)
	return err
}

// DeleteAndSendMessage deletes old message and sends new one
func (s *MessageService) DeleteAndSendMessage(ctx context.Context, chatID int64, messageID int, text string, keyboard interface{}) error {
	// Delete old message (ignore error)
	_ = s.DeleteMessage(ctx, chatID, messageID)

	// Send new message
	return s.SendMessageWithKeyboard(ctx, chatID, text, keyboard)
}

// DeleteAndSendMessageWithMarkdownV2 deletes old message and sends new one with MarkdownV2
func (s *MessageService) DeleteAndSendMessageWithMarkdownV2(ctx context.Context, chatID int64, messageID int, text string, keyboard interface{}) error {
	// Delete old message (ignore error)
	_ = s.DeleteMessage(ctx, chatID, messageID)

	// Send new message with MarkdownV2
	return s.SendMessageWithMarkdownV2(ctx, chatID, text, keyboard)
}

// AnswerCallbackQuery answers a callback query
func (s *MessageService) AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string, showAlert bool) error {
	ack := tgbotapi.NewCallback(callbackQueryID, text)
	ack.ShowAlert = showAlert
	_, err := s.bot.Request(ack)
	return err
}

func (s *MessageService) SendPhotoWithMarkdown(ctx context.Context, chatID int64, imagePath string, caption string, keyboard interface{}) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	photo.Caption = escapeMarkdownV2(caption)
	photo.ParseMode = "MarkdownV2"
	if keyboard != nil {
		if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
			photo.ReplyMarkup = kb
		}
	}
	_, err := s.bot.Send(photo)
	if err != nil {
		slog.Error("Failed to send photo with markdown", "chat_id", chatID, "image_path", imagePath, "error", err)
	}
	return err
}

func (s *MessageService) SendPhotoWithPreEscapedMarkdown(ctx context.Context, chatID int64, imagePath string, caption string, keyboard interface{}) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imagePath))
	photo.Caption = caption
	photo.ParseMode = "MarkdownV2"
	if keyboard != nil {
		if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
			photo.ReplyMarkup = kb
		}
	}
	_, err := s.bot.Send(photo)
	if err != nil {
		slog.Error("Failed to send photo with pre-escaped markdown", "chat_id", chatID, "image_path", imagePath, "error", err)
	}
	return err
}

var (
	markdownV2Replacer = strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	linkPattern = regexp.MustCompile(`\[[^\]]+\]\([^)]+\)`)
)

func escapeMarkdownV2(text string) string {
	links := make(map[string]string)
	matches := linkPattern.FindAllString(text, -1)
	for i, match := range matches {
		placeholder := fmt.Sprintf("\x00LINK%d\x00", i)
		links[placeholder] = match
		text = strings.Replace(text, match, placeholder, 1)
	}
	text = markdownV2Replacer.Replace(text)
	for placeholder, original := range links {
		text = strings.Replace(text, placeholder, original, 1)
	}
	return text
}
