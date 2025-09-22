package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

// ServerRepository реализует domain.ServerRepository
type ServerRepository struct {
	repo *Repository
}

// NewServerRepository создает новый репозиторий серверов
func NewServerRepository(repo *Repository) *ServerRepository {
	return &ServerRepository{
		repo: repo,
	}
}

// Create создает новый сервер
func (r *ServerRepository) Create(ctx context.Context, server *domain.Server) error {
	query := `
		INSERT INTO servers (name, host, port, username, password, status, max_clients, current_clients, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		server.Name,
		server.Host,
		server.Port,
		server.Username,
		server.Password,
		server.Status,
		server.MaxClients,
		server.CurrentClients,
		server.CreatedAt,
		server.UpdatedAt,
	).Scan(&server.ID)

	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return nil
}

// GetByID получает сервер по ID
func (r *ServerRepository) GetByID(ctx context.Context, id int64) (*domain.Server, error) {
	query := `
		SELECT id, name, host, port, username, password, status, max_clients, current_clients, created_at, updated_at
		FROM servers
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	server, err := r.scanServer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrServerNotFound
		}
		return nil, fmt.Errorf("failed to get server by id: %w", err)
	}

	return server, nil
}

// GetAvailable получает доступные серверы
func (r *ServerRepository) GetAvailable(ctx context.Context) ([]*domain.Server, error) {
	query := `
		SELECT id, name, host, port, username, password, status, max_clients, current_clients, created_at, updated_at
		FROM servers
		WHERE status = $1 AND current_clients < max_clients
		ORDER BY current_clients ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, domain.ServerStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get available servers: %w", err)
	}
	defer rows.Close()

	var servers []*domain.Server
	for rows.Next() {
		server, err := r.scanServer(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}
		servers = append(servers, server)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate servers: %w", err)
	}

	return servers, nil
}

// Update обновляет сервер
func (r *ServerRepository) Update(ctx context.Context, server *domain.Server) error {
	query := `
		UPDATE servers
		SET name = $2, host = $3, port = $4, username = $5, password = $6, status = $7, max_clients = $8, current_clients = $9, updated_at = $10
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		server.ID,
		server.Name,
		server.Host,
		server.Port,
		server.Username,
		server.Password,
		server.Status,
		server.MaxClients,
		server.CurrentClients,
		server.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update server: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrServerNotFound
	}

	return nil
}

// Delete удаляет сервер
func (r *ServerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM servers WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrServerNotFound
	}

	return nil
}

// SetStatus устанавливает статус сервера
func (r *ServerRepository) SetStatus(ctx context.Context, id int64, status domain.ServerStatus) error {
	query := `UPDATE servers SET status = $2, updated_at = $3 WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set server status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrServerNotFound
	}

	return nil
}

// IncrementClients увеличивает количество клиентов на сервере
func (r *ServerRepository) IncrementClients(ctx context.Context, id int64) error {
	query := `
		UPDATE servers 
		SET current_clients = current_clients + 1, updated_at = $2 
		WHERE id = $1 AND current_clients < max_clients`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to increment server clients: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrServerOverloaded
	}

	return nil
}

// DecrementClients уменьшает количество клиентов на сервере
func (r *ServerRepository) DecrementClients(ctx context.Context, id int64) error {
	query := `
		UPDATE servers 
		SET current_clients = GREATEST(current_clients - 1, 0), updated_at = $2 
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to decrement server clients: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrServerNotFound
	}

	return nil
}

// scanServer сканирует сервер из строки результата
func (r *ServerRepository) scanServer(row pgx.Row) (*domain.Server, error) {
	var server domain.Server

	err := row.Scan(
		&server.ID,
		&server.Name,
		&server.Host,
		&server.Port,
		&server.Username,
		&server.Password,
		&server.Status,
		&server.MaxClients,
		&server.CurrentClients,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &server, nil
}
