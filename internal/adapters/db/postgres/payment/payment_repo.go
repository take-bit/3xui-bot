package payment

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
)

type Payment struct {
	dbGetter transactorPgx.DBGetter
}

func NewPayment(dbGetter transactorPgx.DBGetter) *Payment {
	return &Payment{
		dbGetter: dbGetter,
	}
}

func (p *Payment) CreatePayment(ctx context.Context, payment *core.Payment) error {
	query := `
		INSERT INTO payments (id, user_id, amount, currency, payment_method, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := p.dbGetter(ctx).Exec(ctx, query,
		payment.ID, payment.UserID, payment.Amount, payment.Currency,
		payment.PaymentMethod, payment.Description, payment.Status,
		payment.CreatedAt, payment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (p *Payment) GetPaymentByID(ctx context.Context, id string) (*core.Payment, error) {
	query := `
		SELECT id, user_id, amount, currency, payment_method, description, status, created_at, updated_at
		FROM payments WHERE id = $1`

	payment := &core.Payment{}
	err := p.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&payment.ID, &payment.UserID, &payment.Amount, &payment.Currency,
		&payment.PaymentMethod, &payment.Description, &payment.Status,
		&payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return payment, nil
}

func (p *Payment) GetPaymentsByUserID(ctx context.Context, userID int64) ([]*core.Payment, error) {
	query := `
		SELECT id, user_id, amount, currency, payment_method, description, status, created_at, updated_at
		FROM payments WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := p.dbGetter(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user ID: %w", err)
	}
	defer rows.Close()

	var payments []*core.Payment
	for rows.Next() {
		payment := &core.Payment{}
		err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Amount, &payment.Currency,
			&payment.PaymentMethod, &payment.Description, &payment.Status,
			&payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	return payments, nil
}

func (p *Payment) UpdatePayment(ctx context.Context, payment *core.Payment) error {
	query := `
		UPDATE payments 
		SET amount = $2, currency = $3, payment_method = $4, 
		    description = $5, status = $6, updated_at = $7
		WHERE id = $1`

	result, err := p.dbGetter(ctx).Exec(ctx, query,
		payment.ID, payment.Amount, payment.Currency, payment.PaymentMethod,
		payment.Description, payment.Status, payment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (p *Payment) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	query := `UPDATE payments SET status = $2, updated_at = NOW() WHERE id = $1`

	result, err := p.dbGetter(ctx).Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (p *Payment) DeletePayment(ctx context.Context, id string) error {
	query := `DELETE FROM payments WHERE id = $1`

	_, err := p.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	return nil
}
