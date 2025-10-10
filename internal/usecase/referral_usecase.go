package usecase

import (
	"3xui-bot/internal/ports"
	"context"
	"time"

	"3xui-bot/internal/core"
)

// ReferralUseCase use case для работы с рефералами
type ReferralUseCase struct {
	referralRepo ports.ReferralRepo
	linkRepo     ports.ReferralLinkRepo
}

// NewReferralUseCase создает новый use case для рефералов
func NewReferralUseCase(referralRepo ports.ReferralRepo, linkRepo ports.ReferralLinkRepo) *ReferralUseCase {
	return &ReferralUseCase{
		referralRepo: referralRepo,
		linkRepo:     linkRepo,
	}
}

// GetReferralLink получает реферальную ссылку пользователя
func (uc *ReferralUseCase) GetReferralLink(ctx context.Context, userID int64) (*core.ReferralLink, error) {
	link, err := uc.linkRepo.GetReferralLinkByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if link == nil {
		// Создаем новую реферальную ссылку
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

// ProcessReferral обрабатывает реферальное приглашение
func (uc *ReferralUseCase) ProcessReferral(ctx context.Context, referrerID, refereeID int64) error {
	newReferral := &core.Referral{
		ReferrerID: referrerID,
		RefereeID:  refereeID,
		CreatedAt:  time.Now(),
	}

	return uc.referralRepo.CreateReferral(ctx, newReferral)
}

// ReferralStats статистика рефералов
type ReferralStats struct {
	TotalReferrals int
	UserID         int64
}

// GetReferralStats получает статистику рефералов пользователя
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

// GetReferrals получает всех рефералов пользователя
func (uc *ReferralUseCase) GetReferrals(ctx context.Context, userID int64) ([]*core.Referral, error) {
	return uc.referralRepo.GetReferralsByReferrerID(ctx, userID)
}

// generateReferralCode генерирует реферальный код
func generateReferralCode(userID int64) string {
	// Простая генерация кода на основе userID
	return "ref_" + string(rune(userID%1000))
}
