package callback

import (
	"context"
	"log/slog"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"
)

func (h *BaseHandler) HandleGetTrial(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling get trial", "user_id", userID)

	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")

		return err
	}

	success, err := h.activateTrial(ctx, userID)
	if err != nil {
		h.logError(err, "ActivateTrial")

		return err
	}

	var text string
	if success {
		text = "🎉 Пробный доступ активирован на 3 дня!"

		err = h.createTrialSubscription(ctx, userID)
		if err != nil {
			h.logError(err, "CreateTrialSubscription")
		}

		user.HasTrial = true
	} else {
		text = "❌ Пробный доступ уже был использован"
	}

	_ = h.msg.DeleteMessage(ctx, chatID, messageID)
	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)

	return h.msg.SendPhotoWithMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}

func (h *BaseHandler) createTrialSubscription(ctx context.Context, userID int64) error {
	now := time.Now()
	dto := usecase.CreateSubscriptionDTO{
		UserID:    userID,
		Name:      "Пробная подписка",
		PlanID:    "trial",
		Days:      3,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, 3),
		IsActive:  true,
	}

	subscription, err := h.createSubscription(ctx, dto)
	if err != nil {

		return err
	}

	if sub, ok := subscription.(*core.Subscription); ok && sub != nil {
		_, err = h.createVPNForSubscription(ctx, userID, sub.ID)
		if err != nil {
			slog.Error("Failed to create VPN for trial subscription", "error", err)
		}
	}

	return nil
}
