package usecase

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/domain"
)

// PaymentUseCase представляет use case для работы с платежами
type PaymentUseCase struct {
	paymentService      domain.PaymentService
	subscriptionService domain.SubscriptionService
	vpnUseCase          *VPNUseCase
	userService         domain.UserService
	notificationService domain.NotificationService
}

// NewPaymentUseCase создает новый Payment use case
func NewPaymentUseCase(
	paymentService domain.PaymentService,
	subscriptionService domain.SubscriptionService,
	vpnUseCase *VPNUseCase,
	userService domain.UserService,
	notificationService domain.NotificationService,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentService:      paymentService,
		subscriptionService: subscriptionService,
		vpnUseCase:          vpnUseCase,
		userService:         userService,
		notificationService: notificationService,
	}
}

// CreatePayment создает новый платеж
func (uc *PaymentUseCase) CreatePayment(ctx context.Context, userID int64, planID int64, amount int, currency string, method domain.PaymentMethod) (*domain.Payment, error) {
	// 1. Проверяем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Создаем платеж
	payment, err := uc.paymentService.Create(ctx, user.ID, planID, amount, currency, method)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 3. Отправляем уведомление о создании платежа
	if uc.notificationService != nil {
		message := fmt.Sprintf("💳 Платеж создан!\n\n💰 Сумма: %d %s\n📋 Метод: %s\n🆔 ID: %d",
			amount, currency, method, payment.ID)
		_ = uc.notificationService.SendToUser(ctx, userID, "Платеж", message, false)
	}

	return payment, nil
}

// ProcessPaymentWebhook обрабатывает webhook от платежной системы
func (uc *PaymentUseCase) ProcessPaymentWebhook(ctx context.Context, externalID string, status domain.PaymentStatus) error {
	// 1. Получаем платеж по внешнему ID
	payment, err := uc.paymentService.GetByExternalID(ctx, externalID)
	if err != nil {
		return fmt.Errorf("failed to get payment by external ID: %w", err)
	}

	// 2. Обрабатываем webhook
	err = uc.paymentService.ProcessWebhook(ctx, externalID, status)
	if err != nil {
		return fmt.Errorf("failed to process webhook: %w", err)
	}

	// 3. Если платеж успешен, создаем подписку и VPN
	if status == domain.PaymentStatusCompleted {
		err = uc.handleSuccessfulPayment(ctx, payment)
		if err != nil {
			return fmt.Errorf("failed to handle successful payment: %w", err)
		}
	}

	// 4. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		user, err := uc.userService.GetByTelegramID(ctx, payment.UserID)
		if err == nil {
			var message string
			switch status {
			case domain.PaymentStatusCompleted:
				message = "✅ Платеж успешно обработан!\n\n🎉 Ваша подписка активирована!"
			case domain.PaymentStatusFailed:
				message = "❌ Платеж не прошел\n\nПопробуйте еще раз или обратитесь в поддержку"
			case domain.PaymentStatusCancelled:
				message = "🚫 Платеж отменен"
			default:
				message = fmt.Sprintf("📋 Статус платежа: %s", status)
			}
			_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Платеж", message, false)
		}
	}

	return nil
}

// CompletePayment завершает платеж вручную
func (uc *PaymentUseCase) CompletePayment(ctx context.Context, paymentID int64) error {
	// 1. Завершаем платеж
	err := uc.paymentService.Complete(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to complete payment: %w", err)
	}

	// 2. Получаем обновленный платеж
	// TODO: Добавить метод GetByID в PaymentService
	// payment, err := uc.paymentService.GetByID(ctx, paymentID)
	// if err != nil {
	//     return fmt.Errorf("failed to get payment: %w", err)
	// }

	// 3. Обрабатываем успешный платеж
	// err = uc.handleSuccessfulPayment(ctx, payment)
	// if err != nil {
	//     return fmt.Errorf("failed to handle successful payment: %w", err)
	// }

	return nil
}

// FailPayment помечает платеж как неудачный
func (uc *PaymentUseCase) FailPayment(ctx context.Context, paymentID int64) error {
	// 1. Помечаем платеж как неудачный
	err := uc.paymentService.Fail(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to fail payment: %w", err)
	}

	// 2. Отправляем уведомление пользователю
	// TODO: Получить userID из платежа и отправить уведомление

	return nil
}

// GetPendingPayments возвращает список ожидающих платежей
func (uc *PaymentUseCase) GetPendingPayments(ctx context.Context) ([]*domain.Payment, error) {
	return uc.paymentService.GetPending(ctx)
}

// GetUserPayments возвращает платежи пользователя
func (uc *PaymentUseCase) GetUserPayments(ctx context.Context, userID int64) ([]*domain.Payment, error) {
	// 1. Получаем пользователя
	_, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем платежи пользователя
	// TODO: Добавить метод GetByUserID в PaymentService
	// return uc.paymentService.GetByUserID(ctx, user.ID)
	return nil, fmt.Errorf("GetByUserID not implemented yet")
}

// handleSuccessfulPayment обрабатывает успешный платеж
func (uc *PaymentUseCase) handleSuccessfulPayment(ctx context.Context, payment *domain.Payment) error {
	// 1. Создаем подписку из платежа
	// TODO: Добавить поле Days в Payment или получать из Plan
	days := 30 // Временное значение
	subscription, err := uc.subscriptionService.CreateFromPayment(ctx, payment.UserID, payment.PlanID, days)
	if err != nil {
		return fmt.Errorf("failed to create subscription from payment: %w", err)
	}

	// 2. Создаем VPN подключение
	_, err = uc.vpnUseCase.CreateVPNConnection(ctx, payment.UserID, "")
	if err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("Failed to create VPN connection after payment: %v\n", err)
	}

	// 3. Отправляем уведомление о создании подписки
	if uc.notificationService != nil {
		user, err := uc.userService.GetByTelegramID(ctx, payment.UserID)
		if err == nil {
			message := fmt.Sprintf("🎉 Подписка активирована!\n\n📅 Действует до: %s\n🔗 VPN подключение создано",
				subscription.EndDate.Format("02.01.2006 15:04"))
			_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Подписка", message, false)
		}
	}

	return nil
}

// ProcessRefund обрабатывает возврат средств
func (uc *PaymentUseCase) ProcessRefund(ctx context.Context, paymentID int64, reason string) error {
	// 1. Получаем платеж
	// TODO: Добавить метод GetByID в PaymentService
	// payment, err := uc.paymentService.GetByID(ctx, paymentID)
	// if err != nil {
	//     return fmt.Errorf("failed to get payment: %w", err)
	// }

	// 2. Проверяем, что платеж можно вернуть
	// if payment.Status != domain.PaymentStatusCompleted {
	//     return fmt.Errorf("payment is not completed, cannot refund")
	// }

	// 3. Создаем возврат
	// TODO: Реализовать логику возврата

	// 4. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		// user, err := uc.userService.GetByTelegramID(ctx, payment.UserID)
		// if err == nil {
		//     message := fmt.Sprintf("💰 Возврат обработан\n\n💸 Сумма: %d %s\n📋 Причина: %s",
		//         payment.Amount, payment.Currency, reason)
		//     _ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Возврат", message, false)
		// }
	}

	return fmt.Errorf("refund processing not implemented yet")
}

// GetPaymentStats возвращает статистику платежей
func (uc *PaymentUseCase) GetPaymentStats(ctx context.Context, from, to time.Time) (*PaymentStats, error) {
	// TODO: Реализовать получение статистики платежей
	return nil, fmt.Errorf("payment stats not implemented yet")
}

// PaymentStats представляет статистику платежей
type PaymentStats struct {
	TotalPayments      int     `json:"total_payments"`
	SuccessfulPayments int     `json:"successful_payments"`
	FailedPayments     int     `json:"failed_payments"`
	TotalAmount        int     `json:"total_amount"`
	AverageAmount      float64 `json:"average_amount"`
	SuccessRate        float64 `json:"success_rate"`
}
