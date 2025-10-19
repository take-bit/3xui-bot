package callbacks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"
)

// HandleMySubscriptions handles the my_subscriptions callback
func (h *BaseHandler) HandleMySubscriptions(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my subscriptions", "user_id", userID)

	// Получаем подписки пользователя
	subscriptions, err := h.getUserSubscriptions(ctx, userID)
	if err != nil {
		h.logError(err, "GetUserSubscriptions")
		return err
	}

	text := ui.GetSubscriptionsText(subscriptions)
	keyboard := ui.GetSubscriptionsKeyboard(subscriptions)
	// Используем deleteAndSendMessageWithMarkdownV2 для корректного форматирования
	return h.msg.DeleteAndSendMessageWithMarkdownV2(ctx, chatID, messageID, text, keyboard)
}

// HandleCreateSubscription handles the create_subscription callback
func (h *BaseHandler) HandleCreateSubscription(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create subscription", "user_id", userID)

	// Получаем планы
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleSelectPlan handles the select_plan callback
func (h *BaseHandler) HandleSelectPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling select plan", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)
	// Используем deleteAndSendMessage, т.к. может быть фото
	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

// HandlePayCard handles the pay_card callback
func (h *BaseHandler) HandlePayCard(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay card", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// MOCK: Картой всегда успешно
	slog.Info("Creating MOCK card payment (auto-success)", "user_id", userID, "plan_id", planID)

	now := time.Now()

	// Генерируем осмысленное название подписки на основе плана
	subscriptionName := h.generateSubscriptionName(plan)

	dto := usecase.CreateSubscriptionDTO{
		UserID:    userID,
		Name:      subscriptionName,
		PlanID:    planID,
		Days:      plan.Days,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.Days),
		IsActive:  true,
	}

	subscription, err := h.createSubscription(ctx, dto)
	if err != nil {
		h.logError(err, "CreateSubscription")
		return h.sendError(chatID, "Ошибка создания подписки")
	}

	// Создаем VPN для подписки
	if sub, ok := subscription.(*core.Subscription); ok && sub != nil {
		_, err = h.createVPNForSubscription(ctx, userID, sub.ID)
		if err != nil {
			h.logError(err, "CreateVPN")
		}
	}

	text := fmt.Sprintf("✅ Оплата успешна!\n\n🎉 Подписка '%s' активирована на %d дней", plan.Name, plan.Days)
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandlePaySBP handles the pay_sbp callback
func (h *BaseHandler) HandlePaySBP(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay SBP", "plan_id", planID, "user_id", userID)
	// SBP использует ту же логику что и карта
	return h.HandlePayCard(ctx, userID, chatID, messageID, planID)
}

// HandlePayStars handles the pay_stars callback
func (h *BaseHandler) HandlePayStars(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay stars", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := fmt.Sprintf("⭐ Оплата Telegram Stars\n\n💰 Сумма: %.0f₽\n⏰ План: %s (%d дней)\n\n🚧 Функция в разработке", plan.Price, plan.Name, plan.Days)
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleViewSubscription handles the view_subscription callback
func (h *BaseHandler) HandleViewSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling view subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Получаем план подписки
	plan, err := h.getPlan(ctx, subscription.PlanID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// Получаем VPN конфигурации подписки
	vpnConfigs, err := h.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		h.logError(err, "GetVPNConnections")
		vpnConfigs = []*core.VPNConnection{} // Пустой массив, если ошибка
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	return h.msg.EditMessageWithMarkdownV2(ctx, chatID, messageID, text, keyboard)
}

// HandleRenameSubscription handles the rename_subscription callback
func (h *BaseHandler) HandleRenameSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling rename subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Сохраняем состояние переименования
	h.mu.Lock()
	h.renamingUsers[userID] = subscriptionID
	h.mu.Unlock()

	text := ui.GetRenameSubscriptionText(subscription)
	return h.msg.SendMessage(ctx, chatID, text)
}

// HandleExtendSubscription handles the extend_subscription callback
func (h *BaseHandler) HandleExtendSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling extend subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Получаем планы для продления
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetExtendSubscriptionText(subscription)
	keyboard := ui.GetExtendSubscriptionKeyboard(subscriptionID, plans)
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleDeleteSubscription handles the delete_subscription callback
func (h *BaseHandler) HandleDeleteSubscription(ctx context.Context, userID, chatID int64, messageID int, description string) error {
	slog.Info("Handling delete subscription", "user_id", userID)

	// Get subscription to pass to the text function
	subscription, err := h.getSubscription(ctx, userID, description)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	text := ui.GetDeleteSubscriptionText(subscription)
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleExtendSubscriptionByPlan handles the extend_subscription_by_plan callback
func (h *BaseHandler) HandleExtendSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID, subscriptionID string) error {
	slog.Info("Handling extend subscription by plan", "plan_id", planID, "subscription_id", subscriptionID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// MOCK: Всегда успешно
	slog.Info("Creating MOCK extension (auto-success)", "user_id", userID, "plan_id", planID, "subscription_id", subscriptionID)

	// Продлеваем подписку
	err = h.extendSubscription(ctx, userID, subscriptionID, plan.Days)
	if err != nil {
		h.logError(err, "ExtendSubscription")
		return h.sendError(chatID, "Ошибка продления подписки")
	}

	text := fmt.Sprintf("✅ Подписка продлена на %d дней!", plan.Days)
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// generateSubscriptionName генерирует осмысленное название подписки на основе плана
func (h *BaseHandler) generateSubscriptionName(plan *core.Plan) string {
	// Базовое название на основе длительности
	var baseName string
	switch {
	case plan.Days <= 7:
		baseName = "Недельная"
	case plan.Days <= 30:
		baseName = "Месячная"
	case plan.Days <= 90:
		baseName = "Квартальная"
	case plan.Days <= 365:
		baseName = "Годовая"
	default:
		baseName = "Долгосрочная"
	}

	// Добавляем дату создания для уникальности
	now := time.Now()
	dateStr := now.Format("02.01")

	return fmt.Sprintf("%s (%s)", baseName, dateStr)
}

// extendSubscription продлевает подписку
func (h *BaseHandler) extendSubscription(ctx context.Context, userID int64, subscriptionID string, days int) error {
	return h.subUC.ExtendSubscription(ctx, userID, subscriptionID, days)
}

// HandleCreateSubscriptionByPlan handles the create_subscription_by_plan callback
func (h *BaseHandler) HandleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling create subscription by plan", "plan_id", planID, "user_id", userID)

	// Get plan
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)
	// Use deleteAndSendMessage as it might be a photo
	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}
