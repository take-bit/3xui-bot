package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	repo *Repository
}

func NewPaymentRepository(repo *Repository) *PaymentRepository {
	return &PaymentRepository{
		repo: repo,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (user_id, plan_id, amount, currency, method, status, external_id, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		payment.UserID,
		payment.PlanID,
		payment.Amount,
		payment.Currency,
		payment.Method,
		payment.Status,
		payment.ExternalID,
		payment.Description,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.CompletedAt,
	).Scan(&payment.ID)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetByID получает платеж по ID
func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `
		SELECT id, user_id, plan_id, amount, currency, method, status, external_id, description, created_at, updated_at, completed_at
		FROM payments
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	payment, err := r.scanPayment(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by id: %w", err)
	}

	return payment, nil
}

// GetByExternalID получает платеж по внешнему ID
func (r *PaymentRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.Payment, error) {
	query := `
		SELECT id, user_id, plan_id, amount, currency, method, status, external_id, description, created_at, updated_at, completed_at
		FROM payments
		WHERE external_id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, externalID)

	payment, err := r.scanPayment(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by external id: %w", err)
	}

	return payment, nil
}

// Update обновляет платеж
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments
		SET user_id = $2, plan_id = $3, amount = $4, currency = $5, method = $6, status = $7, external_id = $8, description = $9, updated_at = $10, completed_at = $11
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		payment.ID,
		payment.UserID,
		payment.PlanID,
		payment.Amount,
		payment.Currency,
		payment.Method,
		payment.Status,
		payment.ExternalID,
		payment.Description,
		payment.UpdatedAt,
		payment.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}

// GetPending получает ожидающие платежи
func (r *PaymentRepository) GetPending(ctx context.Context) ([]*domain.Payment, error) {
	query := `
		SELECT id, user_id, plan_id, amount, currency, method, status, external_id, description, created_at, updated_at, completed_at
		FROM payments
		WHERE status = $1
		ORDER BY created_at ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, domain.PaymentStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending payments: %w", err)
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment, err := r.scanPayment(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate payments: %w", err)
	}

	return payments, nil
}

// GetByUserID получает платежи пользователя
func (r *PaymentRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Payment, error) {
	query := `
		SELECT id, user_id, plan_id, amount, currency, method, status, external_id, description, created_at, updated_at, completed_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user id: %w", err)
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment, err := r.scanPayment(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate payments: %w", err)
	}

	return payments, nil
}

// Complete отмечает платеж как завершенный
func (r *PaymentRepository) Complete(ctx context.Context, id int64) error {
	now := time.Now()
	query := `
		UPDATE payments 
		SET status = $2, completed_at = $3, updated_at = $4 
		WHERE id = $1 AND status = $5`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, domain.PaymentStatusCompleted, &now, now, domain.PaymentStatusPending)
	if err != nil {
		return fmt.Errorf("failed to complete payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPaymentAlreadyPaid
	}

	return nil
}

// Fail отмечает платеж как неудачный
func (r *PaymentRepository) Fail(ctx context.Context, id int64) error {
	query := `
		UPDATE payments 
		SET status = $2, updated_at = $3 
		WHERE id = $1 AND status = $4`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, domain.PaymentStatusFailed, time.Now(), domain.PaymentStatusPending)
	if err != nil {
		return fmt.Errorf("failed to fail payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPaymentAlreadyPaid
	}

	return nil
}

// scanPayment сканирует платеж из строки результата
func (r *PaymentRepository) scanPayment(row pgx.Row) (*domain.Payment, error) {
	var payment domain.Payment
	var completedAt sql.NullTime

	err := row.Scan(
		&payment.ID,
		&payment.UserID,
		&payment.PlanID,
		&payment.Amount,
		&payment.Currency,
		&payment.Method,
		&payment.Status,
		&payment.ExternalID,
		&payment.Description,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		payment.CompletedAt = &completedAt.Time
	}

	return &payment, nil
}
