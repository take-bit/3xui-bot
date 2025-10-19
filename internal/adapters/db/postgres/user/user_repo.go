package user

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
)

type User struct {
	dbGetter transactorPgx.DBGetter
}

func NewUser(dbGetter transactorPgx.DBGetter) *User {
	return &User{
		dbGetter: dbGetter,
	}
}

func (u *User) CreateUser(ctx context.Context, user *core.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, language_code, is_blocked, has_trial, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := u.dbGetter(ctx).Exec(ctx, query,
		user.TelegramID, user.Username, user.FirstName, user.LastName,
		user.LanguageCode, user.IsBlocked, user.HasTrial, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user (telegram_id=%d): %w", user.TelegramID, err)
	}

	return nil
}

func (u *User) GetUserByID(ctx context.Context, id int64) (*core.User, error) {
	// id теперь это telegram_id
	return u.GetUserByTelegramID(ctx, id)
}

func (u *User) GetUserByTelegramID(ctx context.Context, telegramID int64) (*core.User, error) {
	query := `
		SELECT telegram_id, username, first_name, last_name, language_code, is_blocked, has_trial, created_at, updated_at
		FROM users WHERE telegram_id = $1`

	user := &core.User{}
	err := u.dbGetter(ctx).QueryRow(ctx, query, telegramID).Scan(
		&user.TelegramID, &user.Username, &user.FirstName,
		&user.LastName, &user.LanguageCode, &user.IsBlocked, &user.HasTrial,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return user, nil
}

func (u *User) UpdateUser(ctx context.Context, user *core.User) error {
	query := `
		UPDATE users 
		SET username = $2, first_name = $3, last_name = $4, language_code = $5, 
		    is_blocked = $6, has_trial = $7, updated_at = $8
		WHERE telegram_id = $1`

	result, err := u.dbGetter(ctx).Exec(ctx, query,
		user.TelegramID, user.Username, user.FirstName, user.LastName,
		user.LanguageCode, user.IsBlocked, user.HasTrial, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (u *User) DeleteUser(ctx context.Context, telegramID int64) error {
	query := `DELETE FROM users WHERE telegram_id = $1`

	_, err := u.dbGetter(ctx).Exec(ctx, query, telegramID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetUserState удален - данные получаем из основных таблиц

// SetUserState удален - данные обновляем в основных таблицах

// MarkTrialAsUsed отмечает, что пользователь использовал пробный период
func (u *User) MarkTrialAsUsed(ctx context.Context, userID int64) error {
	query := `UPDATE users SET has_trial = TRUE, updated_at = $2 WHERE telegram_id = $1`

	result, err := u.dbGetter(ctx).Exec(ctx, query, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark trial as used: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}
