package domain

import (
	"time"
)

// Referral представляет реферальную связь
type Referral struct {
	ID         int64      `json:"id"`
	ReferrerID int64      `json:"referrer_id"` // кто пригласил
	ReferredID int64      `json:"referred_id"` // кого пригласили
	Level      int        `json:"level"`       // уровень реферала (1 или 2)
	RewardDays int        `json:"reward_days"` // награда в днях
	IsPaid     bool       `json:"is_paid"`     // выплачена ли награда
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
}

// ReferralStats представляет статистику рефералов
type ReferralStats struct {
	UserID            int64 `json:"user_id"`
	TotalReferrals    int   `json:"total_referrals"`
	Level1Referrals   int   `json:"level1_referrals"`
	Level2Referrals   int   `json:"level2_referrals"`
	TotalRewardDays   int   `json:"total_reward_days"`
	PaidRewardDays    int   `json:"paid_reward_days"`
	PendingRewardDays int   `json:"pending_reward_days"`
}
