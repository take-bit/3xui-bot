package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"3xui-bot/internal/domain"

	"github.com/jackc/pgx/v5"
)

// PlanRepository реализует domain.PlanRepository
type PlanRepository struct {
	repo *Repository
}

// NewPlanRepository создает новый репозиторий планов
func NewPlanRepository(repo *Repository) *PlanRepository {
	return &PlanRepository{
		repo: repo,
	}
}

// Create создает новый план
func (r *PlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	query := `
		INSERT INTO plans (name, devices, prices, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	pricesJSON, err := json.Marshal(plan.Prices)
	if err != nil {
		return fmt.Errorf("failed to marshal prices: %w", err)
	}

	db := r.repo.GetDB(ctx)
	err = db.QueryRow(ctx, query,
		plan.Name,
		plan.Devices,
		pricesJSON,
		plan.IsActive,
		plan.CreatedAt,
		plan.UpdatedAt,
	).Scan(&plan.ID)

	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}

	return nil
}

// GetByID получает план по ID
func (r *PlanRepository) GetByID(ctx context.Context, id int64) (*domain.Plan, error) {
	query := `
		SELECT id, name, devices, prices, is_active, created_at, updated_at
		FROM plans
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	plan, err := r.scanPlan(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPlanNotFound
		}
		return nil, fmt.Errorf("failed to get plan by id: %w", err)
	}

	return plan, nil
}

// GetActive получает активные планы
func (r *PlanRepository) GetActive(ctx context.Context) ([]*domain.Plan, error) {
	query := `
		SELECT id, name, devices, prices, is_active, created_at, updated_at
		FROM plans
		WHERE is_active = true
		ORDER BY created_at ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active plans: %w", err)
	}
	defer rows.Close()

	var plans []*domain.Plan
	for rows.Next() {
		plan, err := r.scanPlan(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate plans: %w", err)
	}

	return plans, nil
}

// Update обновляет план
func (r *PlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	query := `
		UPDATE plans
		SET name = $2, devices = $3, prices = $4, is_active = $5, updated_at = $6
		WHERE id = $1`

	pricesJSON, err := json.Marshal(plan.Prices)
	if err != nil {
		return fmt.Errorf("failed to marshal prices: %w", err)
	}

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query,
		plan.ID,
		plan.Name,
		plan.Devices,
		pricesJSON,
		plan.IsActive,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPlanNotFound
	}

	return nil
}

// Delete удаляет план
func (r *PlanRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM plans WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPlanNotFound
	}

	return nil
}

// SetActive активирует/деактивирует план
func (r *PlanRepository) SetActive(ctx context.Context, id int64, active bool) error {
	query := `UPDATE plans SET is_active = $2, updated_at = $3 WHERE id = $1`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, active, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set plan active status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPlanNotFound
	}

	return nil
}

// scanPlan сканирует план из строки результата
func (r *PlanRepository) scanPlan(row pgx.Row) (*domain.Plan, error) {
	var plan domain.Plan
	var pricesJSON []byte

	err := row.Scan(
		&plan.ID,
		&plan.Name,
		&plan.Devices,
		&pricesJSON,
		&plan.IsActive,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pricesJSON, &plan.Prices)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prices: %w", err)
	}

	return &plan, nil
}
