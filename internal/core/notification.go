package core

import (
	"time"
)

// Notification представляет уведомление
type Notification struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationType представляет тип уведомления
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeSuccess NotificationType = "success"
)

// MarkAsRead отмечает уведомление как прочитанное
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// IsUnread проверяет, является ли уведомление непрочитанным
func (n *Notification) IsUnread() bool {
	return !n.IsRead
}

// GetTypeIcon возвращает иконку для типа уведомления
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
