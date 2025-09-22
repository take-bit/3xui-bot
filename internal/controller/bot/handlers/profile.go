package handlers

import (
	"context"
	"fmt"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ProfileHandler обрабатывает команду /profile
type ProfileHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewProfileHandler создает новый обработчик команды /profile
func NewProfileHandler(useCaseManager *usecase.UseCaseManager) *ProfileHandler {
	return &ProfileHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /profile
func (h *ProfileHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)

	// Получаем профиль пользователя
	profile, err := h.UseCaseManager.GetUserUseCase().GetUserProfile(ctx, userID)
	if err != nil {
		h.HandleError(ctx, chatID, err, "❌ Ошибка при получении профиля. Попробуйте позже.")
		return err
	}

	// Формируем сообщение с информацией о профиле
	message := h.formatProfileMessage(profile)

	// Создаем клавиатуру
	keyboard := h.createProfileKeyboard()

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// formatProfileMessage форматирует сообщение с информацией о профиле
func (h *ProfileHandler) formatProfileMessage(profile *usecase.UserProfile) string {
	message := fmt.Sprintf(`
👤 <b>Мой профиль</b>

<b>Основная информация:</b>
🆔 ID: <code>%d</code>
👤 Имя: %s
📧 Username: @%s
🌐 Язык: %s
📅 Регистрация: %s

<b>Подписка:</b>
`, profile.User.ID, profile.User.FirstName, profile.User.Username, profile.User.LanguageCode, profile.RegistrationDate)

	if profile.Subscription != nil {
		message += fmt.Sprintf(`
✅ Статус: Активна
📅 Действует до: %s
⏰ Осталось дней: %d
`, profile.Subscription.EndDate.Format("02.01.2006 15:04"), profile.DaysRemaining)
	} else {
		message += `
❌ Статус: Неактивна
💳 Продлите подписку для продолжения использования
`
	}

	// Добавляем информацию о рефералах
	if profile.ReferralStats != nil {
		message += fmt.Sprintf(`
<b>Реферальная программа:</b>
👥 Всего рефералов: %d
💰 Заработано дней: %d
`, profile.ReferralStats.TotalReferrals, profile.ReferralStats.TotalRewardDays)
	}

	return message
}

// createProfileKeyboard создает клавиатуру для профиля
func (h *ProfileHandler) createProfileKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)
}

// Command возвращает команду обработчика
func (h *ProfileHandler) Command() string {
	return "profile"
}

// Description возвращает описание обработчика
func (h *ProfileHandler) Description() string {
	return "Показать мой профиль"
}
