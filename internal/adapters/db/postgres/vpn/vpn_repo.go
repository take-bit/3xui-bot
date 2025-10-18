package vpn

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
	"github.com/jackc/pgx/v5"
)

type VPNConnection struct {
	dbGetter transactorPgx.DBGetter
}

func NewVPNConnection(dbGetter transactorPgx.DBGetter) *VPNConnection {
	return &VPNConnection{
		dbGetter: dbGetter,
	}
}

// CreateVPNConnection создает VPN подключение
func (v *VPNConnection) CreateVPNConnection(ctx context.Context, conn *core.VPNConnection) error {
	query := `
		INSERT INTO vpn_connections (id, telegram_user_id, marzban_username, name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := v.dbGetter(ctx).Exec(ctx, query,
		conn.ID, conn.TelegramUserID, conn.MarzbanUsername, conn.Name,
		conn.IsActive, conn.CreatedAt, conn.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create VPN connection: %w", err)
	}
	return nil
}

// GetVPNConnectionsByTelegramUserID получает все VPN подключения пользователя
func (v *VPNConnection) GetVPNConnectionsByTelegramUserID(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {
	query := `
		SELECT id, telegram_user_id, marzban_username, name, is_active, created_at, updated_at
		FROM vpn_connections WHERE telegram_user_id = $1 ORDER BY created_at DESC`

	rows, err := v.dbGetter(ctx).Query(ctx, query, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN connections: %w", err)
	}
	defer rows.Close()

	var connections []*core.VPNConnection
	for rows.Next() {
		conn := &core.VPNConnection{}
		err := rows.Scan(
			&conn.ID, &conn.TelegramUserID, &conn.MarzbanUsername, &conn.Name,
			&conn.IsActive, &conn.CreatedAt, &conn.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan VPN connection: %w", err)
		}
		connections = append(connections, conn)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating VPN connections: %w", err)
	}
	return connections, nil
}

// GetVPNConnectionsBySubscriptionID получает все VPN подключения для подписки
// Пока реализовано через получение подписки и VPN по user_id
func (v *VPNConnection) GetVPNConnectionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]*core.VPNConnection, error) {
	// Получаем user_id из подписки
	var telegramUserID int64
	query := `SELECT user_id FROM subscriptions WHERE id = $1`
	err := v.dbGetter(ctx).QueryRow(ctx, query, subscriptionID).Scan(&telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_id from subscription: %w", err)
	}

	// Получаем VPN подключения этого пользователя
	return v.GetVPNConnectionsByTelegramUserID(ctx, telegramUserID)
}

// GetVPNConnectionByID получает VPN подключение по ID
func (v *VPNConnection) GetVPNConnectionByID(ctx context.Context, id string) (*core.VPNConnection, error) {
	query := `
		SELECT id, telegram_user_id, marzban_username, name, is_active, created_at, updated_at
		FROM vpn_connections WHERE id = $1`

	conn := &core.VPNConnection{}
	err := v.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&conn.ID, &conn.TelegramUserID, &conn.MarzbanUsername, &conn.Name,
		&conn.IsActive, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get VPN connection: %w", err)
	}
	return conn, nil
}

// GetVPNConnectionByMarzbanUsername получает VPN подключение по Marzban username
func (v *VPNConnection) GetVPNConnectionByMarzbanUsername(ctx context.Context, marzbanUsername string) (*core.VPNConnection, error) {
	query := `
		SELECT id, telegram_user_id, marzban_username, name, is_active, created_at, updated_at
		FROM vpn_connections WHERE marzban_username = $1`

	conn := &core.VPNConnection{}
	err := v.dbGetter(ctx).QueryRow(ctx, query, marzbanUsername).Scan(
		&conn.ID, &conn.TelegramUserID, &conn.MarzbanUsername, &conn.Name,
		&conn.IsActive, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get VPN connection by Marzban username: %w", err)
	}
	return conn, nil
}

// UpdateVPNConnectionName обновляет только локальное имя
func (v *VPNConnection) UpdateVPNConnectionName(ctx context.Context, id, name string) error {
	query := `UPDATE vpn_connections SET name = $2, updated_at = $3 WHERE id = $1`

	result, err := v.dbGetter(ctx).Exec(ctx, query, id, name, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update VPN connection name: %w", err)
	}
	if result.RowsAffected() == 0 {
		return usecase.ErrNotFound
	}
	return nil
}

// DeleteVPNConnection удаляет VPN подключение
func (v *VPNConnection) DeleteVPNConnection(ctx context.Context, id string) error {
	query := `DELETE FROM vpn_connections WHERE id = $1`

	_, err := v.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete VPN connection: %w", err)
	}
	return nil
}

// DeleteVPNConnectionByMarzbanUsername удаляет VPN подключение по Marzban username
func (v *VPNConnection) DeleteVPNConnectionByMarzbanUsername(ctx context.Context, marzbanUsername string) error {
	query := `DELETE FROM vpn_connections WHERE marzban_username = $1`

	_, err := v.dbGetter(ctx).Exec(ctx, query, marzbanUsername)
	if err != nil {
		return fmt.Errorf("failed to delete VPN connection by Marzban username: %w", err)
	}
	return nil
}

// GetActiveVPNConnections получает только активные VPN подключения пользователя
func (v *VPNConnection) GetActiveVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {
	query := `
		SELECT id, telegram_user_id, marzban_username, name, is_active, created_at, updated_at
		FROM vpn_connections WHERE telegram_user_id = $1 AND is_active = TRUE ORDER BY created_at DESC`

	rows, err := v.dbGetter(ctx).Query(ctx, query, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active VPN connections: %w", err)
	}
	defer rows.Close()

	var connections []*core.VPNConnection
	for rows.Next() {
		conn := &core.VPNConnection{}
		err := rows.Scan(
			&conn.ID, &conn.TelegramUserID, &conn.MarzbanUsername, &conn.Name,
			&conn.IsActive, &conn.CreatedAt, &conn.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan VPN connection: %w", err)
		}
		connections = append(connections, conn)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating VPN connections: %w", err)
	}
	return connections, nil
}
