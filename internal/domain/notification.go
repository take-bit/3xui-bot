package domain

import (
	"time"
)

// NotificationType представляет тип уведомления
type NotificationType string

const (
	NotificationTypeUser   NotificationType = "user"   // конкретному пользователю
	NotificationTypeAll    NotificationType = "all"    // всем пользователям
	NotificationTypeSystem NotificationType = "system" // системное уведомление
)

// NotificationStatus представляет статус уведомления
type NotificationStatus string

const (
	NotificationStatusDraft   NotificationStatus = "draft"
	NotificationStatusSending NotificationStatus = "sending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// Notification представляет уведомление
type Notification struct {
	ID        int64              `json:"id"`
	Type      NotificationType   `json:"type"`
	Status    NotificationStatus `json:"status"`
	UserID    *int64             `json:"user_id,omitempty"` // для уведомлений конкретному пользователю
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	IsHTML    bool               `json:"is_html"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
}

// IsEditable проверяет, можно ли редактировать уведомление
func (n *Notification) IsEditable() bool {
	return n.Status == NotificationStatusDraft
}

// IsSent проверяет, отправлено ли уведомление
func (n *Notification) IsSent() bool {
	return n.Status == NotificationStatusSent
}
