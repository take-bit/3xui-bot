package domain

import (
	"context"
	"time"
)

// UserService определяет интерфейс для бизнес-логики пользователей
type UserService interface {
	CreateOrUpdate(ctx context.Context, telegramID int64, username, firstName, lastName, languageCode string) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	Block(ctx context.Context, telegramID int64) error
	Unblock(ctx context.Context, telegramID int64) error
	GetDisplayName(ctx context.Context, telegramID int64) (string, error)
}

// SubscriptionService определяет интерфейс для бизнес-логики подписок
type SubscriptionService interface {
	CreateTrial(ctx context.Context, userID int64, days int) (*Subscription, error)
	CreateFromPayment(ctx context.Context, userID, planID int64, days int) (*Subscription, error)
	Extend(ctx context.Context, userID int64, days int) error
	GetActive(ctx context.Context, userID int64) (*Subscription, error)
	Expire(ctx context.Context, userID int64) error
	CheckExpired(ctx context.Context) error
	GetDaysRemaining(ctx context.Context, userID int64) (int, error)
}

// PaymentService определяет интерфейс для бизнес-логики платежей
type PaymentService interface {
	Create(ctx context.Context, userID, planID int64, amount int, currency string, method PaymentMethod) (*Payment, error)
	ProcessWebhook(ctx context.Context, externalID string, status PaymentStatus) error
	GetByExternalID(ctx context.Context, externalID string) (*Payment, error)
	Complete(ctx context.Context, paymentID int64) error
	Fail(ctx context.Context, paymentID int64) error
	GetPending(ctx context.Context) ([]*Payment, error)
}

// PromocodeService определяет интерфейс для бизнес-логики промокодов
type PromocodeService interface {
	Create(ctx context.Context, code string, promocodeType PromocodeType, value int, usageLimit int, expiresAt *time.Time) (*Promocode, error)
	Validate(ctx context.Context, code string) (*Promocode, error)
	Use(ctx context.Context, code string, userID int64) error
	GetActive(ctx context.Context) ([]*Promocode, error)
	Deactivate(ctx context.Context, code string) error
}

// ReferralService определяет интерфейс для бизнес-логики рефералов
type ReferralService interface {
	CreateReferral(ctx context.Context, referrerID, referredID int64) error
	ProcessPaymentReward(ctx context.Context, userID int64, amount int) error
	GetStats(ctx context.Context, userID int64) (*ReferralStats, error)
	GetUnpaidRewards(ctx context.Context) ([]*Referral, error)
	PayReward(ctx context.Context, referralID int64) error
}

// ServerService определяет интерфейс для бизнес-логики серверов
type ServerService interface {
	AddClient(ctx context.Context, userID int64) (*Server, error)
	RemoveClient(ctx context.Context, userID int64) error
	GetAvailableServer(ctx context.Context) (*Server, error)
	UpdateStatus(ctx context.Context, serverID int64, status ServerStatus) error
	GetStats(ctx context.Context) ([]*Server, error)
}

// NotificationService определяет интерфейс для бизнес-логики уведомлений
type NotificationService interface {
	SendToUser(ctx context.Context, userID int64, title, message string, isHTML bool) error
	SendToAll(ctx context.Context, title, message string, isHTML bool) error
	CreateDraft(ctx context.Context, notificationType NotificationType, userID *int64, title, message string, isHTML bool) (*Notification, error)
	UpdateDraft(ctx context.Context, id int64, title, message string, isHTML bool) error
	SendDraft(ctx context.Context, id int64) error
	GetUserNotifications(ctx context.Context, userID int64) ([]*Notification, error)
}

// BotService определяет основной интерфейс бота
type BotService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	HandleMessage(ctx context.Context, userID int64, message string) error
	HandleCallback(ctx context.Context, userID int64, callbackData string) error
	HandlePayment(ctx context.Context, userID int64, paymentID int64) error
}
