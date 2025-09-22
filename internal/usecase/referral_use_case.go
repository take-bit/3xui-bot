package usecase

import (
	"context"
	"fmt"

	"3xui-bot/internal/domain"
)

// ReferralUseCase представляет use case для работы с реферальной программой
type ReferralUseCase struct {
	referralService     domain.ReferralService
	userService         domain.UserService
	subscriptionService domain.SubscriptionService
	notificationService domain.NotificationService
}

// NewReferralUseCase создает новый Referral use case
func NewReferralUseCase(
	referralService domain.ReferralService,
	userService domain.UserService,
	subscriptionService domain.SubscriptionService,
	notificationService domain.NotificationService,
) *ReferralUseCase {
	return &ReferralUseCase{
		referralService:     referralService,
		userService:         userService,
		subscriptionService: subscriptionService,
		notificationService: notificationService,
	}
}

// CreateReferral создает реферальную связь
func (uc *ReferralUseCase) CreateReferral(ctx context.Context, referrerID int64, referredID int64) error {
	// 1. Проверяем, что пользователи существуют
	referrer, err := uc.userService.GetByTelegramID(ctx, referrerID)
	if err != nil {
		return fmt.Errorf("failed to get referrer: %w", err)
	}

	referred, err := uc.userService.GetByTelegramID(ctx, referredID)
	if err != nil {
		return fmt.Errorf("failed to get referred user: %w", err)
	}

	// 2. Проверяем, что пользователь не рефералит сам себя
	if referrer.ID == referred.ID {
		return fmt.Errorf("user cannot refer themselves")
	}

	// 3. Создаем реферальную связь
	err = uc.referralService.CreateReferral(ctx, referrer.ID, referred.ID)
	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}

	// 4. Отправляем уведомления
	if uc.notificationService != nil {
		// Уведомление рефереру
		referrerMessage := fmt.Sprintf("🎉 У вас новый реферал!\n\n👤 Пользователь: %s\n💰 Вы получите вознаграждение после его первого платежа", referred.FirstName)
		_ = uc.notificationService.SendToUser(ctx, referrer.TelegramID, "Новый реферал", referrerMessage, false)

		// Уведомление рефералу
		referredMessage := fmt.Sprintf("👋 Добро пожаловать!\n\n🎁 Вы пришли по реферальной ссылке от %s\n💰 После первого платежа ваш реферер получит вознаграждение", referrer.FirstName)
		_ = uc.notificationService.SendToUser(ctx, referred.TelegramID, "Реферальная программа", referredMessage, false)
	}

	return nil
}

// ProcessPaymentReward обрабатывает вознаграждение за платеж реферала
func (uc *ReferralUseCase) ProcessPaymentReward(ctx context.Context, userID int64, amount int) error {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Обрабатываем вознаграждение
	err = uc.referralService.ProcessPaymentReward(ctx, user.ID, amount)
	if err != nil {
		return fmt.Errorf("failed to process payment reward: %w", err)
	}

	// 3. Отправляем уведомления о вознаграждении
	// TODO: Получить информацию о рефералах и отправить уведомления

	return nil
}

// GetReferralStats получает статистику рефералов пользователя
func (uc *ReferralUseCase) GetReferralStats(ctx context.Context, userID int64) (*ReferralStatsInfo, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем статистику рефералов
	stats, err := uc.referralService.GetStats(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referral stats: %w", err)
	}

	// 3. Формируем информацию для пользователя
	info := &ReferralStatsInfo{
		User:            user,
		ReferralStats:   stats,
		ReferralLink:    fmt.Sprintf("https://t.me/your_bot?start=ref_%d", user.ID),
		TotalEarnings:   stats.TotalRewardDays, // Используем TotalRewardDays как earnings
		PendingEarnings: stats.PendingRewardDays,
		PaidEarnings:    stats.PaidRewardDays,
	}

	return info, nil
}

// GetUnpaidRewards получает список неоплаченных вознаграждений
func (uc *ReferralUseCase) GetUnpaidRewards(ctx context.Context) ([]*domain.Referral, error) {
	return uc.referralService.GetUnpaidRewards(ctx)
}

// PayReward выплачивает вознаграждение
func (uc *ReferralUseCase) PayReward(ctx context.Context, referralID int64) error {
	// 1. Выплачиваем вознаграждение
	err := uc.referralService.PayReward(ctx, referralID)
	if err != nil {
		return fmt.Errorf("failed to pay reward: %w", err)
	}

	// 2. Отправляем уведомление о выплате
	// TODO: Получить информацию о реферале и отправить уведомление

	return nil
}

// GetReferralLink генерирует реферальную ссылку для пользователя
func (uc *ReferralUseCase) GetReferralLink(ctx context.Context, userID int64) (string, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Генерируем реферальную ссылку
	referralLink := fmt.Sprintf("https://t.me/your_bot?start=ref_%d", user.ID)

	return referralLink, nil
}

// GetReferralLeaderboard получает таблицу лидеров по рефералам
func (uc *ReferralUseCase) GetReferralLeaderboard(ctx context.Context, limit int) (*ReferralLeaderboard, error) {
	// TODO: Реализовать получение таблицы лидеров
	return nil, fmt.Errorf("referral leaderboard not implemented yet")
}

// SendReferralReminder отправляет напоминание о реферальной программе
func (uc *ReferralUseCase) SendReferralReminder(ctx context.Context, userID int64) error {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем статистику рефералов
	stats, err := uc.referralService.GetStats(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get referral stats: %w", err)
	}

	// 3. Формируем сообщение
	referralLink := fmt.Sprintf("https://t.me/your_bot?start=ref_%d", user.ID)
	message := fmt.Sprintf("💰 Реферальная программа\n\n👥 Ваши рефералы: %d\n💸 Заработано: %d дней\n\n🔗 Ваша реферальная ссылка:\n%s\n\n📢 Поделитесь с друзьями и получайте вознаграждения!",
		stats.TotalReferrals, stats.TotalRewardDays, referralLink)

	// 4. Отправляем напоминание
	if uc.notificationService != nil {
		_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Реферальная программа", message, false)
	}

	return nil
}

// ProcessReferralSignup обрабатывает регистрацию по реферальной ссылке
func (uc *ReferralUseCase) ProcessReferralSignup(ctx context.Context, referrerID int64, referredID int64) error {
	// 1. Создаем реферальную связь
	err := uc.CreateReferral(ctx, referrerID, referredID)
	if err != nil {
		return fmt.Errorf("failed to create referral: %w", err)
	}

	// 2. Создаем расширенную пробную подписку для реферала
	referred, err := uc.userService.GetByTelegramID(ctx, referredID)
	if err == nil {
		// Создаем пробную подписку на 7 дней вместо 3
		_, err = uc.subscriptionService.CreateTrial(ctx, referred.ID, 7)
		if err != nil {
			fmt.Printf("Failed to create extended trial for referral: %v\n", err)
		}
	}

	return nil
}

// ReferralStatsInfo представляет информацию о статистике рефералов
type ReferralStatsInfo struct {
	User            *domain.User          `json:"user"`
	ReferralStats   *domain.ReferralStats `json:"referral_stats"`
	ReferralLink    string                `json:"referral_link"`
	TotalEarnings   int                   `json:"total_earnings"`
	PendingEarnings int                   `json:"pending_earnings"`
	PaidEarnings    int                   `json:"paid_earnings"`
}

// ReferralLeaderboard представляет таблицу лидеров по рефералам
type ReferralLeaderboard struct {
	Users []ReferralLeader `json:"users"`
	Total int              `json:"total"`
}

// ReferralLeader представляет лидера по рефералам
type ReferralLeader struct {
	UserID         int64  `json:"user_id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	TotalReferrals int    `json:"total_referrals"`
	TotalEarnings  int    `json:"total_earnings"`
	Rank           int    `json:"rank"`
}
