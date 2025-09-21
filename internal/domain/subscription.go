package domain

import (
	"time"
)

// SubscriptionStatus представляет статус подписки
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
	SubscriptionStatusDisabled SubscriptionStatus = "disabled"
)

// Subscription представляет подписку пользователя
type Subscription struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	PlanID    int64              `json:"plan_id"`
	Status    SubscriptionStatus `json:"status"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	IsTrial   bool               `json:"is_trial"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// IsActive проверяет, активна ли подписка
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.EndDate)
}

// IsExpired проверяет, истекла ли подписка
func (s *Subscription) IsExpired() bool {
	return s.Status == SubscriptionStatusExpired || time.Now().After(s.EndDate)
}

// GetDaysRemaining возвращает количество оставшихся дней
func (s *Subscription) GetDaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	remaining := time.Until(s.EndDate)
	return int(remaining.Hours() / 24)
}
