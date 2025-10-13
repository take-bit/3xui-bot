package handlers

import (
	"log/slog"
	"3xui-bot/internal/adapters/bot/telegram/ui"
	"context"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	controller interface{}
}

// NewStartHandler создает новый обработчик команды /start
func NewStartHandler(controller interface{}) *StartHandler {
	return &StartHandler{controller: controller}
}

// CanHandle проверяет, может ли обработчик обработать обновление
func (h *StartHandler) CanHandle(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == ui.CommandStart
}

// Handle обрабатывает команду /start
func (h *StartHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.getUserID(update)
	chatID := h.getChatID(update)

	slog.Info("Handling /start command for user %d", userID)

	// Создаем пользователя
	createUserDTO := usecase.CreateUserDTO{
		TelegramID:   userID,
		Username:     update.Message.From.UserName,
		FirstName:    update.Message.From.FirstName,
		LastName:     update.Message.From.LastName,
		LanguageCode: update.Message.From.LanguageCode,
	}

	_, err := h.createUser(ctx, createUserDTO)
	if err != nil {
		h.logError(err, "CreateUser")
		return err
	}

	// Отправляем приветственное сообщение с клавиатурой
	text := ui.GetWelcomeText()
	return h.sendMessageWithKeyboard(ctx, chatID, text, ui.GetWelcomeKeyboard())
}

// Вспомогательные методы
func (h *StartHandler) getUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	return 0
}

func (h *StartHandler) getChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	return 0
}

func (h *StartHandler) sendMessage(ctx context.Context, chatID int64, text string) error {
	if bot, ok := h.controller.(interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
	}); ok {
		return bot.SendMessage(ctx, chatID, text)
	}
	return nil
}

func (h *StartHandler) sendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	if bot, ok := h.controller.(interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
	}); ok {
		return bot.SendMessage(ctx, chatID, text)
	}
	return nil
}

func (h *StartHandler) createUser(ctx context.Context, dto usecase.CreateUserDTO) (interface{}, error) {
	if userUC, ok := h.controller.(interface {
		UserUC() *usecase.UserUseCase
	}); ok {
		return userUC.UserUC().CreateUser(ctx, dto)
	}
	return nil, nil
}

func (h *StartHandler) logError(err error, context string) {
	if logger, ok := h.controller.(interface {
		LogError(err error, context string)
	}); ok {
		logger.LogError(err, context)
	}
}
