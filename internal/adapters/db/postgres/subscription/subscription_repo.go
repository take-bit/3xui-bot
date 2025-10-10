package subscription

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
)

type Subscription struct {
	dbGetter transactorPgx.DBGetter
}

func NewSubscription(dbGetter transactorPgx.DBGetter) *Subscription {
	return &Subscription{
		dbGetter: dbGetter,
	}
}

func (s *Subscription) CreateSubscription(ctx context.Context, subscription *core.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, name, plan_id, start_date, end_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := s.dbGetter(ctx).Exec(ctx, query,
		subscription.ID, subscription.UserID, subscription.Name, subscription.PlanID,
		subscription.StartDate, subscription.EndDate, subscription.IsActive,
		subscription.CreatedAt, subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (s *Subscription) GetSubscriptionByID(ctx context.Context, id string) (*core.Subscription, error) {
	query := `
		SELECT id, user_id, name, plan_id, start_date, end_date, is_active, created_at, updated_at
		FROM subscriptions WHERE id = $1`

	subscription := &core.Subscription{}
	err := s.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Name, &subscription.PlanID,
		&subscription.StartDate, &subscription.EndDate, &subscription.IsActive,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return subscription, nil
}

func (s *Subscription) GetSubscriptionsByUserID(ctx context.Context, userID int64) ([]*core.Subscription, error) {
	query := `
		SELECT id, user_id, name, plan_id, start_date, end_date, is_active, created_at, updated_at
		FROM subscriptions WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := s.dbGetter(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by user ID: %w", err)
	}
	defer rows.Close()

	var subscriptions []*core.Subscription
	for rows.Next() {
		subscription := &core.Subscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.Name, &subscription.PlanID,
			&subscription.StartDate, &subscription.EndDate, &subscription.IsActive,
			&subscription.CreatedAt, &subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (s *Subscription) GetActiveSubscriptionByUserID(ctx context.Context, userID int64) (*core.Subscription, error) {
	query := `
		SELECT id, user_id, name, plan_id, start_date, end_date, is_active, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND is_active = true AND end_date > NOW()
		ORDER BY created_at DESC
		LIMIT 1`

	subscription := &core.Subscription{}
	err := s.dbGetter(ctx).QueryRow(ctx, query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Name, &subscription.PlanID,
		&subscription.StartDate, &subscription.EndDate, &subscription.IsActive,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return subscription, nil
}

func (s *Subscription) UpdateSubscription(ctx context.Context, subscription *core.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET name = $2, plan_id = $3, start_date = $4, end_date = $5, 
		    is_active = $6, updated_at = $7
		WHERE id = $1`

	result, err := s.dbGetter(ctx).Exec(ctx, query,
		subscription.ID, subscription.Name, subscription.PlanID,
		subscription.StartDate, subscription.EndDate, subscription.IsActive,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (s *Subscription) DeleteSubscription(ctx context.Context, id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := s.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// Plan repository
type Plan struct {
	dbGetter transactorPgx.DBGetter
}

func NewPlan(dbGetter transactorPgx.DBGetter) *Plan {
	return &Plan{
		dbGetter: dbGetter,
	}
}

func (p *Plan) GetAll(ctx context.Context) ([]*core.Plan, error) {
	query := `
		SELECT id, name, description, price, days, is_active
		FROM plans WHERE is_active = true
		ORDER BY days ASC`

	rows, err := p.dbGetter(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all plans: %w", err)
	}
	defer rows.Close()

	var plans []*core.Plan
	for rows.Next() {
		plan := &core.Plan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.Price,
			&plan.Days, &plan.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating plans: %w", err)
	}

	return plans, nil
}

func (p *Plan) GetPlanByID(ctx context.Context, id string) (*core.Plan, error) {
	query := `
		SELECT id, name, description, price, days, is_active
		FROM plans WHERE id = $1`

	plan := &core.Plan{}
	err := p.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.Price,
		&plan.Days, &plan.IsActive,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return plan, nil
}
