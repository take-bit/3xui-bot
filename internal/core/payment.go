package core

import (
	"time"
)

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

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

func (p *Payment) IsPending() bool {

	return p.Status == string(PaymentStatusPending)
}

func (p *Payment) IsCompleted() bool {

	return p.Status == string(PaymentStatusCompleted)
}

func (p *Payment) IsFailed() bool {

	return p.Status == string(PaymentStatusFailed)
}

func (p *Payment) IsCancelled() bool {

	return p.Status == string(PaymentStatusCancelled)
}
