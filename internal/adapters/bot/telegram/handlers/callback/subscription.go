package callback

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"
)

func (h *BaseHandler) HandleMySubscriptions(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my subscriptions", "user_id", userID)

	subscriptions, err := h.getUserSubscriptions(ctx, userID)
	if err != nil {
		h.logError(err, "GetUserSubscriptions")

		return err
	}

	text := ui.GetSubscriptionsText(subscriptions)
	keyboard := ui.GetSubscriptionsKeyboard(subscriptions)
	_ = h.msg.DeleteMessage(ctx, chatID, messageID)

	return h.msg.SendPhotoWithPreEscapedMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}

func (h *BaseHandler) HandleCreateSubscription(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create subscription", "user_id", userID)

	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")

		return err
	}

	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleSelectPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling select plan", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandlePayCard(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay card", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	slog.Info("Creating MOCK card payment (auto-success)", "user_id", userID, "plan_id", planID)

	now := time.Now()

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

	if sub, ok := subscription.(*core.Subscription); ok && sub != nil {
		_, err = h.createVPNForSubscription(ctx, userID, sub.ID)
		if err != nil {
			h.logError(err, "CreateVPN")
		}
	}

	text := fmt.Sprintf("✅ Оплата успешна!\n\n🎉 Подписка '%s' активирована на %d дней", plan.Name, plan.Days)
	keyboard := ui.GetBackToSubscriptionsKeyboard()

	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandlePaySBP(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay SBP", "plan_id", planID, "user_id", userID)

	return h.HandlePayCard(ctx, userID, chatID, messageID, planID)
}

func (h *BaseHandler) HandlePayStars(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay stars", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	text := fmt.Sprintf("⭐ Оплата Telegram Stars\n\n💰 Сумма: %.0f₽\n⏰ План: %s (%d дней)\n\n🚧 Функция в разработке", plan.Price, plan.Name, plan.Days)
	keyboard := ui.GetBackToPricingKeyboard()

	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleViewSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling view subscription", "subscription_id", subscriptionID, "user_id", userID)

	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")

		return err
	}

	plan, err := h.getPlan(ctx, subscription.PlanID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	vpnConfigs, err := h.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		h.logError(err, "GetVPNConnections")
		vpnConfigs = []*core.VPNConnection{}
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)

	return h.msg.DeleteAndSendMessageWithMarkdownV2(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleRenameSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling rename subscription", "subscription_id", subscriptionID, "user_id", userID)

	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")

		return err
	}

	h.mu.Lock()
	h.renamingUsers[userID] = subscriptionID
	h.mu.Unlock()

	text := ui.GetRenameSubscriptionText(subscription)

	return h.msg.SendMessage(ctx, chatID, text)
}

func (h *BaseHandler) HandleExtendSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling extend subscription", "subscription_id", subscriptionID, "user_id", userID)

	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")

		return err
	}

	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")

		return err
	}

	text := ui.GetExtendSubscriptionText(subscription)
	keyboard := ui.GetExtendSubscriptionKeyboard(subscriptionID, plans)

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleDeleteSubscription(ctx context.Context, userID, chatID int64, messageID int, description string) error {
	slog.Info("Handling delete subscription", "user_id", userID)

	subscription, err := h.getSubscription(ctx, userID, description)
	if err != nil {
		h.logError(err, "GetSubscription")

		return err
	}

	text := ui.GetDeleteSubscriptionText(subscription)
	keyboard := ui.GetBackToPricingKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleExtendSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID, subscriptionID string) error {
	slog.Info("Handling extend subscription by plan", "plan_id", planID, "subscription_id", subscriptionID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	slog.Info("Creating MOCK extension (auto-success)", "user_id", userID, "plan_id", planID, "subscription_id", subscriptionID)

	err = h.extendSubscription(ctx, userID, subscriptionID, plan.Days)
	if err != nil {
		h.logError(err, "ExtendSubscription")

		return h.sendError(chatID, "Ошибка продления подписки")
	}

	text := fmt.Sprintf("✅ Подписка продлена на %d дней!", plan.Days)
	keyboard := ui.GetBackToSubscriptionsKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) generateSubscriptionName(plan *core.Plan) string {
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

	now := time.Now()
	dateStr := now.Format("02.01")

	return fmt.Sprintf("%s (%s)", baseName, dateStr)
}

func (h *BaseHandler) extendSubscription(ctx context.Context, userID int64, subscriptionID string, days int) error {

	return h.subUC.ExtendSubscription(ctx, userID, subscriptionID, days)
}

func (h *BaseHandler) HandleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling create subscription by plan", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")

		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}
