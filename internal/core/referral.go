package core

import (
	"time"
)

// Referral представляет реферальную связь
type Referral struct {
	ID         int64     `json:"id"`
	ReferrerID int64     `json:"referrer_id"`
	RefereeID  int64     `json:"referee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// ReferralLink представляет реферальную ссылку
type ReferralLink struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Link      string    `json:"link"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsExpired проверяет, истекла ли реферальная ссылка (опционально)
func (r *ReferralLink) IsExpired() bool {
	// TODO: Добавить логику истечения ссылок при необходимости
	return false
}

// GetShortLink возвращает укороченную версию ссылки для отображения
func (r *ReferralLink) GetShortLink() string {
	if len(r.Link) > 50 {
		return r.Link[:47] + "..."
	}
	return r.Link
}
