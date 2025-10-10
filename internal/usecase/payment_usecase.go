package usecase

import (
	"3xui-bot/internal/ports"
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
)

// PaymentProvider интерфейс для платежного провайдера (внешний сервис)
type PaymentProvider interface {
	CreatePayment(ctx context.Context, amount float64, currency, description string) (paymentURL string, paymentID string, err error)
	CheckPaymentStatus(ctx context.Context, paymentID string) (status string, err error)
}

// PaymentUseCase use case для работы с платежами
type PaymentUseCase struct {
	paymentRepo    ports.PaymentRepo
	subscriptionUC *SubscriptionUseCase
	vpnUC          *VPNUseCase
	notifUC        *NotificationUseCase
	provider       PaymentProvider
}

// NewPaymentUseCase создает новый use case для платежей
func NewPaymentUseCase(
	paymentRepo ports.PaymentRepo,
	subscriptionUC *SubscriptionUseCase,
	vpnUC *VPNUseCase,
	notifUC *NotificationUseCase,
	provider PaymentProvider,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:    paymentRepo,
		subscriptionUC: subscriptionUC,
		vpnUC:          vpnUC,
		notifUC:        notifUC,
		provider:       provider,
	}
}

// CreatePayment создает новый платеж
func (uc *PaymentUseCase) CreatePayment(ctx context.Context, dto CreatePaymentDTO) (*core.Payment, error) {
	newPayment := &core.Payment{
		UserID:        dto.UserID,
		Amount:        dto.Amount,
		Currency:      dto.Currency,
		PaymentMethod: dto.PaymentMethod,
		Description:   dto.Description,
		Status:        string(core.PaymentStatusPending),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := uc.paymentRepo.CreatePayment(ctx, newPayment)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}

// GetPayment получает платеж по ID
func (uc *PaymentUseCase) GetPayment(ctx context.Context, paymentID string) (*core.Payment, error) {
	return uc.paymentRepo.GetPaymentByID(ctx, paymentID)
}

// GetUserPayments получает все платежи пользователя
func (uc *PaymentUseCase) GetUserPayments(ctx context.Context, userID int64) ([]*core.Payment, error) {
	return uc.paymentRepo.GetPaymentsByUserID(ctx, userID)
}

// CompletePayment завершает платеж
func (uc *PaymentUseCase) CompletePayment(ctx context.Context, paymentID string) error {
	return uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusCompleted))
}

// FailPayment отмечает платеж как неудачный
func (uc *PaymentUseCase) FailPayment(ctx context.Context, paymentID string) error {
	return uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusFailed))
}

// CancelPayment отменяет платеж
func (uc *PaymentUseCase) CancelPayment(ctx context.Context, paymentID string) error {
	return uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusCancelled))
}

// CreatePaymentForPlan создает платеж для выбранного плана (бизнес-логика)
func (uc *PaymentUseCase) CreatePaymentForPlan(ctx context.Context, userID int64, planID string) (*core.Payment, string, error) {
	// Получаем план
	plan, err := uc.subscriptionUC.GetPlan(ctx, planID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get plan: %w", err)
	}

	// Создаем платеж в БД
	payment := &core.Payment{
		ID:            id.Generate(),
		UserID:        userID,
		Amount:        plan.Price,
		Currency:      "RUB",
		PaymentMethod: "mock", // В реальности будет yookassa/stripe
		Description:   fmt.Sprintf("Подписка: %s", plan.Name),
		Status:        string(core.PaymentStatusPending),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, "", fmt.Errorf("failed to create payment: %w", err)
	}

	// Создаем платеж в провайдере
	paymentURL, externalID, err := uc.provider.CreatePayment(
		ctx,
		plan.Price,
		"RUB",
		payment.Description,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create payment in provider: %w", err)
	}

	// Сохраняем external ID (в реальности нужно добавить поле в Payment)
	_ = externalID

	return payment, paymentURL, nil
}

// ProcessPaymentSuccess обрабатывает успешную оплату (оркестрация)
func (uc *PaymentUseCase) ProcessPaymentSuccess(ctx context.Context, paymentID string, planID string) error {
	// Получаем платеж
	payment, err := uc.paymentRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Проверяем что платеж еще не обработан
	if payment.IsCompleted() {
		return fmt.Errorf("payment already completed")
	}

	// Обновляем статус платежа
	if err := uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusCompleted)); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Получаем план
	plan, err := uc.subscriptionUC.GetPlan(ctx, planID)
	if err != nil {
		return fmt.Errorf("failed to get plan: %w", err)
	}

	// Создаем подписку через SubscriptionUseCase
	subscriptionDTO := CreateSubscriptionDTO{
		UserID:    payment.UserID,
		Name:      "Основная подписка",
		PlanID:    planID,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, plan.Days),
		IsActive:  true,
	}

	subscription, err := uc.subscriptionUC.CreateSubscription(ctx, subscriptionDTO)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Создаем VPN подключение через VPNUseCase
	vpnConn, err := uc.vpnUC.CreateVPNForSubscription(ctx, payment.UserID, subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to create VPN: %w", err)
	}

	// Отправляем уведомление через NotificationUseCase
	notifDTO := CreateNotificationDTO{
		UserID:  payment.UserID,
		Type:    "payment",
		Title:   "✅ Платеж успешен",
		Message: fmt.Sprintf("Ваш платеж на сумму %.2f ₽ успешно обработан. VPN подключение \"%s\" активировано!", payment.Amount, vpnConn.Name),
	}

	if err := uc.notifUC.CreateNotification(ctx, notifDTO); err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("failed to send notification: %v\n", err)
	}

	return nil
}

// ProcessPaymentFailure обрабатывает неудачную оплату
func (uc *PaymentUseCase) ProcessPaymentFailure(ctx context.Context, paymentID string) error {
	return uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusFailed))
}

// ProcessPaymentCancellation обрабатывает отмену платежа
func (uc *PaymentUseCase) ProcessPaymentCancellation(ctx context.Context, paymentID string) error {
	return uc.paymentRepo.UpdatePaymentStatus(ctx, paymentID, string(core.PaymentStatusCancelled))
}
