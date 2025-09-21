package domain

import (
	"context"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Block(ctx context.Context, id int64) error
	Unblock(ctx context.Context, id int64) error
}

// SubscriptionRepository определяет интерфейс для работы с подписками
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	GetByID(ctx context.Context, id int64) (*Subscription, error)
	GetByUserID(ctx context.Context, userID int64) (*Subscription, error)
	GetActiveByUserID(ctx context.Context, userID int64) (*Subscription, error)
	Update(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, id int64) error
	ListExpired(ctx context.Context) ([]*Subscription, error)
	Extend(ctx context.Context, id int64, days int) error
}

// PlanRepository определяет интерфейс для работы с планами
type PlanRepository interface {
	Create(ctx context.Context, plan *Plan) error
	GetByID(ctx context.Context, id int64) (*Plan, error)
	GetActive(ctx context.Context) ([]*Plan, error)
	Update(ctx context.Context, plan *Plan) error
	Delete(ctx context.Context, id int64) error
	SetActive(ctx context.Context, id int64, active bool) error
}

// ServerRepository определяет интерфейс для работы с серверами
type ServerRepository interface {
	Create(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id int64) (*Server, error)
	GetAvailable(ctx context.Context) ([]*Server, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id int64) error
	SetStatus(ctx context.Context, id int64, status ServerStatus) error
	IncrementClients(ctx context.Context, id int64) error
	DecrementClients(ctx context.Context, id int64) error
}

// PaymentRepository определяет интерфейс для работы с платежами
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id int64) (*Payment, error)
	GetByExternalID(ctx context.Context, externalID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	GetPending(ctx context.Context) ([]*Payment, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Payment, error)
	Complete(ctx context.Context, id int64) error
	Fail(ctx context.Context, id int64) error
}

// PromocodeRepository определяет интерфейс для работы с промокодами
type PromocodeRepository interface {
	Create(ctx context.Context, promocode *Promocode) error
	GetByID(ctx context.Context, id int64) (*Promocode, error)
	GetByCode(ctx context.Context, code string) (*Promocode, error)
	Update(ctx context.Context, promocode *Promocode) error
	Delete(ctx context.Context, id int64) error
	IncrementUsage(ctx context.Context, id int64) error
	GetActive(ctx context.Context) ([]*Promocode, error)
}

// ReferralRepository определяет интерфейс для работы с рефералами
type ReferralRepository interface {
	Create(ctx context.Context, referral *Referral) error
	GetByID(ctx context.Context, id int64) (*Referral, error)
	GetByReferredID(ctx context.Context, referredID int64) (*Referral, error)
	GetByReferrerID(ctx context.Context, referrerID int64) ([]*Referral, error)
	GetStats(ctx context.Context, userID int64) (*ReferralStats, error)
	MarkAsPaid(ctx context.Context, id int64) error
	GetUnpaidRewards(ctx context.Context) ([]*Referral, error)
}

// NotificationRepository определяет интерфейс для работы с уведомлениями
type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id int64) (*Notification, error)
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id int64) error
	GetDraft(ctx context.Context) ([]*Notification, error)
	GetByUserID(ctx context.Context, userID int64) ([]*Notification, error)
	MarkAsSent(ctx context.Context, id int64) error
}
