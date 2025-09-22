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

// ReferralRepository реализует domain.ReferralRepository
type ReferralRepository struct {
	repo *Repository
}

// NewReferralRepository создает новый репозиторий рефералов
func NewReferralRepository(repo *Repository) *ReferralRepository {
	return &ReferralRepository{
		repo: repo,
	}
}

// Create создает новую реферальную связь
func (r *ReferralRepository) Create(ctx context.Context, referral *domain.Referral) error {
	query := `
		INSERT INTO referrals (referrer_id, referred_id, level, reward_days, is_paid, created_at, updated_at, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	db := r.repo.GetDB(ctx)
	err := db.QueryRow(ctx, query,
		referral.ReferrerID,
		referral.ReferredID,
		referral.Level,
		referral.RewardDays,
		referral.IsPaid,
		referral.CreatedAt,
		referral.UpdatedAt,
		referral.PaidAt,
	).Scan(&referral.ID)

	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}

	return nil
}

// GetByID получает реферальную связь по ID
func (r *ReferralRepository) GetByID(ctx context.Context, id int64) (*domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referred_id, level, reward_days, is_paid, created_at, updated_at, paid_at
		FROM referrals
		WHERE id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, id)

	referral, err := r.scanReferral(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReferralNotFound
		}
		return nil, fmt.Errorf("failed to get referral by id: %w", err)
	}

	return referral, nil
}

// GetByReferredID получает реферальную связь по ID приглашенного
func (r *ReferralRepository) GetByReferredID(ctx context.Context, referredID int64) (*domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referred_id, level, reward_days, is_paid, created_at, updated_at, paid_at
		FROM referrals
		WHERE referred_id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, referredID)

	referral, err := r.scanReferral(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReferralNotFound
		}
		return nil, fmt.Errorf("failed to get referral by referred id: %w", err)
	}

	return referral, nil
}

// GetByReferrerID получает реферальные связи по ID приглашающего
func (r *ReferralRepository) GetByReferrerID(ctx context.Context, referrerID int64) ([]*domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referred_id, level, reward_days, is_paid, created_at, updated_at, paid_at
		FROM referrals
		WHERE referrer_id = $1
		ORDER BY created_at DESC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query, referrerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referrals by referrer id: %w", err)
	}
	defer rows.Close()

	var referrals []*domain.Referral
	for rows.Next() {
		referral, err := r.scanReferral(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referral: %w", err)
		}
		referrals = append(referrals, referral)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate referrals: %w", err)
	}

	return referrals, nil
}

// GetStats получает статистику рефералов пользователя
func (r *ReferralRepository) GetStats(ctx context.Context, userID int64) (*domain.ReferralStats, error) {
	query := `
		SELECT 
			$1 as user_id,
			COUNT(*) as total_referrals,
			COUNT(CASE WHEN level = 1 THEN 1 END) as level1_referrals,
			COUNT(CASE WHEN level = 2 THEN 1 END) as level2_referrals,
			COALESCE(SUM(reward_days), 0) as total_reward_days,
			COALESCE(SUM(CASE WHEN is_paid THEN reward_days ELSE 0 END), 0) as paid_reward_days,
			COALESCE(SUM(CASE WHEN NOT is_paid THEN reward_days ELSE 0 END), 0) as pending_reward_days
		FROM referrals
		WHERE referrer_id = $1`

	db := r.repo.GetDB(ctx)
	row := db.QueryRow(ctx, query, userID)

	var stats domain.ReferralStats
	err := row.Scan(
		&stats.UserID,
		&stats.TotalReferrals,
		&stats.Level1Referrals,
		&stats.Level2Referrals,
		&stats.TotalRewardDays,
		&stats.PaidRewardDays,
		&stats.PendingRewardDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get referral stats: %w", err)
	}

	return &stats, nil
}

// MarkAsPaid отмечает реферальную награду как выплаченную
func (r *ReferralRepository) MarkAsPaid(ctx context.Context, id int64) error {
	now := time.Now()
	query := `
		UPDATE referrals 
		SET is_paid = true, paid_at = $2, updated_at = $3 
		WHERE id = $1 AND NOT is_paid`

	db := r.repo.GetDB(ctx)
	result, err := db.Exec(ctx, query, id, &now, now)
	if err != nil {
		return fmt.Errorf("failed to mark referral as paid: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrReferralNotFound
	}

	return nil
}

// GetUnpaidRewards получает невыплаченные награды
func (r *ReferralRepository) GetUnpaidRewards(ctx context.Context) ([]*domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referred_id, level, reward_days, is_paid, created_at, updated_at, paid_at
		FROM referrals
		WHERE is_paid = false
		ORDER BY created_at ASC`

	db := r.repo.GetDB(ctx)
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get unpaid rewards: %w", err)
	}
	defer rows.Close()

	var referrals []*domain.Referral
	for rows.Next() {
		referral, err := r.scanReferral(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referral: %w", err)
		}
		referrals = append(referrals, referral)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate referrals: %w", err)
	}

	return referrals, nil
}

// scanReferral сканирует реферальную связь из строки результата
func (r *ReferralRepository) scanReferral(row pgx.Row) (*domain.Referral, error) {
	var referral domain.Referral
	var paidAt sql.NullTime

	err := row.Scan(
		&referral.ID,
		&referral.ReferrerID,
		&referral.ReferredID,
		&referral.Level,
		&referral.RewardDays,
		&referral.IsPaid,
		&referral.CreatedAt,
		&referral.UpdatedAt,
		&paidAt,
	)

	if err != nil {
		return nil, err
	}

	if paidAt.Valid {
		referral.PaidAt = &paidAt.Time
	}

	return &referral, nil
}
