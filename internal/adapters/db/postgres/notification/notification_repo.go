package notification

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
)

type Notification struct {
	dbGetter transactorPgx.DBGetter
}

func NewNotification(dbGetter transactorPgx.DBGetter) *Notification {
	return &Notification{
		dbGetter: dbGetter,
	}
}

func (n *Notification) CreateNotification(ctx context.Context, notification *core.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := n.dbGetter(ctx).Exec(ctx, query,
		notification.ID, notification.UserID, notification.Type,
		notification.Title, notification.Message, notification.IsRead,
		notification.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func (n *Notification) GetNotificationByID(ctx context.Context, id string) (*core.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, created_at
		FROM notifications WHERE id = $1`

	notification := &core.Notification{}
	err := n.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&notification.ID, &notification.UserID, &notification.Type,
		&notification.Title, &notification.Message, &notification.IsRead,
		&notification.CreatedAt,
	)

	if err != nil {
		return nil, usecase.ErrNotFound
	}

	return notification, nil
}

func (n *Notification) GetNotificationsByUserID(ctx context.Context, userID int64) ([]*core.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := n.dbGetter(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications by user ID: %w", err)
	}
	defer rows.Close()

	var notifications []*core.Notification
	for rows.Next() {
		notification := &core.Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Title, &notification.Message, &notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

func (n *Notification) GetUnreadNotificationsByUserID(ctx context.Context, userID int64) ([]*core.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, created_at
		FROM notifications WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC`

	rows, err := n.dbGetter(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications by user ID: %w", err)
	}
	defer rows.Close()

	var notifications []*core.Notification
	for rows.Next() {
		notification := &core.Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Title, &notification.Message, &notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, nil
}

func (n *Notification) UpdateNotification(ctx context.Context, notification *core.Notification) error {
	query := `
		UPDATE notifications 
		SET type = $2, title = $3, message = $4, is_read = $5
		WHERE id = $1`

	result, err := n.dbGetter(ctx).Exec(ctx, query,
		notification.ID, notification.Type, notification.Title,
		notification.Message, notification.IsRead,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (n *Notification) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1`

	result, err := n.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

func (n *Notification) DeleteNotification(ctx context.Context, id string) error {
	query := `DELETE FROM notifications WHERE id = $1`

	_, err := n.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}
