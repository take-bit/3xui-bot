package handlers

import (
	"context"
	"fmt"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SubscriptionHandler обрабатывает команду /subscription
type SubscriptionHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewSubscriptionHandler создает новый обработчик команды /subscription
func NewSubscriptionHandler(useCaseManager *usecase.UseCaseManager) *SubscriptionHandler {
	return &SubscriptionHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /subscription
func (h *SubscriptionHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)

	// Получаем информацию о подписке
	subscriptionInfo, err := h.UseCaseManager.GetSubscriptionUseCase().GetActiveSubscription(ctx, userID)
	if err != nil {
		// Если активной подписки нет, показываем варианты продления
		return h.handleNoActiveSubscription(ctx, chatID, userID)
	}

	// Отправляем информацию о подписке
	return h.sendSubscriptionInfo(ctx, chatID, subscriptionInfo)
}

// handleNoActiveSubscription обрабатывает случай, когда нет активной подписки
func (h *SubscriptionHandler) handleNoActiveSubscription(ctx context.Context, chatID, userID int64) error {
	message := `
💳 <b>Подписка</b>

❌ У вас нет активной подписки

💡 <b>Выберите тарифный план:</b>

<b>📅 1 месяц</b>
💰 Цена: 100 руб.
⏰ Срок: 30 дней

<b>📅 3 месяца</b>
💰 Цена: 250 руб.
⏰ Срок: 90 дней
💸 <b>Экономия: 50 руб.</b>

<b>📅 6 месяцев</b>
💰 Цена: 450 руб.
⏰ Срок: 180 дней
💸 <b>Экономия: 150 руб.</b>

<b>📅 1 год</b>
💰 Цена: 800 руб.
⏰ Срок: 365 дней
💸 <b>Экономия: 400 руб.</b>
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 1 месяц - 100₽", "subscription_1_month"),
			tgbotapi.NewInlineKeyboardButtonData("📅 3 месяца - 250₽", "subscription_3_months"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 6 месяцев - 450₽", "subscription_6_months"),
			tgbotapi.NewInlineKeyboardButtonData("📅 1 год - 800₽", "subscription_1_year"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// sendSubscriptionInfo отправляет информацию о подписке
func (h *SubscriptionHandler) sendSubscriptionInfo(ctx context.Context, chatID int64, subscriptionInfo *usecase.SubscriptionInfo) error {
	message := fmt.Sprintf(`
💳 <b>Подписка</b>

✅ <b>Статус:</b> Активна
📅 <b>Действует до:</b> %s
⏰ <b>Осталось дней:</b> %d

<b>VPN подключение:</b>
`,
		subscriptionInfo.ExpiresAt.Format("02.01.2006 15:04"),
		subscriptionInfo.DaysRemaining,
	)

	if subscriptionInfo.VPNConnection != nil {
		message += fmt.Sprintf(`
✅ <b>Статус:</b> Активно
🔗 <b>Конфигурация:</b> <code>%s</code>
`, subscriptionInfo.VPNConnection.ConfigURL)
	} else {
		message += `
❌ <b>Статус:</b> Неактивно
💡 Создайте VPN подключение в разделе "VPN"
`
	}

	// Добавляем информацию о продлении
	if subscriptionInfo.DaysRemaining <= 7 {
		message += `
⚠️ <b>Внимание!</b> Подписка скоро истекает!
💳 Продлите подписку, чтобы не потерять доступ к VPN
`
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Продлить подписку", "subscription_extend"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
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
func (h *SubscriptionHandler) Command() string {
	return "subscription"
}

// Description возвращает описание обработчика
func (h *SubscriptionHandler) Description() string {
	return "Управление подпиской"
}
