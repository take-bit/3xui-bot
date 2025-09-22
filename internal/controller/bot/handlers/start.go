package handlers

import (
	"context"
	"fmt"

	"3xui-bot/internal/interfaces"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewStartHandler создает новый обработчик команды /start
func NewStartHandler(useCaseManager *usecase.UseCaseManager) *StartHandler {
	return &StartHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /start
func (h *StartHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)
	username := h.GetUsername(update)
	firstName := h.GetFirstName(update)
	lastName := h.GetLastName(update)
	languageCode := h.GetLanguageCode(update)

	// Проверяем, есть ли реферальный код в команде
	referrerID := h.extractReferrerID(update)

	var user *usecase.UserProfile
	var err error

	if referrerID > 0 {
		// Регистрация по реферальной ссылке
		user, err = h.processReferralRegistration(ctx, userID, referrerID, username, firstName, lastName, languageCode)
	} else {
		// Обычная регистрация
		user, err = h.processNormalRegistration(ctx, userID, username, firstName, lastName, languageCode)
	}

	if err != nil {
		h.HandleError(ctx, chatID, err, "❌ Ошибка при регистрации. Попробуйте позже.")
		return err
	}

	// Отправляем приветственное сообщение
	return h.sendWelcomeMessage(ctx, chatID, user, referrerID > 0)
}

// extractReferrerID извлекает ID реферера из команды
func (h *StartHandler) extractReferrerID(update tgbotapi.Update) int64 {
	text := h.GetText(update)

	// Проверяем, есть ли реферальный код в команде /start
	if len(text) > 6 && text[:6] == "/start" {
		// Ищем реферальный код после /start
		parts := text[6:]
		if len(parts) > 0 {
			// Убираем пробелы и проверяем формат ref_123456
			parts = parts[1:] // убираем пробел
			if len(parts) > 4 && parts[:4] == "ref_" {
				// Извлекаем ID реферера
				var referrerID int64
				if _, err := fmt.Sscanf(parts[4:], "%d", &referrerID); err == nil {
					return referrerID
				}
			}
		}
	}

	return 0
}

// processReferralRegistration обрабатывает регистрацию по реферальной ссылке
func (h *StartHandler) processReferralRegistration(ctx context.Context, userID, referrerID int64, username, firstName, lastName, languageCode string) (*usecase.UserProfile, error) {
	// Регистрируем пользователя по реферальной ссылке
	_, err := h.UseCaseManager.ProcessReferralRegistration(ctx, referrerID, userID, username, firstName, lastName, languageCode)
	if err != nil {
		return nil, fmt.Errorf("failed to process referral registration: %w", err)
	}

	// Получаем профиль пользователя
	profile, err := h.UseCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return profile, nil
}

// processNormalRegistration обрабатывает обычную регистрацию
func (h *StartHandler) processNormalRegistration(ctx context.Context, userID int64, username, firstName, lastName, languageCode string) (*usecase.UserProfile, error) {
	// Регистрируем пользователя
	_, err := h.UseCaseManager.ProcessUserRegistration(ctx, userID, username, firstName, lastName, languageCode)
	if err != nil {
		return nil, fmt.Errorf("failed to process user registration: %w", err)
	}

	// Получаем профиль пользователя
	profile, err := h.UseCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return profile, nil
}

// sendWelcomeMessage отправляет приветственное сообщение
func (h *StartHandler) sendWelcomeMessage(ctx context.Context, chatID int64, user *usecase.UserProfile, isReferral bool) error {
	var message string
	var keyboard tgbotapi.InlineKeyboardMarkup

	if isReferral {
		message = fmt.Sprintf(`
🎉 Добро пожаловать, %s!

🎁 Вы пришли по реферальной ссылке!
⏰ Вам предоставлен расширенный пробный период на 7 дней

🚀 Начните пользоваться VPN прямо сейчас!

📱 Используйте кнопки ниже для управления подпиской:
`, user.User.FirstName)
	} else {
		message = fmt.Sprintf(`
👋 Добро пожаловать, %s!

🎁 Вам предоставлен пробный период на 3 дня
🚀 Начните пользоваться VPN прямо сейчас!

📱 Используйте кнопки ниже для управления подпиской:
`, user.User.FirstName)
	}

	// Создаем клавиатуру
	keyboard = h.createMainKeyboard()

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// createMainKeyboard создает основную клавиатуру
func (h *StartHandler) createMainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)
}

// Command возвращает команду обработчика
func (h *StartHandler) Command() string {
	return "start"
}

// Description возвращает описание обработчика
func (h *StartHandler) Description() string {
	return "Начать работу с ботом"
}
