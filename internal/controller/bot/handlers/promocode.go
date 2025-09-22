package handlers

import (
	"context"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PromocodeHandler обрабатывает команду /promocode
type PromocodeHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewPromocodeHandler создает новый обработчик команды /promocode
func NewPromocodeHandler(useCaseManager *usecase.UseCaseManager) *PromocodeHandler {
	return &PromocodeHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /promocode
func (h *PromocodeHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	chatID := h.GetChatID(update)

	// Отправляем информацию о промокодах
	message := `
🎁 <b>Промокоды</b>

💡 <b>Как использовать промокод:</b>
1️⃣ Введите промокод в поле ниже
2️⃣ Нажмите кнопку "Применить"
3️⃣ Получите бонусы на свой аккаунт

<b>Типы промокодов:</b>
⏰ <b>Дополнительные дни</b> - продлевают подписку
💰 <b>Скидка</b> - уменьшают стоимость подписки

<b>Примеры промокодов:</b>
• <code>WELCOME10</code> - 10 дополнительных дней
• <code>DISCOUNT20</code> - 20%% скидка на подписку
• <code>NEWUSER</code> - 7 дополнительных дней

<b>Введите промокод:</b>
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *PromocodeHandler) Command() string {
	return "promocode"
}

// Description возвращает описание обработчика
func (h *PromocodeHandler) Description() string {
	return "Применить промокод"
}
