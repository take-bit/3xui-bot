package core

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int64     `json:"id" db:"id"`
	TelegramID   int64     `json:"telegram_id" db:"telegram_id"`
	Username     string    `json:"username" db:"username"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	LanguageCode string    `json:"language_code" db:"language_code"`
	IsBlocked    bool      `json:"is_blocked" db:"is_blocked"`
	HasTrial     bool      `json:"has_trial" db:"has_trial"` // Использовал ли пробный период
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// GetDisplayName возвращает отображаемое имя пользователя
func (u *User) GetDisplayName() string {
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.Username != "" {
		return "@" + u.Username
	}
	return "Пользователь"
}

// IsActive проверяет, активен ли пользователь
func (u *User) IsActive() bool {
	return !u.IsBlocked
}
