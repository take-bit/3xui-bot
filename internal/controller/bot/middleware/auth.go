package middleware

import (
	"context"
	"fmt"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AuthMiddleware представляет middleware для аутентификации
type AuthMiddleware struct {
	useCaseManager *usecase.UseCaseManager
}

// NewAuthMiddleware создает новый middleware для аутентификации
func NewAuthMiddleware(useCaseManager *usecase.UseCaseManager) *AuthMiddleware {
	return &AuthMiddleware{
		useCaseManager: useCaseManager,
	}
}

// Process обрабатывает обновление с проверкой аутентификации
func (m *AuthMiddleware) Process(ctx context.Context, update tgbotapi.Update, next func(ctx context.Context, update tgbotapi.Update) error) error {
	// Получаем ID пользователя
	userID := m.getUserID(update)
	if userID == 0 {
		return fmt.Errorf("unable to get user ID from update")
	}

	// Проверяем, зарегистрирован ли пользователь
	profile, err := m.useCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		// Пользователь не зарегистрирован, регистрируем его
		username := m.getUsername(update)
		firstName := m.getFirstName(update)
		lastName := m.getLastName(update)
		languageCode := m.getLanguageCode(update)

		_, err = m.useCaseManager.ProcessUserRegistration(ctx, userID, username, firstName, lastName, languageCode)
		if err != nil {
			return fmt.Errorf("failed to register user: %w", err)
		}
	} else if profile.IsBlocked {
		// Пользователь заблокирован
		chatID := m.getChatID(update)
		message := "🚫 Ваш аккаунт заблокирован. Обратитесь в поддержку для получения дополнительной информации."

		// Отправляем сообщение о блокировке
		// TODO: Отправить сообщение пользователю
		_ = chatID
		_ = message

		return fmt.Errorf("user is blocked")
	}

	// Продолжаем обработку
	return next(ctx, update)
}

// getUserID получает ID пользователя из обновления
func (m *AuthMiddleware) getUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

// getChatID получает ID чата из обновления
func (m *AuthMiddleware) getChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return 0
}

// getUsername получает имя пользователя из обновления
func (m *AuthMiddleware) getUsername(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.UserName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.UserName
	}
	return ""
}

// getFirstName получает имя пользователя из обновления
func (m *AuthMiddleware) getFirstName(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.FirstName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.FirstName
	}
	return ""
}

// getLastName получает фамилию пользователя из обновления
func (m *AuthMiddleware) getLastName(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.LastName
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LastName
	}
	return ""
}

// getLanguageCode получает код языка пользователя из обновления
func (m *AuthMiddleware) getLanguageCode(update tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.LanguageCode
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LanguageCode
	}
	return ""
}
