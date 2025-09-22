package usecase

import (
	"context"
	"fmt"

	"3xui-bot/internal/domain"
)

// UseCaseManager управляет всеми Use Case
type UseCaseManager struct {
	// Use Cases
	VPNUseCase          *VPNUseCase
	PaymentUseCase      *PaymentUseCase
	UserUseCase         *UserUseCase
	SubscriptionUseCase *SubscriptionUseCase
	PromocodeUseCase    *PromocodeUseCase
	ReferralUseCase     *ReferralUseCase
	NotificationUseCase *NotificationUseCase
	ServerUseCase       *ServerUseCase
}

// NewUseCaseManager создает новый менеджер Use Case
func NewUseCaseManager(
	// Services
	vpnService domain.VPNService,
	paymentService domain.PaymentService,
	userService domain.UserService,
	subscriptionService domain.SubscriptionService,
	promocodeService domain.PromocodeService,
	referralService domain.ReferralService,
	notificationService domain.NotificationService,
	serverService domain.ServerService,
) *UseCaseManager {
	// Создаем Use Cases
	notificationUseCase := NewNotificationUseCase(notificationService, userService)
	serverUseCase := NewServerUseCase(serverService, notificationService)
	vpnUseCase := NewVPNUseCase(vpnService, serverService, userService, subscriptionService, notificationService)
	paymentUseCase := NewPaymentUseCase(paymentService, subscriptionService, vpnUseCase, userService, notificationService)
	userUseCase := NewUserUseCase(userService, subscriptionService, referralService, notificationService)
	subscriptionUseCase := NewSubscriptionUseCase(subscriptionService, vpnUseCase, userService, notificationService)
	promocodeUseCase := NewPromocodeUseCase(promocodeService, subscriptionService, userService, notificationService)
	referralUseCase := NewReferralUseCase(referralService, userService, subscriptionService, notificationService)

	return &UseCaseManager{
		VPNUseCase:          vpnUseCase,
		PaymentUseCase:      paymentUseCase,
		UserUseCase:         userUseCase,
		SubscriptionUseCase: subscriptionUseCase,
		PromocodeUseCase:    promocodeUseCase,
		ReferralUseCase:     referralUseCase,
		NotificationUseCase: notificationUseCase,
		ServerUseCase:       serverUseCase,
	}
}

// Initialize инициализирует все Use Cases
func (m *UseCaseManager) Initialize(ctx context.Context) error {
	// Запускаем мониторинг серверов
	err := m.ServerUseCase.StartHealthMonitoring(ctx)
	if err != nil {
		return fmt.Errorf("failed to start server health monitoring: %w", err)
	}

	// TODO: Добавить другие инициализации

	return nil
}

// Shutdown корректно завершает работу всех Use Cases
func (m *UseCaseManager) Shutdown() error {
	// Останавливаем мониторинг серверов
	err := m.ServerUseCase.StopHealthMonitoring()
	if err != nil {
		return fmt.Errorf("failed to stop server health monitoring: %w", err)
	}

	// TODO: Добавить другие завершения

	return nil
}

// GetVPNUseCase возвращает VPN Use Case
func (m *UseCaseManager) GetVPNUseCase() *VPNUseCase {
	return m.VPNUseCase
}

// GetPaymentUseCase возвращает Payment Use Case
func (m *UseCaseManager) GetPaymentUseCase() *PaymentUseCase {
	return m.PaymentUseCase
}

// GetUserUseCase возвращает User Use Case
func (m *UseCaseManager) GetUserUseCase() *UserUseCase {
	return m.UserUseCase
}

// GetSubscriptionUseCase возвращает Subscription Use Case
func (m *UseCaseManager) GetSubscriptionUseCase() *SubscriptionUseCase {
	return m.SubscriptionUseCase
}

// GetPromocodeUseCase возвращает Promocode Use Case
func (m *UseCaseManager) GetPromocodeUseCase() *PromocodeUseCase {
	return m.PromocodeUseCase
}

// GetReferralUseCase возвращает Referral Use Case
func (m *UseCaseManager) GetReferralUseCase() *ReferralUseCase {
	return m.ReferralUseCase
}

// GetNotificationUseCase возвращает Notification Use Case
func (m *UseCaseManager) GetNotificationUseCase() *NotificationUseCase {
	return m.NotificationUseCase
}

// GetServerUseCase возвращает Server Use Case
func (m *UseCaseManager) GetServerUseCase() *ServerUseCase {
	return m.ServerUseCase
}

// ProcessUserRegistration обрабатывает регистрацию пользователя
func (m *UseCaseManager) ProcessUserRegistration(ctx context.Context, telegramID int64, username, firstName, lastName, languageCode string) (*domain.User, error) {
	// 1. Регистрируем пользователя
	user, err := m.UserUseCase.RegisterUser(ctx, telegramID, username, firstName, lastName, languageCode)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// 2. Создаем VPN подключение
	_, err = m.VPNUseCase.CreateVPNConnection(ctx, telegramID, "")
	if err != nil {
		// Логируем ошибку, но не прерываем регистрацию
		fmt.Printf("Failed to create VPN connection during registration: %v\n", err)
	}

	return user, nil
}

// ProcessReferralRegistration обрабатывает регистрацию по реферальной ссылке
func (m *UseCaseManager) ProcessReferralRegistration(ctx context.Context, referrerID, referredID int64, username, firstName, lastName, languageCode string) (*domain.User, error) {
	// 1. Регистрируем пользователя
	user, err := m.UserUseCase.RegisterUser(ctx, referredID, username, firstName, lastName, languageCode)
	if err != nil {
		return nil, fmt.Errorf("failed to register referred user: %w", err)
	}

	// 2. Создаем реферальную связь
	err = m.ReferralUseCase.ProcessReferralSignup(ctx, referrerID, referredID)
	if err != nil {
		// Логируем ошибку, но не прерываем регистрацию
		fmt.Printf("Failed to process referral signup: %v\n", err)
	}

	return user, nil
}

// ProcessPaymentWebhook обрабатывает webhook от платежной системы
func (m *UseCaseManager) ProcessPaymentWebhook(ctx context.Context, externalID string, status domain.PaymentStatus) error {
	return m.PaymentUseCase.ProcessPaymentWebhook(ctx, externalID, status)
}

// ProcessPromocodeApplication обрабатывает применение промокода
func (m *UseCaseManager) ProcessPromocodeApplication(ctx context.Context, userID int64, code string) (*PromocodeResult, error) {
	return m.PromocodeUseCase.ApplyPromocode(ctx, userID, code)
}

// GetUserDashboard возвращает данные для дашборда пользователя
func (m *UseCaseManager) GetUserDashboard(ctx context.Context, userID int64) (*UserDashboard, error) {
	// 1. Получаем профиль пользователя
	profile, err := m.UserUseCase.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// 2. Получаем информацию о VPN подключении
	var vpnConnection *domain.VPNConnection
	vpnConnection, err = m.VPNUseCase.GetVPNConnectionInfo(ctx, userID)
	if err != nil {
		vpnConnection = nil // VPN подключение не найдено
	}

	// 3. Получаем статистику рефералов
	var referralStats *ReferralStatsInfo
	referralStats, err = m.ReferralUseCase.GetReferralStats(ctx, userID)
	if err != nil {
		referralStats = nil // Реферальная статистика недоступна
	}

	dashboard := &UserDashboard{
		User:          profile.User,
		Subscription:  profile.Subscription,
		DaysRemaining: profile.DaysRemaining,
		VPNConnection: vpnConnection,
		ReferralStats: referralStats,
		IsBlocked:     profile.IsBlocked,
	}

	return dashboard, nil
}

// UserDashboard представляет данные для дашборда пользователя
type UserDashboard struct {
	User          *domain.User          `json:"user"`
	Subscription  *domain.Subscription  `json:"subscription,omitempty"`
	DaysRemaining int                   `json:"days_remaining"`
	VPNConnection *domain.VPNConnection `json:"vpn_connection,omitempty"`
	ReferralStats *ReferralStatsInfo    `json:"referral_stats,omitempty"`
	IsBlocked     bool                  `json:"is_blocked"`
}
