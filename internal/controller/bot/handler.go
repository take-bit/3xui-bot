package bot

import (
	"context"
	"fmt"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BaseHandler представляет базовый обработчик
type BaseHandler struct {
	UseCaseManager *usecase.UseCaseManager
	bot            *Bot
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler(useCaseManager *usecase.UseCaseManager, bot *Bot) *BaseHandler {
	return &BaseHandler{
		UseCaseManager: useCaseManager,
		bot:            bot,
	}
}

// GetUserID получает ID пользователя из обновления
func (h *BaseHandler) GetUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

// GetChatID получает ID чата из обновления
func (h *BaseHandler) GetChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return 0
}

// GetMessageID получает ID сообщения из обновления
func (h *BaseHandler) GetMessageID(update tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.MessageID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.MessageID
	}
	return 0
}

// GetText получает текст из обновления
func (h *BaseHandler) GetText(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.Text
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}
	return ""
}

// GetUsername получает имя пользователя из обновления
func (h *BaseHandler) GetUsername(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.UserName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.UserName
	}
	return ""
}

// GetFirstName получает имя пользователя из обновления
func (h *BaseHandler) GetFirstName(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.FirstName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.FirstName
	}
	return ""
}

// GetLastName получает фамилию пользователя из обновления
func (h *BaseHandler) GetLastName(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.LastName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LastName
	}
	return ""
}

// GetLanguageCode получает код языка пользователя из обновления
func (h *BaseHandler) GetLanguageCode(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.LanguageCode
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LanguageCode
	}
	return ""
}

// SendMessage отправляет сообщение пользователю
func (h *BaseHandler) SendMessage(ctx context.Context, chatID int64, text string, replyMarkup interface{}) error {
	return h.bot.SendMessage(ctx, chatID, text, replyMarkup)
}

// SendPhoto отправляет фото пользователю
func (h *BaseHandler) SendPhoto(ctx context.Context, chatID int64, photo tgbotapi.FileBytes, caption string, replyMarkup interface{}) error {
	return h.bot.SendPhoto(ctx, chatID, photo, caption, replyMarkup)
}

// AnswerCallbackQuery отвечает на callback query
func (h *BaseHandler) AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string, showAlert bool) error {
	return h.bot.AnswerCallbackQuery(ctx, callbackQueryID, text, showAlert)
}

// EditMessageText редактирует текст сообщения
func (h *BaseHandler) EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	return h.bot.EditMessageText(ctx, chatID, messageID, text, replyMarkup)
}

// DeleteMessage удаляет сообщение
func (h *BaseHandler) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	return h.bot.DeleteMessage(ctx, chatID, messageID)
}

// HandleError обрабатывает ошибку
func (h *BaseHandler) HandleError(ctx context.Context, chatID int64, err error, userMessage string) {
	// Логируем ошибку
	fmt.Printf("Error in handler: %v\n", err)

	// Отправляем сообщение пользователю
	errorMessage := "❌ Произошла ошибка. Попробуйте позже или обратитесь в поддержку."
	if userMessage != "" {
		errorMessage = userMessage
	}

	h.SendMessage(ctx, chatID, errorMessage, nil)
}

// IsUserBlocked проверяет, заблокирован ли пользователь
func (h *BaseHandler) IsUserBlocked(ctx context.Context, userID int64) bool {
	profile, err := h.UseCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		return false
	}
	return profile.IsBlocked
}

// EnsureUserRegistered обеспечивает регистрацию пользователя
func (h *BaseHandler) EnsureUserRegistered(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	username := h.GetUsername(update)
	firstName := h.GetFirstName(update)
	lastName := h.GetLastName(update)
	languageCode := h.GetLanguageCode(update)

	// Проверяем, зарегистрирован ли пользователь
	_, err := h.UseCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		// Пользователь не зарегистрирован, регистрируем его
		_, err = h.UseCaseManager.ProcessUserRegistration(ctx, userID, username, firstName, lastName, languageCode)
		if err != nil {
			return fmt.Errorf("failed to register user: %w", err)
		}
	}

	return nil
}
