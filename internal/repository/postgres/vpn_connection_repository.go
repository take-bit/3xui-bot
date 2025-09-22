package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

// VpnConnectionRepository реализует интерфейс для работы с VPN подключениями
type VpnConnectionRepository struct {
	repo *Repository
}

// NewVpnConnectionRepository создает новый репозиторий VPN подключений
func NewVpnConnectionRepository(repo *Repository) *VpnConnectionRepository {
	return &VpnConnectionRepository{repo: repo}
}

// Create создает новое VPN подключение
func (r *VpnConnectionRepository) Create(ctx context.Context, connection *domain.VPNConnection) error {
	query := `
		INSERT INTO vpn_connections (
			user_id, server_id, xui_inbound_id, xui_client_id, 
			uuid, email, config_url, created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		connection.UserID,
		connection.ServerID,
		connection.XUIInboundID,
		connection.XUIClientID,
		connection.UUID,
		connection.Email,
		connection.ConfigURL,
		connection.CreatedAt,
		connection.ExpiresAt,
	).Scan(&connection.UserID) // Используем UserID как ID для простоты

	return err
}

// GetByID получает VPN подключение по ID
func (r *VpnConnectionRepository) GetByID(ctx context.Context, id int64) (*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE user_id = $1`

	var connection domain.VPNConnection
	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query, id).Scan(
		&connection.UserID,
		&connection.ServerID,
		&connection.XUIInboundID,
		&connection.XUIClientID,
		&connection.UUID,
		&connection.Email,
		&connection.ConfigURL,
		&connection.CreatedAt,
		&connection.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &connection, nil
}

// GetByUserID получает VPN подключение по ID пользователя
func (r *VpnConnectionRepository) GetByUserID(ctx context.Context, userID int64) (*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE user_id = $1`

	var connection domain.VPNConnection
	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query, userID).Scan(
		&connection.UserID,
		&connection.ServerID,
		&connection.XUIInboundID,
		&connection.XUIClientID,
		&connection.UUID,
		&connection.Email,
		&connection.ConfigURL,
		&connection.CreatedAt,
		&connection.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &connection, nil
}

// GetByServerID получает все VPN подключения для сервера
func (r *VpnConnectionRepository) GetByServerID(ctx context.Context, serverID int64) ([]*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE server_id = $1
		ORDER BY created_at DESC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*domain.VPNConnection
	for rows.Next() {
		var connection domain.VPNConnection
		err := rows.Scan(
			&connection.UserID,
			&connection.ServerID,
			&connection.XUIInboundID,
			&connection.XUIClientID,
			&connection.UUID,
			&connection.Email,
			&connection.ConfigURL,
			&connection.CreatedAt,
			&connection.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		connections = append(connections, &connection)
	}

	return connections, rows.Err()
}

// GetByXUIInboundID получает VPN подключение по ID сервера и inbound
func (r *VpnConnectionRepository) GetByXUIInboundID(ctx context.Context, serverID int64, inboundID int) (*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE server_id = $1 AND xui_inbound_id = $2`

	var connection domain.VPNConnection
	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query, serverID, inboundID).Scan(
		&connection.UserID,
		&connection.ServerID,
		&connection.XUIInboundID,
		&connection.XUIClientID,
		&connection.UUID,
		&connection.Email,
		&connection.ConfigURL,
		&connection.CreatedAt,
		&connection.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &connection, nil
}

// Update обновляет VPN подключение
func (r *VpnConnectionRepository) Update(ctx context.Context, connection *domain.VPNConnection) error {
	query := `
		UPDATE vpn_connections 
		SET server_id = $2, xui_inbound_id = $3, xui_client_id = $4, 
		    uuid = $5, email = $6, config_url = $7, expires_at = $8
		WHERE user_id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		connection.UserID,
		connection.ServerID,
		connection.XUIInboundID,
		connection.XUIClientID,
		connection.UUID,
		connection.Email,
		connection.ConfigURL,
		connection.ExpiresAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete удаляет VPN подключение по ID
func (r *VpnConnectionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM vpn_connections WHERE user_id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteByUserID удаляет VPN подключение по ID пользователя
func (r *VpnConnectionRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM vpn_connections WHERE user_id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetExpired получает все истекшие VPN подключения
func (r *VpnConnectionRepository) GetExpired(ctx context.Context) ([]*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE expires_at < $1
		ORDER BY expires_at ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*domain.VPNConnection
	for rows.Next() {
		var connection domain.VPNConnection
		err := rows.Scan(
			&connection.UserID,
			&connection.ServerID,
			&connection.XUIInboundID,
			&connection.XUIClientID,
			&connection.UUID,
			&connection.Email,
			&connection.ConfigURL,
			&connection.CreatedAt,
			&connection.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		connections = append(connections, &connection)
	}

	return connections, rows.Err()
}

// GetActiveByUserID получает активное VPN подключение пользователя
func (r *VpnConnectionRepository) GetActiveByUserID(ctx context.Context, userID int64) (*domain.VPNConnection, error) {
	query := `
		SELECT user_id, server_id, xui_inbound_id, xui_client_id, 
		       uuid, email, config_url, created_at, expires_at
		FROM vpn_connections 
		WHERE user_id = $1 AND expires_at > $2`

	var connection domain.VPNConnection
	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query, userID, time.Now()).Scan(
		&connection.UserID,
		&connection.ServerID,
		&connection.XUIInboundID,
		&connection.XUIClientID,
		&connection.UUID,
		&connection.Email,
		&connection.ConfigURL,
		&connection.CreatedAt,
		&connection.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &connection, nil
}

// Extend продлевает VPN подключение на указанное количество дней
func (r *VpnConnectionRepository) Extend(ctx context.Context, id int64, days int) error {
	query := `
		UPDATE vpn_connections 
		SET expires_at = expires_at + INTERVAL '%d days'
		WHERE user_id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, fmt.Sprintf(query, days), id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
