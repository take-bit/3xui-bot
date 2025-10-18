package usecase

import (
	"context"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
	"3xui-bot/internal/ports"
)

// SubscriptionUseCase use case для работы с подписками
type SubscriptionUseCase struct {
	subRepo  ports.SubscriptionRepo
	planRepo ports.PlanRepo
}

// NewSubscriptionUseCase создает новый use case для подписок
func NewSubscriptionUseCase(subRepo ports.SubscriptionRepo, planRepo ports.PlanRepo) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subRepo:  subRepo,
		planRepo: planRepo,
	}
}

// CreateSubscription создает новую подписку
func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (*core.Subscription, error) {
	// Создаем новую подписку
	newSub := &core.Subscription{
		ID:        id.Generate(), // Генерируем уникальный ID
		UserID:    dto.UserID,
		Name:      dto.Name,
		PlanID:    dto.PlanID,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
		IsActive:  dto.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.subRepo.CreateSubscription(ctx, newSub)
	if err != nil {
		return nil, err
	}

	return newSub, nil
}

// GetUserSubscriptions получает все подписки пользователя
func (uc *SubscriptionUseCase) GetUserSubscriptions(ctx context.Context, userID int64) ([]*core.Subscription, error) {
	return uc.subRepo.GetSubscriptionsByUserID(ctx, userID)
}

// GetSubscription получает подписку по ID
func (uc *SubscriptionUseCase) GetSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.Subscription, error) {
	return uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
}

// GetSubscriptionByID получает подписку по ID (алиас для удобства)
func (uc *SubscriptionUseCase) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*core.Subscription, error) {
	return uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
}

// GetActiveSubscription получает активную подписку пользователя
func (uc *SubscriptionUseCase) GetActiveSubscription(ctx context.Context, userID int64) (*core.Subscription, error) {
	return uc.subRepo.GetActiveSubscriptionByUserID(ctx, userID)
}

// UpdateSubscriptionName обновляет название подписки
func (uc *SubscriptionUseCase) UpdateSubscriptionName(ctx context.Context, userID int64, subscriptionID, name string) error {
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if sub.UserID != userID {
		return ErrUnauthorized
	}

	sub.Name = name
	sub.UpdatedAt = time.Now()

	return uc.subRepo.UpdateSubscription(ctx, sub)
}

// ExtendSubscription продлевает подписку
func (uc *SubscriptionUseCase) ExtendSubscription(ctx context.Context, userID int64, subscriptionID string, days int) error {
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if sub.UserID != userID {
		return ErrUnauthorized
	}

	// Продлеваем подписку
	if sub.EndDate.After(time.Now()) {
		// Если подписка еще активна, продлеваем от текущей даты окончания
		sub.EndDate = sub.EndDate.AddDate(0, 0, days)
	} else {
		// Если подписка истекла, продлеваем от текущей даты
		sub.EndDate = time.Now().AddDate(0, 0, days)
	}

	sub.IsActive = true
	sub.UpdatedAt = time.Now()

	return uc.subRepo.UpdateSubscription(ctx, sub)
}

// CancelSubscription отменяет подписку
func (uc *SubscriptionUseCase) CancelSubscription(ctx context.Context, userID int64, subscriptionID string) error {
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if sub.UserID != userID {
		return ErrUnauthorized
	}

	sub.IsActive = false
	sub.UpdatedAt = time.Now()

	return uc.subRepo.UpdateSubscription(ctx, sub)
}

// DeleteSubscription удаляет подписку
func (uc *SubscriptionUseCase) DeleteSubscription(ctx context.Context, userID int64, subscriptionID string) error {
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if sub.UserID != userID {
		return ErrUnauthorized
	}

	return uc.subRepo.DeleteSubscription(ctx, subscriptionID)
}

// GetPlans получает все доступные планы
func (uc *SubscriptionUseCase) GetPlans(ctx context.Context) ([]*core.Plan, error) {
	return uc.planRepo.GetAll(ctx)
}

// GetPlan получает план по ID
func (uc *SubscriptionUseCase) GetPlan(ctx context.Context, planID string) (*core.Plan, error) {
	return uc.planRepo.GetPlanByID(ctx, planID)
}

// GetPlanByID получает план по ID (алиас для удобства)
func (uc *SubscriptionUseCase) GetPlanByID(ctx context.Context, planID string) (*core.Plan, error) {
	return uc.planRepo.GetPlanByID(ctx, planID)
}
