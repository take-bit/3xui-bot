package usecase

import (
	"3xui-bot/internal/core"
	"time"
)

// CreateUserDTO данные для создания пользователя
type CreateUserDTO struct {
	TelegramID   int64
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
}

// CreateSubscriptionDTO данные для создания подписки
type CreateSubscriptionDTO struct {
	UserID    int64
	Name      string
	PlanID    string
	Days      int
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

// CreatePaymentDTO данные для создания платежа
type CreatePaymentDTO struct {
	UserID        int64
	Amount        float64
	Currency      string
	PaymentMethod string
	Description   string
}

// CreateConfigDTO данные для создания VPN конфигурации
type CreateConfigDTO struct {
	UserID     int64
	Name       string
	ConfigType string
}

// SendNotificationDTO данные для отправки уведомления
type SendNotificationDTO struct {
	UserID  int64
	Type    core.NotificationType
	Title   string
	Message string
}

// CreateNotificationDTO данные для создания уведомления
type CreateNotificationDTO struct {
	UserID  int64
	Type    string
	Title   string
	Message string
}
