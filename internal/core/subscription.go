package core

import (
	"time"
)

type Subscription struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	PlanID    string    `json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Subscription) IsExpired() bool {

	return time.Now().After(s.EndDate)
}

func (s *Subscription) DaysRemaining() int {
	if s.IsExpired() {

		return 0
	}

	return int(time.Until(s.EndDate).Hours() / 24)
}

func (s *Subscription) GetDisplayName() string {
	if s.Name != "" {

		return s.Name
	}
	if len(s.ID) >= 8 {

		return "Подписка " + s.ID[:8]
	}
	if s.ID != "" {

		return "Подписка " + s.ID
	}

	return "Подписка"
}

func (s *Subscription) GetStatusText() string {
	if !s.IsActive {

		return "Неактивна"
	}
	if s.IsExpired() {

		return "Истекла"
	}

	return "Активна"
}

func (s *Subscription) Extend(days int) {
	s.EndDate = s.EndDate.AddDate(0, 0, days)
	s.UpdatedAt = time.Now()
}

func (s *Subscription) GetStatus() SubscriptionStatus {
	if !s.IsActive {

		return StatusInactive
	}
	if s.IsExpired() {

		return StatusExpired
	}

	return StatusActive
}

type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusExpired  SubscriptionStatus = "expired"
	StatusInactive SubscriptionStatus = "inactive"
)

type SubscriptionPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (p *SubscriptionPeriod) IsValid() bool {

	return p.StartDate.Before(p.EndDate)
}

func (p *SubscriptionPeriod) Duration() int {

	return int(p.EndDate.Sub(p.StartDate).Hours() / 24)
}

type Plan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Days        int     `json:"days"`
	IsActive    bool    `json:"is_active"`
}

func (p *Plan) GetPricePerDay() float64 {
	if p.Days == 0 {

		return 0
	}

	return p.Price / float64(p.Days)
}

func (p *Plan) GetDiscount() float64 {
	monthlyPrice := 5.0
	monthlyDays := 30

	if p.Days <= monthlyDays {

		return 0
	}

	expectedPrice := (monthlyPrice / float64(monthlyDays)) * float64(p.Days)
	if expectedPrice == 0 {

		return 0
	}

	discount := ((expectedPrice - p.Price) / expectedPrice) * 100
	if discount < 0 {

		return 0
	}

	return discount
}
