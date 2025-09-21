package domain

import (
	"time"
)

// PromocodeType представляет тип промокода
type PromocodeType string

const (
	PromocodeTypeExtraDays PromocodeType = "extra_days"
	PromocodeTypeDiscount  PromocodeType = "discount"
)

// Promocode представляет промокод
type Promocode struct {
	ID         int64         `json:"id"`
	Code       string        `json:"code"`
	Type       PromocodeType `json:"type"`
	Value      int           `json:"value"` // количество дней или процент скидки
	IsActive   bool          `json:"is_active"`
	UsageLimit int           `json:"usage_limit"` // 0 = без ограничений
	UsedCount  int           `json:"used_count"`
	ExpiresAt  *time.Time    `json:"expires_at,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// IsValid проверяет, действителен ли промокод
func (p *Promocode) IsValid() bool {
	if !p.IsActive {
		return false
	}

	if p.UsageLimit > 0 && p.UsedCount >= p.UsageLimit {
		return false
	}

	if p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt) {
		return false
	}

	return true
}

// CanBeUsed проверяет, можно ли использовать промокод
func (p *Promocode) CanBeUsed() bool {
	return p.IsValid() && (p.UsageLimit == 0 || p.UsedCount < p.UsageLimit)
}
