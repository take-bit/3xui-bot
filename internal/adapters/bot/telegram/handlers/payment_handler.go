package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PaymentHandler обработчик платежей
type PaymentHandler struct {
	bot       *tgbotapi.BotAPI
	paymentUC *usecase.PaymentUseCase
}

// NewPaymentHandler создает новый обработчик платежей
func NewPaymentHandler(
	bot *tgbotapi.BotAPI,
	paymentUC *usecase.PaymentUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		bot:       bot,
		paymentUC: paymentUC,
	}
}

// HandleSelectPlan обрабатывает выбор плана подписки
func (h *PaymentHandler) HandleSelectPlan(ctx context.Context, userID int64, chatID int64, planID string) error {
	slog.Info("User selected plan", "user_id", userID, "plan_id", planID)

	// Создаем платеж через UseCase
	payment, paymentURL, err := h.paymentUC.CreatePaymentForPlan(ctx, userID, planID)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// Отправляем сообщение с ссылкой на оплату
	message := fmt.Sprintf(
		"💳 *Оплата подписки*\n\n"+
			"Сумма: %.2f ₽\n"+
			"ID платежа: %s\n\n"+
			"⚠️ Для оплаты используйте кнопку ниже.\n"+
			"После успешной оплаты подписка будет активирована автоматически.",
		payment.Amount,
		payment.ID,
	)

	// Создаем клавиатуру с кнопкой оплаты
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

// HandlePaymentCheck обрабатывает проверку статуса платежа
func (h *PaymentHandler) HandlePaymentCheck(ctx context.Context, userID int64, chatID int64, messageID int, paymentID string, planID string) error {
	slog.Info("Checking payment %s for user %d", paymentID, userID)

	// Обрабатываем успешный платеж через UseCase (вся бизнес-логика внутри)
	if err := h.paymentUC.ProcessPaymentSuccess(ctx, paymentID, planID); err != nil {
		// Отправляем сообщение об ошибке
		msg := tgbotapi.NewMessage(chatID, "❌ Платеж не найден или еще не обработан. Попробуйте позже.")
		h.bot.Send(msg)
		return fmt.Errorf("failed to process payment: %w", err)
	}

	// Удаляем предыдущее сообщение
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	h.bot.Send(deleteMsg)

	// Отправляем успешное сообщение
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

// HandlePaymentCancel обрабатывает отмену платежа
func (h *PaymentHandler) HandlePaymentCancel(ctx context.Context, userID int64, chatID int64, messageID int, paymentID string) error {
	slog.Info("Cancelling payment %s for user %d", paymentID, userID)

	// Отменяем платеж через UseCase
	if err := h.paymentUC.ProcessPaymentCancellation(ctx, paymentID); err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	// Удаляем сообщение с платежом
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	h.bot.Send(deleteMsg)

	// Отправляем подтверждение отмены
	msg := tgbotapi.NewMessage(chatID, "❌ Платеж отменен.")
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// HandlePaymentWebhook обрабатывает webhook от платежной системы
// TODO: Реализовать HTTP endpoint для webhook
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
