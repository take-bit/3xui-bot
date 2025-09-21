package domain

import (
	"time"
)

// PaymentStatus представляет статус платежа
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod представляет способ оплаты
type PaymentMethod string

const (
	PaymentMethodCryptomus     PaymentMethod = "cryptomus"
	PaymentMethodHeleket       PaymentMethod = "heleket"
	PaymentMethodYooKassa      PaymentMethod = "yookassa"
	PaymentMethodYooMoney      PaymentMethod = "yoomoney"
	PaymentMethodTelegramStars PaymentMethod = "telegram_stars"
)

// Payment представляет платеж
type Payment struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id"`
	PlanID      int64         `json:"plan_id"`
	Amount      int           `json:"amount"`
	Currency    string        `json:"currency"`
	Method      PaymentMethod `json:"method"`
	Status      PaymentStatus `json:"status"`
	ExternalID  string        `json:"external_id"` // ID платежа во внешней системе
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
}

// IsCompleted проверяет, завершен ли платеж
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsPending проверяет, ожидает ли платеж обработки
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}
