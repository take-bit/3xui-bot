package usecase

import (
	"context"
	"fmt"

	"3xui-bot/internal/domain"
)

// UserUseCase представляет use case для работы с пользователями
type UserUseCase struct {
	userService         domain.UserService
	subscriptionService domain.SubscriptionService
	referralService     domain.ReferralService
	notificationService domain.NotificationService
}

// NewUserUseCase создает новый User use case
func NewUserUseCase(
	userService domain.UserService,
	subscriptionService domain.SubscriptionService,
	referralService domain.ReferralService,
	notificationService domain.NotificationService,
) *UserUseCase {
	return &UserUseCase{
		userService:         userService,
		subscriptionService: subscriptionService,
		referralService:     referralService,
		notificationService: notificationService,
	}
}

// RegisterUser регистрирует нового пользователя
func (uc *UserUseCase) RegisterUser(ctx context.Context, telegramID int64, username, firstName, lastName, languageCode string) (*domain.User, error) {
	// 1. Создаем или обновляем пользователя
	user, err := uc.userService.CreateOrUpdate(ctx, telegramID, username, firstName, lastName, languageCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update user: %w", err)
	}

	// 2. Создаем пробную подписку, если у пользователя ее нет
	_, err = uc.subscriptionService.GetActive(ctx, user.ID)
	if err != nil {
		// Если активной подписки нет, создаем пробную
		_, err = uc.subscriptionService.CreateTrial(ctx, user.ID, 3) // 3 дня пробного периода
		if err != nil {
			// Логируем ошибку, но не прерываем регистрацию
			fmt.Printf("Failed to create trial subscription for user %d: %v\n", user.ID, err)
		}
	}

	// 3. Отправляем приветственное сообщение
	if uc.notificationService != nil {
		message := fmt.Sprintf("👋 Добро пожаловать, %s!\n\n🎁 Вам предоставлен пробный период на 3 дня\n\n🚀 Начните пользоваться VPN прямо сейчас!", firstName)
		_ = uc.notificationService.SendToUser(ctx, telegramID, "Добро пожаловать!", message, false)
	}

	return user, nil
}

// GetUserProfile получает профиль пользователя
func (uc *UserUseCase) GetUserProfile(ctx context.Context, telegramID int64) (*UserProfile, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем активную подписку
	subscription, err := uc.subscriptionService.GetActive(ctx, user.ID)
	if err != nil {
		subscription = nil // У пользователя нет активной подписки
	}

	// 3. Получаем статистику рефералов
	var referralStats *domain.ReferralStats
	if uc.referralService != nil {
		referralStats, err = uc.referralService.GetStats(ctx, user.ID)
		if err != nil {
			referralStats = nil // Игнорируем ошибки реферальной системы
		}
	}

	// 4. Вычисляем дни до истечения подписки
	daysRemaining := 0
	if subscription != nil {
		daysRemaining, _ = uc.subscriptionService.GetDaysRemaining(ctx, user.ID)
	}

	profile := &UserProfile{
		User:             user,
		Subscription:     subscription,
		DaysRemaining:    daysRemaining,
		ReferralStats:    referralStats,
		IsBlocked:        user.IsBlocked,
		RegistrationDate: user.CreatedAt.Format("02.01.2006"),
	}

	return profile, nil
}

// BlockUser блокирует пользователя
func (uc *UserUseCase) BlockUser(ctx context.Context, telegramID int64, reason string) error {
	// 1. Блокируем пользователя
	err := uc.userService.Block(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	// 2. Отправляем уведомление о блокировке
	if uc.notificationService != nil {
		message := fmt.Sprintf("🚫 Ваш аккаунт заблокирован\n\n📋 Причина: %s\n\n❓ По вопросам обращайтесь в поддержку", reason)
		_ = uc.notificationService.SendToUser(ctx, telegramID, "Блокировка", message, false)
	}

	return nil
}

// UnblockUser разблокирует пользователя
func (uc *UserUseCase) UnblockUser(ctx context.Context, telegramID int64) error {
	// 1. Разблокируем пользователя
	err := uc.userService.Unblock(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	// 2. Отправляем уведомление о разблокировке
	if uc.notificationService != nil {
		message := "✅ Ваш аккаунт разблокирован!\n\n🎉 Добро пожаловать обратно!"
		_ = uc.notificationService.SendToUser(ctx, telegramID, "Разблокировка", message, false)
	}

	return nil
}

// GetUserDisplayName получает отображаемое имя пользователя
func (uc *UserUseCase) GetUserDisplayName(ctx context.Context, telegramID int64) (string, error) {
	return uc.userService.GetDisplayName(ctx, telegramID)
}

// UpdateUserLanguage обновляет язык пользователя
func (uc *UserUseCase) UpdateUserLanguage(ctx context.Context, telegramID int64, languageCode string) error {
	// 1. Получаем пользователя
	_, err := uc.userService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Обновляем язык
	// TODO: Добавить метод UpdateLanguage в UserService
	// err = uc.userService.UpdateLanguage(ctx, user.ID, languageCode)
	// if err != nil {
	//     return fmt.Errorf("failed to update language: %w", err)
	// }

	// 3. Отправляем подтверждение
	if uc.notificationService != nil {
		message := fmt.Sprintf("🌐 Язык изменен на: %s", languageCode)
		_ = uc.notificationService.SendToUser(ctx, telegramID, "Настройки", message, false)
	}

	return fmt.Errorf("language update not implemented yet")
}

// GetUserStats возвращает статистику пользователя
func (uc *UserUseCase) GetUserStats(ctx context.Context, telegramID int64) (*UserStats, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем статистику рефералов
	var referralStats *domain.ReferralStats
	if uc.referralService != nil {
		referralStats, err = uc.referralService.GetStats(ctx, user.ID)
		if err != nil {
			referralStats = nil
		}
	}

	// 3. Получаем дни до истечения подписки
	daysRemaining := 0
	_, err = uc.subscriptionService.GetActive(ctx, user.ID)
	if err == nil {
		daysRemaining, _ = uc.subscriptionService.GetDaysRemaining(ctx, user.ID)
	}

	stats := &UserStats{
		UserID:           user.ID,
		TelegramID:       user.TelegramID,
		Username:         user.Username,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		LanguageCode:     user.LanguageCode,
		IsBlocked:        user.IsBlocked,
		RegistrationDate: user.CreatedAt.Format("02.01.2006"),
		DaysRemaining:    daysRemaining,
		ReferralStats:    referralStats,
	}

	return stats, nil
}

// SendUserNotification отправляет уведомление пользователю
func (uc *UserUseCase) SendUserNotification(ctx context.Context, telegramID int64, title, message string, isHTML bool) error {
	if uc.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	return uc.notificationService.SendToUser(ctx, telegramID, title, message, isHTML)
}

// UserProfile представляет профиль пользователя
type UserProfile struct {
	User             *domain.User          `json:"user"`
	Subscription     *domain.Subscription  `json:"subscription,omitempty"`
	DaysRemaining    int                   `json:"days_remaining"`
	ReferralStats    *domain.ReferralStats `json:"referral_stats,omitempty"`
	IsBlocked        bool                  `json:"is_blocked"`
	RegistrationDate string                `json:"registration_date"`
}

// UserStats представляет статистику пользователя
type UserStats struct {
	UserID           int64                 `json:"user_id"`
	TelegramID       int64                 `json:"telegram_id"`
	Username         string                `json:"username"`
	FirstName        string                `json:"first_name"`
	LastName         string                `json:"last_name"`
	LanguageCode     string                `json:"language_code"`
	IsBlocked        bool                  `json:"is_blocked"`
	RegistrationDate string                `json:"registration_date"`
	DaysRemaining    int                   `json:"days_remaining"`
	ReferralStats    *domain.ReferralStats `json:"referral_stats,omitempty"`
}
