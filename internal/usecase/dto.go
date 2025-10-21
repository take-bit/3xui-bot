package usecase

import (
	"3xui-bot/internal/core"
	"time"
)

type CreateUserDTO struct {
	TelegramID   int64
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
}

type CreateSubscriptionDTO struct {
	UserID    int64
	Name      string
	PlanID    string
	Days      int
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

type CreatePaymentDTO struct {
	UserID        int64
	Amount        float64
	Currency      string
	PaymentMethod string
	Description   string
}

type CreateConfigDTO struct {
	UserID     int64
	Name       string
	ConfigType string
}

type SendNotificationDTO struct {
	UserID  int64
	Type    core.NotificationType
	Title   string
	Message string
}

type CreateNotificationDTO struct {
	UserID  int64
	Type    string
	Title   string
	Message string
}
