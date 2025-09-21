package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

// subscriptionRepository реализует domain.SubscriptionRepository
type subscriptionRepository struct {
	dbGetter DBGetter
}

// NewSubscriptionRepository создает новый репозиторий подписок
func NewSubscriptionRepository(dbGetter DBGetter) domain.SubscriptionRepository {
	return &subscriptionRepository{
		dbGetter: dbGetter,
	}
}

// Create создает новую подписку
func (r *subscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, status, start_date, end_date, is_trial, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	db := r.dbGetter.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		subscription.UserID,
		subscription.PlanID,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.IsTrial,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).Scan(&subscription.ID)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// GetByID получает подписку по ID
func (r *subscriptionRepository) GetByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, start_date, end_date, is_trial, created_at, updated_at
		FROM subscriptions
		WHERE id = $1`

	db := r.dbGetter.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	subscription, err := r.scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription by id: %w", err)
	}

	return subscription, nil
}

// GetByUserID получает подписку пользователя
func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, start_date, end_date, is_trial, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	db := r.dbGetter.GetDB(ctx)
	row := db.QueryRow(ctx, query, userID)

	subscription, err := r.scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get subscription by user id: %w", err)
	}

	return subscription, nil
}

// GetActiveByUserID получает активную подписку пользователя
func (r *subscriptionRepository) GetActiveByUserID(ctx context.Context, userID int64) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, start_date, end_date, is_trial, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1 AND status = $2 AND end_date > $3
		ORDER BY created_at DESC
		LIMIT 1`

	db := r.dbGetter.GetDB(ctx)
	row := db.QueryRow(ctx, query, userID, domain.SubscriptionStatusActive, time.Now())

	subscription, err := r.scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to get active subscription by user id: %w", err)
	}

	return subscription, nil
}

// Update обновляет подписку
func (r *subscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET plan_id = $2, status = $3, start_date = $4, end_date = $5, is_trial = $6, updated_at = $7
		WHERE id = $1`

	db := r.dbGetter.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		subscription.ID,
		subscription.PlanID,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.IsTrial,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// Delete удаляет подписку
func (r *subscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	db := r.dbGetter.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// ListExpired получает список истекших подписок
func (r *subscriptionRepository) ListExpired(ctx context.Context) ([]*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, start_date, end_date, is_trial, created_at, updated_at
		FROM subscriptions
		WHERE status = $1 AND end_date <= $2
		ORDER BY end_date ASC`

	db := r.dbGetter.GetDB(ctx)
	rows, err := db.Query(ctx, query, domain.SubscriptionStatusActive, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to list expired subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*domain.Subscription
	for rows.Next() {
		subscription, err := r.scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subscriptions: %w", err)
	}

	return subscriptions, nil
}

// Extend продлевает подписку на указанное количество дней
func (r *subscriptionRepository) Extend(ctx context.Context, id int64, days int) error {
	query := `
		UPDATE subscriptions
		SET end_date = end_date + INTERVAL '%d days', updated_at = $2
		WHERE id = $1`

	db := r.dbGetter.GetDB(ctx)
	result, err := db.Exec(ctx, fmt.Sprintf(query, days), id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to extend subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

// scanSubscription сканирует подписку из строки результата
func (r *subscriptionRepository) scanSubscription(row pgx.Row) (*domain.Subscription, error) {
	var subscription domain.Subscription

	err := row.Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PlanID,
		&subscription.Status,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.IsTrial,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}
