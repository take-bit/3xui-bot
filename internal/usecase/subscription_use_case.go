package usecase

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/domain"
)

// SubscriptionUseCase представляет use case для работы с подписками
type SubscriptionUseCase struct {
	subscriptionService domain.SubscriptionService
	vpnUseCase          *VPNUseCase
	userService         domain.UserService
	notificationService domain.NotificationService
}

// NewSubscriptionUseCase создает новый Subscription use case
func NewSubscriptionUseCase(
	subscriptionService domain.SubscriptionService,
	vpnUseCase *VPNUseCase,
	userService domain.UserService,
	notificationService domain.NotificationService,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionService: subscriptionService,
		vpnUseCase:          vpnUseCase,
		userService:         userService,
		notificationService: notificationService,
	}
}

// CreateTrialSubscription создает пробную подписку
func (uc *SubscriptionUseCase) CreateTrialSubscription(ctx context.Context, userID int64, days int) (*domain.Subscription, error) {
	// 1. Создаем пробную подписку
	subscription, err := uc.subscriptionService.CreateTrial(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to create trial subscription: %w", err)
	}

	// 2. Создаем VPN подключение
	_, err = uc.vpnUseCase.CreateVPNConnection(ctx, userID, "")
	if err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("Failed to create VPN connection for trial: %v\n", err)
	}

	// 3. Отправляем уведомление
	if uc.notificationService != nil {
		user, err := uc.userService.GetByTelegramID(ctx, userID)
		if err == nil {
			message := fmt.Sprintf("🎁 Пробная подписка активирована!\n\n📅 Действует %d дней\n🔗 VPN подключение создано\n\n⏰ Подписка истекает: %s",
				days, subscription.EndDate.Format("02.01.2006 15:04"))
			_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Пробный период", message, false)
		}
	}

	return subscription, nil
}

// ExtendSubscription продлевает подписку
func (uc *SubscriptionUseCase) ExtendSubscription(ctx context.Context, userID int64, days int) error {
	// 1. Продлеваем подписку
	err := uc.subscriptionService.Extend(ctx, userID, days)
	if err != nil {
		return fmt.Errorf("failed to extend subscription: %w", err)
	}

	// 2. Обновляем время истечения VPN подключения
	// TODO: Добавить метод UpdateConnectionExpiry в VPNUseCase
	// err = uc.vpnUseCase.UpdateConnectionExpiry(ctx, userID, time.Now().AddDate(0, 0, days))
	// if err != nil {
	//     // Логируем ошибку, но не прерываем процесс
	//     fmt.Printf("Failed to update VPN connection expiry: %v\n", err)
	// }

	// 3. Отправляем уведомление
	if uc.notificationService != nil {
		user, err := uc.userService.GetByTelegramID(ctx, userID)
		if err == nil {
			subscription, err := uc.subscriptionService.GetActive(ctx, user.ID)
			if err == nil {
				message := fmt.Sprintf("✅ Подписка продлена на %d дней!\n\n📅 Новый срок действия: %s",
					days, subscription.EndDate.Format("02.01.2006 15:04"))
				_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Продление", message, false)
			}
		}
	}

	return nil
}

// GetActiveSubscription получает активную подписку пользователя
func (uc *SubscriptionUseCase) GetActiveSubscription(ctx context.Context, userID int64) (*SubscriptionInfo, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем активную подписку
	subscription, err := uc.subscriptionService.GetActive(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("no active subscription found: %w", err)
	}

	// 3. Вычисляем дни до истечения
	daysRemaining, err := uc.subscriptionService.GetDaysRemaining(ctx, user.ID)
	if err != nil {
		daysRemaining = 0
	}

	// 4. Получаем информацию о VPN подключении
	var vpnConnection *domain.VPNConnection
	vpnConnection, err = uc.vpnUseCase.GetVPNConnectionInfo(ctx, userID)
	if err != nil {
		vpnConnection = nil // VPN подключение не найдено
	}

	info := &SubscriptionInfo{
		Subscription:  subscription,
		DaysRemaining: daysRemaining,
		VPNConnection: vpnConnection,
		IsActive:      true,
		ExpiresAt:     subscription.EndDate,
	}

	return info, nil
}

// ExpireSubscription истекает подписку
func (uc *SubscriptionUseCase) ExpireSubscription(ctx context.Context, userID int64) error {
	// 1. Истекает подписку
	err := uc.subscriptionService.Expire(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to expire subscription: %w", err)
	}

	// 2. Удаляем VPN подключение
	err = uc.vpnUseCase.DeleteVPNConnection(ctx, userID)
	if err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("Failed to delete VPN connection after expiration: %v\n", err)
	}

	// 3. Отправляем уведомление
	if uc.notificationService != nil {
		user, err := uc.userService.GetByTelegramID(ctx, userID)
		if err == nil {
			message := "⏰ Ваша подписка истекла\n\n🔄 Продлите подписку, чтобы продолжить пользоваться VPN\n\n💳 Доступные способы оплаты в меню"
			_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Подписка истекла", message, false)
		}
	}

	return nil
}

// CheckExpiredSubscriptions проверяет и обрабатывает истекшие подписки
func (uc *SubscriptionUseCase) CheckExpiredSubscriptions(ctx context.Context) error {
	// 1. Проверяем истекшие подписки
	err := uc.subscriptionService.CheckExpired(ctx)
	if err != nil {
		return fmt.Errorf("failed to check expired subscriptions: %w", err)
	}

	// 2. Отправляем уведомления об истечении
	// TODO: Получить список истекших подписок и отправить уведомления

	return nil
}

// GetSubscriptionHistory получает историю подписок пользователя
func (uc *SubscriptionUseCase) GetSubscriptionHistory(ctx context.Context, userID int64) ([]*domain.Subscription, error) {
	// 1. Получаем пользователя
	_, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем историю подписок
	// TODO: Добавить метод GetHistory в SubscriptionService
	// return uc.subscriptionService.GetHistory(ctx, user.ID)
	return nil, fmt.Errorf("subscription history not implemented yet")
}

// GetDaysRemaining получает количество дней до истечения подписки
func (uc *SubscriptionUseCase) GetDaysRemaining(ctx context.Context, userID int64) (int, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем дни до истечения
	return uc.subscriptionService.GetDaysRemaining(ctx, user.ID)
}

// SendExpiryNotifications отправляет уведомления о скором истечении подписки
func (uc *SubscriptionUseCase) SendExpiryNotifications(ctx context.Context) error {
	// TODO: Реализовать отправку уведомлений о скором истечении
	// 1. Получить пользователей с подписками, истекающими через 1, 3, 7 дней
	// 2. Отправить уведомления

	return fmt.Errorf("expiry notifications not implemented yet")
}

// GetSubscriptionStats возвращает статистику подписок
func (uc *SubscriptionUseCase) GetSubscriptionStats(ctx context.Context) (*SubscriptionStats, error) {
	// TODO: Реализовать получение статистики подписок
	return nil, fmt.Errorf("subscription stats not implemented yet")
}

// SubscriptionInfo представляет информацию о подписке
type SubscriptionInfo struct {
	Subscription  *domain.Subscription  `json:"subscription"`
	DaysRemaining int                   `json:"days_remaining"`
	VPNConnection *domain.VPNConnection `json:"vpn_connection,omitempty"`
	IsActive      bool                  `json:"is_active"`
	ExpiresAt     time.Time             `json:"expires_at"`
}

// SubscriptionStats представляет статистику подписок
type SubscriptionStats struct {
	TotalSubscriptions   int     `json:"total_subscriptions"`
	ActiveSubscriptions  int     `json:"active_subscriptions"`
	ExpiredSubscriptions int     `json:"expired_subscriptions"`
	TrialSubscriptions   int     `json:"trial_subscriptions"`
	AverageDuration      float64 `json:"average_duration"`
	RenewalRate          float64 `json:"renewal_rate"`
}
