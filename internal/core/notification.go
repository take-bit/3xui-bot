package core

import (
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeSuccess NotificationType = "success"
)

func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

func (n *Notification) IsUnread() bool {

	return !n.IsRead
}

func (n *Notification) GetTypeIcon() string {
	switch NotificationType(n.Type) {
	case NotificationTypeInfo:

		return "ℹ️"
	case NotificationTypeWarning:

		return "⚠️"
	case NotificationTypeError:

		return "❌"
	case NotificationTypeSuccess:

		return "✅"
	default:

		return "📌"
	}
}
