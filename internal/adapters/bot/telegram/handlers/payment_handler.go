package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PaymentHandler struct {
	bot       *tgbotapi.BotAPI
	paymentUC *usecase.PaymentUseCase
}

func NewPaymentHandler(
	bot *tgbotapi.BotAPI,
	paymentUC *usecase.PaymentUseCase,
) *PaymentHandler {

	return &PaymentHandler{
		bot:       bot,
		paymentUC: paymentUC,
	}
}

func (h *PaymentHandler) HandleSelectPlan(ctx context.Context, userID int64, chatID int64, planID string) error {
	slog.Info("User selected plan", "user_id", userID, "plan_id", planID)

	payment, paymentURL, err := h.paymentUC.CreatePaymentForPlan(ctx, userID, planID)
	if err != nil {

		return fmt.Errorf("failed to create payment: %w", err)
	}

	message := fmt.Sprintf(
		"💳 *Оплата подписки*\n\n"+
			"Сумма: %.2f ₽\n"+
			"ID платежа: %s\n\n"+
			"⚠️ Для оплаты используйте кнопку ниже.\n"+
			"После успешной оплаты подписка будет активирована автоматически.",
		payment.Amount,
		payment.ID,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("💳 Оплатить", paymentURL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Я оплатил", fmt.Sprintf("payment_check_%s", payment.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", fmt.Sprintf("payment_cancel_%s", payment.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {

		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (h *PaymentHandler) HandlePaymentCheck(ctx context.Context, userID int64, chatID int64, messageID int, paymentID string, planID string) error {
	slog.Info("Checking payment %s for user %d", paymentID, userID)

	if err := h.paymentUC.ProcessPaymentSuccess(ctx, paymentID, planID); err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Платеж не найден или еще не обработан. Попробуйте позже.")
		h.bot.Send(msg)

		return fmt.Errorf("failed to process payment: %w", err)
	}

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	h.bot.Send(deleteMsg)

	successMsg := tgbotapi.NewMessage(chatID,
		"✅ *Платеж успешно обработан!*\n\n"+
			"Ваша подписка активирована.\n"+
			"VPN подключение создано.\n\n"+
			"Используйте /vpn для получения конфигурации.",
	)
	successMsg.ParseMode = "Markdown"

	if _, err := h.bot.Send(successMsg); err != nil {

		return fmt.Errorf("failed to send success message: %w", err)
	}

	return nil
}

func (h *PaymentHandler) HandlePaymentCancel(ctx context.Context, userID int64, chatID int64, messageID int, paymentID string) error {
	slog.Info("Cancelling payment %s for user %d", paymentID, userID)

	if err := h.paymentUC.ProcessPaymentCancellation(ctx, paymentID); err != nil {

		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	h.bot.Send(deleteMsg)

	msg := tgbotapi.NewMessage(chatID, "❌ Платеж отменен.")
	if _, err := h.bot.Send(msg); err != nil {

		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (h *PaymentHandler) HandlePaymentWebhook(ctx context.Context, paymentID string, planID string, status string) error {
	slog.Info("Received webhook for payment %s with status %s", paymentID, status)

	switch status {
	case "succeeded", "completed":

		return h.paymentUC.ProcessPaymentSuccess(ctx, paymentID, planID)
	case "failed":

		return h.paymentUC.ProcessPaymentFailure(ctx, paymentID)
	case "cancelled", "canceled":

		return h.paymentUC.ProcessPaymentCancellation(ctx, paymentID)
	default:
		slog.Info("Unknown payment status", "status", status)

		return nil
	}
}
