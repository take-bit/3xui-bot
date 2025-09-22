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

// NotificationRepository реализует domain.NotificationRepository
type NotificationRepository struct {
	repo *Repository
}

// NewNotificationRepository создает новый репозиторий уведомлений
func NewNotificationRepository(repo *Repository) *NotificationRepository {
	return &NotificationRepository{
		repo: repo,
	}
}

// Create создает новое уведомление
func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (type, status, user_id, title, message, is_html, created_at, updated_at, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		notification.Type,
		notification.Status,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.IsHTML,
		notification.CreatedAt,
		notification.UpdatedAt,
		notification.SentAt,
	).Scan(&notification.ID)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetByID получает уведомление по ID
func (r *NotificationRepository) GetByID(ctx context.Context, id int64) (*domain.Notification, error) {
	query := `
		SELECT id, type, status, user_id, title, message, is_html, created_at, updated_at, sent_at
		FROM notifications
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	notification, err := r.scanNotification(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, fmt.Errorf("failed to get notification by id: %w", err)
	}

	return notification, nil
}

// Update обновляет уведомление
func (r *NotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	query := `
		UPDATE notifications
		SET type = $2, status = $3, user_id = $4, title = $5, message = $6, is_html = $7, updated_at = $8, sent_at = $9
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		notification.ID,
		notification.Type,
		notification.Status,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.IsHTML,
		notification.UpdatedAt,
		notification.SentAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}

// Delete удаляет уведомление
func (r *NotificationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notifications WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}

	return nil
}

// GetDraft получает черновики уведомлений
func (r *NotificationRepository) GetDraft(ctx context.Context) ([]*domain.Notification, error) {
	query := `
		SELECT id, type, status, user_id, title, message, is_html, created_at, updated_at, sent_at
		FROM notifications
		WHERE status = $1
		ORDER BY created_at DESC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, domain.NotificationStatusDraft)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		notification, err := r.scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notifications: %w", err)
	}

	return notifications, nil
}

// GetByUserID получает уведомления пользователя
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Notification, error) {
	query := `
		SELECT id, type, status, user_id, title, message, is_html, created_at, updated_at, sent_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by user id: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		notification, err := r.scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsSent отмечает уведомление как отправленное
func (r *NotificationRepository) MarkAsSent(ctx context.Context, id int64) error {
	now := time.Now()
	query := `
		UPDATE notifications 
		SET status = $2, sent_at = $3, updated_at = $4 
		WHERE id = $1 AND status = $5`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, domain.NotificationStatusSent, &now, now, domain.NotificationStatusDraft)
	if err != nil {
		return fmt.Errorf("failed to mark notification as sent: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotificationSent
	}

	return nil
}

// scanNotification сканирует уведомление из строки результата
func (r *NotificationRepository) scanNotification(row pgx.Row) (*domain.Notification, error) {
	var notification domain.Notification
	var userID sql.NullInt64
	var sentAt sql.NullTime

	err := row.Scan(
		&notification.ID,
		&notification.Type,
		&notification.Status,
		&userID,
		&notification.Title,
		&notification.Message,
		&notification.IsHTML,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&sentAt,
	)

	if err != nil {
		return nil, err
	}

	if userID.Valid {
		notification.UserID = &userID.Int64
	}

	if sentAt.Valid {
		notification.SentAt = &sentAt.Time
	}

	return &notification, nil
}
