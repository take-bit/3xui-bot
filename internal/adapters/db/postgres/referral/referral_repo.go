package referral

import (
	"context"
	"fmt"

	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
	"github.com/jackc/pgx/v5"
)

type Referral struct {
	dbGetter transactorPgx.DBGetter
}

func NewReferral(dbGetter transactorPgx.DBGetter) *Referral {

	return &Referral{
		dbGetter: dbGetter,
	}
}

func (r *Referral) CreateReferral(ctx context.Context, referral *core.Referral) error {
	query := `
		INSERT INTO referrals (referrer_id, referee_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.dbGetter(ctx).QueryRow(ctx, query,
		referral.ReferrerID, referral.RefereeID, referral.CreatedAt,
	).Scan(&referral.ID)

	if err != nil {

		return fmt.Errorf("failed to create referral: %w", err)
	}

	return nil
}

func (r *Referral) GetReferralByID(ctx context.Context, id int64) (*core.Referral, error) {
	query := `
		SELECT id, referrer_id, referee_id, created_at
		FROM referrals WHERE id = $1`

	referral := &core.Referral{}
	err := r.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&referral.ID, &referral.ReferrerID, &referral.RefereeID, &referral.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {

			return nil, usecase.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get referral by ID: %w", err)
	}

	return referral, nil
}

func (r *Referral) GetReferralsByReferrerID(ctx context.Context, referrerID int64) ([]*core.Referral, error) {
	query := `
		SELECT id, referrer_id, referee_id, created_at
		FROM referrals WHERE referrer_id = $1
		ORDER BY created_at DESC`

	rows, err := r.dbGetter(ctx).Query(ctx, query, referrerID)
	if err != nil {

		return nil, fmt.Errorf("failed to get referrals by referrer ID: %w", err)
	}
	defer rows.Close()

	var referrals []*core.Referral
	for rows.Next() {
		referral := &core.Referral{}
		err := rows.Scan(
			&referral.ID, &referral.ReferrerID, &referral.RefereeID, &referral.CreatedAt,
		)
		if err != nil {

			return nil, fmt.Errorf("failed to scan referral: %w", err)
		}
		referrals = append(referrals, referral)
	}

	if err = rows.Err(); err != nil {

		return nil, fmt.Errorf("error iterating referrals: %w", err)
	}

	return referrals, nil
}

func (r *Referral) GetReferralByRefereeID(ctx context.Context, refereeID int64) (*core.Referral, error) {
	query := `
		SELECT id, referrer_id, referee_id, created_at
		FROM referrals WHERE referee_id = $1`

	referral := &core.Referral{}
	err := r.dbGetter(ctx).QueryRow(ctx, query, refereeID).Scan(
		&referral.ID, &referral.ReferrerID, &referral.RefereeID, &referral.CreatedAt,
	)

	if err != nil {

		return nil, usecase.ErrNotFound
	}

	return referral, nil
}

func (r *Referral) DeleteReferral(ctx context.Context, id int64) error {
	query := `DELETE FROM referrals WHERE id = $1`

	_, err := r.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {

		return fmt.Errorf("failed to delete referral: %w", err)
	}

	return nil
}

type ReferralLink struct {
	dbGetter transactorPgx.DBGetter
}

func NewReferralLink(dbGetter transactorPgx.DBGetter) *ReferralLink {

	return &ReferralLink{
		dbGetter: dbGetter,
	}
}

func (rl *ReferralLink) CreateReferralLink(ctx context.Context, link *core.ReferralLink) error {
	query := `
		INSERT INTO referral_links (user_id, link, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := rl.dbGetter(ctx).QueryRow(ctx, query,
		link.UserID, link.Link, link.IsActive, link.CreatedAt, link.UpdatedAt,
	).Scan(&link.ID)

	if err != nil {

		return fmt.Errorf("failed to create referral link: %w", err)
	}

	return nil
}

func (rl *ReferralLink) GetReferralLinkByID(ctx context.Context, id int64) (*core.ReferralLink, error) {
	query := `
		SELECT id, user_id, link, is_active, created_at, updated_at
		FROM referral_links WHERE id = $1`

	link := &core.ReferralLink{}
	err := rl.dbGetter(ctx).QueryRow(ctx, query, id).Scan(
		&link.ID, &link.UserID, &link.Link, &link.IsActive,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {

			return nil, usecase.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get referral link by ID: %w", err)
	}

	return link, nil
}

func (rl *ReferralLink) GetReferralLinkByUserID(ctx context.Context, userID int64) (*core.ReferralLink, error) {
	query := `
		SELECT id, user_id, link, is_active, created_at, updated_at
		FROM referral_links WHERE user_id = $1`

	link := &core.ReferralLink{}
	err := rl.dbGetter(ctx).QueryRow(ctx, query, userID).Scan(
		&link.ID, &link.UserID, &link.Link, &link.IsActive,
		&link.CreatedAt, &link.UpdatedAt,
	)

	if err != nil {

		return nil, usecase.ErrNotFound
	}

	return link, nil
}

func (rl *ReferralLink) GetReferralLinkByLink(ctx context.Context, link string) (*core.ReferralLink, error) {
	query := `
		SELECT id, user_id, link, is_active, created_at, updated_at
		FROM referral_links WHERE link = $1`

	referralLink := &core.ReferralLink{}
	err := rl.dbGetter(ctx).QueryRow(ctx, query, link).Scan(
		&referralLink.ID, &referralLink.UserID, &referralLink.Link, &referralLink.IsActive,
		&referralLink.CreatedAt, &referralLink.UpdatedAt,
	)

	if err != nil {

		return nil, usecase.ErrNotFound
	}

	return referralLink, nil
}

func (rl *ReferralLink) UpdateReferralLink(ctx context.Context, link *core.ReferralLink) error {
	query := `
		UPDATE referral_links
		SET link = $2, is_active = $3, updated_at = $4
		WHERE id = $1`

	result, err := rl.dbGetter(ctx).Exec(ctx, query,
		link.ID, link.Link, link.IsActive, link.UpdatedAt,
	)

	if err != nil {

		return fmt.Errorf("failed to update referral link: %w", err)
	}

	if result.RowsAffected() == 0 {

		return usecase.ErrNotFound
	}

	return nil
}

func (rl *ReferralLink) DeleteReferralLink(ctx context.Context, id int64) error {
	query := `DELETE FROM referral_links WHERE id = $1`

	_, err := rl.dbGetter(ctx).Exec(ctx, query, id)
	if err != nil {

		return fmt.Errorf("failed to delete referral link: %w", err)
	}

	return nil
}
