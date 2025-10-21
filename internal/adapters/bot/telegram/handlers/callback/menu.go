package callback

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
)

func (h *BaseHandler) HandleOpenMenu(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open menu", "user_id", userID)
	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")

		return err
	}
	subscriptions, err := h.getUserSubscriptions(ctx, userID)
	if err != nil {
		subscriptions = []*core.Subscription{}
	}
	isPremium := false
	for _, sub := range subscriptions {
		if sub.IsActive && !sub.IsExpired() {
			isPremium = true
			break
		}
	}
	_ = h.msg.DeleteMessage(ctx, chatID, messageID)
	text := ui.GetMainMenuWithProfileText(user, subscriptions)
	keyboard := ui.GetMainMenuWithProfileKeyboard(isPremium)

	return h.msg.SendPhotoWithPreEscapedMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}

func (h *BaseHandler) HandleOpenProfile(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open profile", "user_id", userID)
	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")

		return err
	}
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

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleOpenPricing(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open pricing", "user_id", userID)
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")

		return err
	}
	_ = h.msg.DeleteMessage(ctx, chatID, messageID)
	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)

	return h.msg.SendPhotoWithMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}

func (h *BaseHandler) HandleShowInstruction(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling show instruction", "user_id", userID)
	text := ui.GetInstructionText()
	keyboard := ui.GetBackToMenuKeyboard()
	_ = h.msg.DeleteMessage(ctx, chatID, messageID)

	return h.msg.SendPhotoWithMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}
