package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

// UserRepository реализует domain.UserRepository
type UserRepository struct {
	repo *Repository
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(repo *Repository) *UserRepository {
	return &UserRepository{
		repo: repo,
	}
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, language_code, is_blocked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
		user.IsBlocked,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, language_code, is_blocked, created_at, updated_at
		FROM users
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	user, err := r.scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByTelegramID получает пользователя по Telegram ID
func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, language_code, is_blocked, created_at, updated_at
		FROM users
		WHERE telegram_id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, telegramID)

	user, err := r.scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by telegram id: %w", err)
	}

	return user, nil
}

// Update обновляет пользователя
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, first_name = $3, last_name = $4, language_code = $5, is_blocked = $6, updated_at = $7
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
		user.IsBlocked,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List получает список пользователей с пагинацией
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, language_code, is_blocked, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return users, nil
}

// Block блокирует пользователя
func (r *UserRepository) Block(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_blocked = true, updated_at = $2 WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Unblock разблокирует пользователя
func (r *UserRepository) Unblock(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_blocked = false, updated_at = $2 WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// scanUser сканирует пользователя из строки результата
func (r *UserRepository) scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	var username, firstName, lastName, languageCode sql.NullString

	err := row.Scan(
		&user.ID,
		&user.TelegramID,
		&username,
		&firstName,
		&lastName,
		&languageCode,
		&user.IsBlocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.Username = username.String
	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.LanguageCode = languageCode.String

	return &user, nil
}
