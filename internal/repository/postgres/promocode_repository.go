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

// PromoRepository реализует domain.PromoRepository
type PromoRepository struct {
	repo *Repository
}

func NewPromoRepository(repo *Repository) *PromoRepository {
	return &PromoRepository{
		repo: repo,
	}
}

// Create создает новый промокод
func (r *PromoRepository) Create(ctx context.Context, promocode *domain.Promocode) error {
	query := `
		INSERT INTO promocodes (code, type, value, is_active, usage_limit, used_count, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		promocode.Code,
		promocode.Type,
		promocode.Value,
		promocode.IsActive,
		promocode.UsageLimit,
		promocode.UsedCount,
		promocode.ExpiresAt,
		promocode.CreatedAt,
		promocode.UpdatedAt,
	).Scan(&promocode.ID)

	if err != nil {
		return fmt.Errorf("failed to create promocode: %w", err)
	}

	return nil
}

// GetByID получает промокод по ID
func (r *PromoRepository) GetByID(ctx context.Context, id int64) (*domain.Promocode, error) {
	query := `
		SELECT id, code, type, value, is_active, usage_limit, used_count, expires_at, created_at, updated_at
		FROM promocodes
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	promocode, err := r.scanPromocode(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPromocodeNotFound
		}
		return nil, fmt.Errorf("failed to get promocode by id: %w", err)
	}

	return promocode, nil
}

// GetByCode получает промокод по коду
func (r *PromoRepository) GetByCode(ctx context.Context, code string) (*domain.Promocode, error) {
	query := `
		SELECT id, code, type, value, is_active, usage_limit, used_count, expires_at, created_at, updated_at
		FROM promocodes
		WHERE code = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, code)

	promocode, err := r.scanPromocode(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPromocodeNotFound
		}
		return nil, fmt.Errorf("failed to get promocode by code: %w", err)
	}

	return promocode, nil
}

// Update обновляет промокод
func (r *PromoRepository) Update(ctx context.Context, promocode *domain.Promocode) error {
	query := `
		UPDATE promocodes
		SET code = $2, type = $3, value = $4, is_active = $5, usage_limit = $6, used_count = $7, expires_at = $8, updated_at = $9
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		promocode.ID,
		promocode.Code,
		promocode.Type,
		promocode.Value,
		promocode.IsActive,
		promocode.UsageLimit,
		promocode.UsedCount,
		promocode.ExpiresAt,
		promocode.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update promocode: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPromocodeNotFound
	}

	return nil
}

// Delete удаляет промокод
func (r *PromoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM promocodes WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete promocode: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPromocodeNotFound
	}

	return nil
}

// IncrementUsage увеличивает счетчик использования промокода
func (r *PromoRepository) IncrementUsage(ctx context.Context, id int64) error {
	query := `
		UPDATE promocodes 
		SET used_count = used_count + 1, updated_at = $2 
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to increment promocode usage: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPromocodeNotFound
	}

	return nil
}

// GetActive получает активные промокоды
func (r *PromoRepository) GetActive(ctx context.Context) ([]*domain.Promocode, error) {
	query := `
		SELECT id, code, type, value, is_active, usage_limit, used_count, expires_at, created_at, updated_at
		FROM promocodes
		WHERE is_active = true AND (expires_at IS NULL OR expires_at > $1)
		ORDER BY created_at ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get active promocodes: %w", err)
	}
	defer rows.Close()

	var promocodes []*domain.Promocode
	for rows.Next() {
		promocode, err := r.scanPromocode(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan promocode: %w", err)
		}
		promocodes = append(promocodes, promocode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate promocodes: %w", err)
	}

	return promocodes, nil
}

// scanPromocode сканирует промокод из строки результата
func (r *PromoRepository) scanPromocode(row pgx.Row) (*domain.Promocode, error) {
	var promocode domain.Promocode
	var expiresAt sql.NullTime

	err := row.Scan(
		&promocode.ID,
		&promocode.Code,
		&promocode.Type,
		&promocode.Value,
		&promocode.IsActive,
		&promocode.UsageLimit,
		&promocode.UsedCount,
		&expiresAt,
		&promocode.CreatedAt,
		&promocode.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		promocode.ExpiresAt = &expiresAt.Time
	}

	return &promocode, nil
}
