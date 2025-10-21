package usecase

import (
	"3xui-bot/internal/ports"
	"context"
	"time"

	"3xui-bot/internal/core"
)

type ReferralUseCase struct {
	referralRepo ports.ReferralRepo
	linkRepo     ports.ReferralLinkRepo
}

func NewReferralUseCase(referralRepo ports.ReferralRepo, linkRepo ports.ReferralLinkRepo) *ReferralUseCase {

	return &ReferralUseCase{
		referralRepo: referralRepo,
		linkRepo:     linkRepo,
	}
}

func (uc *ReferralUseCase) GetReferralLink(ctx context.Context, userID int64) (*core.ReferralLink, error) {
	link, err := uc.linkRepo.GetReferralLinkByUserID(ctx, userID)
	if err != nil {

		return nil, err
	}

	if link == nil {
		newLink := &core.ReferralLink{
			UserID:    userID,
			Link:      generateReferralCode(userID),
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = uc.linkRepo.CreateReferralLink(ctx, newLink)
		if err != nil {

			return nil, err
		}

		return newLink, nil
	}

	return link, nil
}

func (uc *ReferralUseCase) ProcessReferral(ctx context.Context, referrerID, refereeID int64) error {
	newReferral := &core.Referral{
		ReferrerID: referrerID,
		RefereeID:  refereeID,
		CreatedAt:  time.Now(),
	}

	return uc.referralRepo.CreateReferral(ctx, newReferral)
}

type ReferralStats struct {
	TotalReferrals int
	UserID         int64
}

func (uc *ReferralUseCase) GetReferralStats(ctx context.Context, userID int64) (*ReferralStats, error) {
	referrals, err := uc.referralRepo.GetReferralsByReferrerID(ctx, userID)
	if err != nil {

		return nil, err
	}

	stats := &ReferralStats{
		TotalReferrals: len(referrals),
		UserID:         userID,
	}

	return stats, nil
}

func (uc *ReferralUseCase) GetReferrals(ctx context.Context, userID int64) ([]*core.Referral, error) {

	return uc.referralRepo.GetReferralsByReferrerID(ctx, userID)
}

func generateReferralCode(userID int64) string {

	return "ref_" + string(rune(userID%1000))
}
