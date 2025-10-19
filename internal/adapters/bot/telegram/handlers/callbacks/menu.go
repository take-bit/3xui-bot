package callbacks

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
)

// HandleOpenMenu handles the open_menu callback
func (h *BaseHandler) HandleOpenMenu(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open menu", "user_id", userID)

	// Получаем пользователя
	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}

	// Проверяем активные подписки
	subscriptions, _ := h.getUserSubscriptions(ctx, userID)
	isPremium := false
	statusText := "🆓 Бесплатный"
	subUntilText := ""

	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if !sub.IsExpired() {
				isPremium = true
				statusText = "⭐ Premium"
				subUntilText = sub.EndDate.Format("02.01.2006")
				break
			}
		}
	}

	text := ui.GetMainMenuWithProfileText(user, isPremium, statusText, subUntilText)
	keyboard := ui.GetMainMenuWithProfileKeyboard(isPremium)
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleOpenProfile handles the open_profile callback
func (h *BaseHandler) HandleOpenProfile(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open profile", "user_id", userID)

	// Получаем пользователя
	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}

	// Check if user is premium
	subscriptions, _ := h.getUserSubscriptions(ctx, userID)
	isPremium := false
	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if !sub.IsExpired() {
				isPremium = true
				break
			}
		}
	}

	text := ui.GetProfileText(user, isPremium, "", "")
	keyboard := ui.GetProfileKeyboard(isPremium)
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleOpenPricing handles the open_pricing callback
func (h *BaseHandler) HandleOpenPricing(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open pricing", "user_id", userID)

	// Получаем планы
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)
	// Используем deleteAndSendMessage, т.к. главное меню может быть с фото
	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

// HandleOpenSupport handles the open_support callback
func (h *BaseHandler) HandleOpenSupport(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open support", "user_id", userID)

	text := ui.GetSupportText()
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}
