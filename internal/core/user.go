package core

import (
	"time"
)

type User struct {
	TelegramID   int64     `json:"telegram_id" db:"telegram_id"`
	Username     string    `json:"username" db:"username"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	LanguageCode string    `json:"language_code" db:"language_code"`
	IsBlocked    bool      `json:"is_blocked" db:"is_blocked"`
	HasTrial     bool      `json:"has_trial" db:"has_trial"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) GetDisplayName() string {
	if u.FirstName != "" {

		return u.FirstName
	}
	if u.Username != "" {

		return "@" + u.Username
	}

	return "Пользователь"
}

func (u *User) IsActive() bool {

	return !u.IsBlocked
}
