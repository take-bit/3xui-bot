package callback

import (
	"context"
	"log/slog"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"
)

// HandleGetTrial handles the get_trial callback
func (h *BaseHandler) HandleGetTrial(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling get trial", "user_id", userID)

	// Получаем пользователя
	user, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}

	// Активируем пробный доступ
	success, err := h.activateTrial(ctx, userID)
	if err != nil {
		h.logError(err, "ActivateTrial")
		return err
	}

	var text string
	if success {
		text = "🎉 Пробный доступ активирован на 3 дня!"

		// Создаем пробную подписку
		err = h.createTrialSubscription(ctx, userID)
		if err != nil {
			h.logError(err, "CreateTrialSubscription")
			// Не возвращаем ошибку, т.к. пробный доступ уже активирован
		}

		// Обновляем user.HasTrial для правильного отображения клавиатуры
		user.HasTrial = true
	} else {
		text = "❌ Пробный доступ уже был использован"
	}

	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
	// Используем deleteAndSendMessage, т.к. приветственное сообщение может быть с фото
	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

// createTrialSubscription создает пробную подписку на 3 дня
func (h *BaseHandler) createTrialSubscription(ctx context.Context, userID int64) error {
	now := time.Now()
	dto := usecase.CreateSubscriptionDTO{
		UserID:    userID,
		Name:      "Пробная подписка",
		PlanID:    "trial", // Специальный ID для пробной подписки
		Days:      3,       // 3 дня пробного доступа
		StartDate: now,
		EndDate:   now.AddDate(0, 0, 3),
		IsActive:  true,
	}

	subscription, err := h.createSubscription(ctx, dto)
	if err != nil {
		return err
	}

	// Создаем VPN для пробной подписки
	if sub, ok := subscription.(*core.Subscription); ok && sub != nil {
		_, err = h.createVPNForSubscription(ctx, userID, sub.ID)
		if err != nil {
			// Логируем ошибку, но не возвращаем её
			slog.Error("Failed to create VPN for trial subscription", "error", err)
		}
	}

	return nil
}
