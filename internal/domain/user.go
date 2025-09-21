package domain

import (
	"time"
)

// User представляет пользователя бота
type User struct {
	ID           int64     `json:"id"`
	TelegramID   int64     `json:"telegram_id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	LanguageCode string    `json:"language_code"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsActive проверяет, активен ли пользователь
func (u *User) IsActive() bool {
	return !u.IsBlocked
}

// GetDisplayName возвращает отображаемое имя пользователя
func (u *User) GetDisplayName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.Username != "" {
		return "@" + u.Username
	}
	return "Пользователь"
}
