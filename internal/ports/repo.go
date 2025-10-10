package ports

import (
	"context"

	"3xui-bot/internal/core"
)

// UserRepo интерфейс для работы с пользователями
type UserRepo interface {
	CreateUser(ctx context.Context, user *core.User) error
	GetUserByID(ctx context.Context, id int64) (*core.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*core.User, error)
	UpdateUser(ctx context.Context, user *core.User) error
	MarkTrialAsUsed(ctx context.Context, userID int64) error
}

// SubscriptionRepo интерфейс для работы с подписками
type SubscriptionRepo interface {
	CreateSubscription(ctx context.Context, subscription *core.Subscription) error
	GetSubscriptionByID(ctx context.Context, id string) (*core.Subscription, error)
	GetSubscriptionsByUserID(ctx context.Context, userID int64) ([]*core.Subscription, error)
	GetActiveSubscriptionByUserID(ctx context.Context, userID int64) (*core.Subscription, error)
	UpdateSubscription(ctx context.Context, subscription *core.Subscription) error
	DeleteSubscription(ctx context.Context, id string) error
}

// PlanRepo интерфейс для работы с планами подписки
type PlanRepo interface {
	GetPlanByID(ctx context.Context, id string) (*core.Plan, error)
	GetAll(ctx context.Context) ([]*core.Plan, error)
}

// PaymentRepo интерфейс для работы с платежами
type PaymentRepo interface {
	CreatePayment(ctx context.Context, payment *core.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*core.Payment, error)
	GetPaymentsByUserID(ctx context.Context, userID int64) ([]*core.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id, status string) error
	DeletePayment(ctx context.Context, id string) error
}

// ReferralRepo интерфейс для работы с рефералами
type ReferralRepo interface {
	CreateReferral(ctx context.Context, referral *core.Referral) error
	GetReferralByID(ctx context.Context, id int64) (*core.Referral, error)
	GetReferralsByReferrerID(ctx context.Context, referrerID int64) ([]*core.Referral, error)
	GetReferralByRefereeID(ctx context.Context, refereeID int64) (*core.Referral, error)
}

// ReferralLinkRepo интерфейс для работы с реферальными ссылками
type ReferralLinkRepo interface {
	CreateReferralLink(ctx context.Context, link *core.ReferralLink) error
	GetReferralLinkByID(ctx context.Context, id int64) (*core.ReferralLink, error)
	GetReferralLinkByUserID(ctx context.Context, userID int64) (*core.ReferralLink, error)
	GetReferralLinkByLink(ctx context.Context, link string) (*core.ReferralLink, error)
	UpdateReferralLink(ctx context.Context, link *core.ReferralLink) error
	DeleteReferralLink(ctx context.Context, id int64) error
}

// VPNRepo интерфейс для работы с VPN подключениями
type VPNRepo interface {
	CreateVPNConnection(ctx context.Context, conn *core.VPNConnection) error
	GetVPNConnectionsByTelegramUserID(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error)
	GetVPNConnectionByID(ctx context.Context, id string) (*core.VPNConnection, error)
	GetVPNConnectionByMarzbanUsername(ctx context.Context, marzbanUsername string) (*core.VPNConnection, error)
	UpdateVPNConnectionName(ctx context.Context, id, name string) error
	DeleteVPNConnection(ctx context.Context, id string) error
	DeleteVPNConnectionByMarzbanUsername(ctx context.Context, marzbanUsername string) error
	GetActiveVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error)
}

// NotificationRepo интерфейс для работы с уведомлениями
type NotificationRepo interface {
	CreateNotification(ctx context.Context, notification *core.Notification) error
	GetNotificationByID(ctx context.Context, id string) (*core.Notification, error)
	GetNotificationsByUserID(ctx context.Context, userID int64) ([]*core.Notification, error)
	GetUnreadNotificationsByUserID(ctx context.Context, userID int64) ([]*core.Notification, error)
	UpdateNotification(ctx context.Context, notification *core.Notification) error
	MarkAsRead(ctx context.Context, id string) error
	DeleteNotification(ctx context.Context, id string) error
}
