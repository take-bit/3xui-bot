package usecase

import (
	"context"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
	"3xui-bot/internal/ports"
)

type SubscriptionUseCase struct {
	subRepo  ports.SubscriptionRepo
	planRepo ports.PlanRepo
}

func NewSubscriptionUseCase(subRepo ports.SubscriptionRepo, planRepo ports.PlanRepo) *SubscriptionUseCase {

	return &SubscriptionUseCase{
		subRepo:  subRepo,
		planRepo: planRepo,
	}
}

func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, dto CreateSubscriptionDTO) (*core.Subscription, error) {
	newSub := &core.Subscription{
		ID:        id.Generate(),
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

func (uc *SubscriptionUseCase) GetUserSubscriptions(ctx context.Context, userID int64) ([]*core.Subscription, error) {

	return uc.subRepo.GetSubscriptionsByUserID(ctx, userID)
}

func (uc *SubscriptionUseCase) GetSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.Subscription, error) {

	return uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
}

func (uc *SubscriptionUseCase) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*core.Subscription, error) {

	return uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
}

func (uc *SubscriptionUseCase) GetActiveSubscription(ctx context.Context, userID int64) (*core.Subscription, error) {

	return uc.subRepo.GetActiveSubscriptionByUserID(ctx, userID)
}

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

func (uc *SubscriptionUseCase) ExtendSubscription(ctx context.Context, userID int64, subscriptionID string, days int) error {
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {

		return err
	}

	if sub.UserID != userID {

		return ErrUnauthorized
	}

	if sub.EndDate.After(time.Now()) {
		sub.EndDate = sub.EndDate.AddDate(0, 0, days)
	} else {
		sub.EndDate = time.Now().AddDate(0, 0, days)
	}

	sub.IsActive = true
	sub.UpdatedAt = time.Now()

	return uc.subRepo.UpdateSubscription(ctx, sub)
}

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

func (uc *SubscriptionUseCase) GetPlans(ctx context.Context) ([]*core.Plan, error) {

	return uc.planRepo.GetAll(ctx)
}

func (uc *SubscriptionUseCase) GetPlan(ctx context.Context, planID string) (*core.Plan, error) {

	return uc.planRepo.GetPlanByID(ctx, planID)
}

func (uc *SubscriptionUseCase) GetPlanByID(ctx context.Context, planID string) (*core.Plan, error) {

	return uc.planRepo.GetPlanByID(ctx, planID)
}
