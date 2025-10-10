package core

import (
	"time"
)

// Payment представляет платеж
type Payment struct {
	ID            string    `json:"id"`
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentStatus представляет статус платежа
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// IsPending проверяет, находится ли платеж в ожидании
func (p *Payment) IsPending() bool {
	return p.Status == string(PaymentStatusPending)
}

// IsCompleted проверяет, завершен ли платеж
func (p *Payment) IsCompleted() bool {
	return p.Status == string(PaymentStatusCompleted)
}

// IsFailed проверяет, провален ли платеж
func (p *Payment) IsFailed() bool {
	return p.Status == string(PaymentStatusFailed)
}

// IsCancelled проверяет, отменен ли платеж
func (p *Payment) IsCancelled() bool {
	return p.Status == string(PaymentStatusCancelled)
}
