package core

import (
	"time"
)

// Subscription представляет подписку пользователя
type Subscription struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"` // Название подписки (например, "Основная", "Работа", "Дом")
	PlanID    string    `json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsExpired проверяет, истекла ли подписка
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}

// DaysRemaining возвращает количество дней до истечения
func (s *Subscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	return int(time.Until(s.EndDate).Hours() / 24)
}

// GetDisplayName возвращает отображаемое название подписки
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

// GetStatusText возвращает текстовое описание статуса подписки
func (s *Subscription) GetStatusText() string {
	if !s.IsActive {
		return "Неактивна"
	}
	if s.IsExpired() {
		return "Истекла"
	}
	return "Активна"
}

// Extend продлевает подписку на указанное количество дней
func (s *Subscription) Extend(days int) {
	s.EndDate = s.EndDate.AddDate(0, 0, days)
	s.UpdatedAt = time.Now()
}

// GetStatus возвращает статус подписки
func (s *Subscription) GetStatus() SubscriptionStatus {
	if !s.IsActive {
		return StatusInactive
	}
	if s.IsExpired() {
		return StatusExpired
	}
	return StatusActive
}

// SubscriptionStatus представляет статус подписки
type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusExpired  SubscriptionStatus = "expired"
	StatusInactive SubscriptionStatus = "inactive"
)

// SubscriptionPeriod представляет период подписки
type SubscriptionPeriod struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// IsValid проверяет валидность периода
func (p *SubscriptionPeriod) IsValid() bool {
	return p.StartDate.Before(p.EndDate)
}

// Duration возвращает продолжительность периода в днях
func (p *SubscriptionPeriod) Duration() int {
	return int(p.EndDate.Sub(p.StartDate).Hours() / 24)
}

// Plan представляет тарифный план
type Plan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Days        int     `json:"days"`
	IsActive    bool    `json:"is_active"`
}

// GetPricePerDay возвращает стоимость за день
func (p *Plan) GetPricePerDay() float64 {
	if p.Days == 0 {
		return 0
	}
	return p.Price / float64(p.Days)
}

// GetDiscount возвращает скидку в процентах по сравнению с месячным планом
func (p *Plan) GetDiscount() float64 {
	monthlyPrice := 5.0 // Базовая цена за месяц
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
