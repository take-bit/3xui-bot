package usecase

import (
	"context"
	"fmt"

	"3xui-bot/internal/domain"
)

// VPNUseCase представляет use case для работы с VPN подключениями
type VPNUseCase struct {
	vpnService          domain.VPNService
	serverService       domain.ServerService
	userService         domain.UserService
	subscriptionService domain.SubscriptionService
	notificationService domain.NotificationService
}

// NewVPNUseCase создает новый VPN use case
func NewVPNUseCase(
	vpnService domain.VPNService,
	serverService domain.ServerService,
	userService domain.UserService,
	subscriptionService domain.SubscriptionService,
	notificationService domain.NotificationService,
) *VPNUseCase {
	return &VPNUseCase{
		vpnService:          vpnService,
		serverService:       serverService,
		userService:         userService,
		subscriptionService: subscriptionService,
		notificationService: notificationService,
	}
}

// CreateVPNConnection создает VPN подключение для пользователя
func (uc *VPNUseCase) CreateVPNConnection(ctx context.Context, userID int64, region string) (*domain.VPNConnection, error) {
	// 1. Проверяем и получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Проверяем активную подписку
	_, err = uc.subscriptionService.GetActive(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("no active subscription found: %w", err)
	}

	// 3. Создаем VPN подключение через VPN Service
	connection, err := uc.vpnService.CreateConnectionWithRegion(ctx, user.ID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPN connection: %w", err)
	}

	// 4. Получаем информацию о сервере для уведомления
	serverHealth, err := uc.serverService.GetServerHealth(ctx, connection.ServerID)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Failed to get server health: %v\n", err)
	}

	// 5. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		message := fmt.Sprintf("✅ VPN подключение создано!\n\n🌍 Сервер: %s\n📅 Действует до: %s",
			connection.ConfigURL,
			connection.ExpiresAt.Format("02.01.2006 15:04"))

		if serverHealth != nil && !serverHealth.IsHealthy {
			message += "\n⚠️ Сервер временно недоступен, но подключение создано"
		}

		_ = uc.notificationService.SendToUser(ctx, userID, "VPN Подключение", message, false)
	}

	return connection, nil
}

// GetVPNConnectionInfo получает информацию о VPN подключении
func (uc *VPNUseCase) GetVPNConnectionInfo(ctx context.Context, userID int64) (*domain.VPNConnection, error) {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем информацию о подключении
	connection, err := uc.vpnService.GetConnectionInfo(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN connection info: %w", err)
	}

	// 3. Получаем дополнительную информацию о сервере
	serverHealth, err := uc.serverService.GetServerHealth(ctx, connection.ServerID)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Failed to get server health: %v\n", err)
	}

	// TODO: Добавить информацию о здоровье сервера в ответ
	_ = serverHealth

	return connection, nil
}

// DeleteVPNConnection удаляет VPN подключение
func (uc *VPNUseCase) DeleteVPNConnection(ctx context.Context, userID int64) error {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Удаляем подключение
	err = uc.vpnService.DeleteConnection(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete VPN connection: %w", err)
	}

	// 3. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		message := "❌ VPN подключение удалено"
		_ = uc.notificationService.SendToUser(ctx, userID, "VPN Подключение", message, false)
	}

	return nil
}

// ExtendVPNConnection продлевает VPN подключение
func (uc *VPNUseCase) ExtendVPNConnection(ctx context.Context, userID int64, days int) error {
	// 1. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Продлеваем подписку
	err = uc.subscriptionService.Extend(ctx, user.ID, days)
	if err != nil {
		return fmt.Errorf("failed to extend subscription: %w", err)
	}

	// 3. Обновляем время истечения VPN подключения
	subscription, err := uc.subscriptionService.GetActive(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get updated subscription: %w", err)
	}

	err = uc.vpnService.UpdateConnectionExpiry(ctx, user.ID, subscription.EndDate)
	if err != nil {
		return fmt.Errorf("failed to update VPN connection expiry: %w", err)
	}

	// 4. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		message := fmt.Sprintf("✅ VPN подключение продлено на %d дней!\n\n📅 Новый срок действия: %s",
			days,
			subscription.EndDate.Format("02.01.2006 15:04"))
		_ = uc.notificationService.SendToUser(ctx, userID, "VPN Подключение", message, false)
	}

	return nil
}

// GetAvailableServers возвращает список доступных серверов
func (uc *VPNUseCase) GetAvailableServers(ctx context.Context) ([]*domain.Server, error) {
	return uc.serverService.GetAvailableServers(ctx)
}

// GetServerHealth возвращает состояние сервера
func (uc *VPNUseCase) GetServerHealth(ctx context.Context, serverID int64) (*domain.ServerHealth, error) {
	return uc.serverService.GetServerHealth(ctx, serverID)
}

// GetAllServersHealth возвращает состояние всех серверов
func (uc *VPNUseCase) GetAllServersHealth(ctx context.Context) (map[int64]*domain.ServerHealth, error) {
	return uc.serverService.GetAllServersHealth(ctx)
}

// GetServerStats возвращает статистику сервера
func (uc *VPNUseCase) GetServerStats(ctx context.Context, serverID int64) (*domain.ServerStats, error) {
	return uc.serverService.GetServerStats(ctx, serverID)
}

// GetAllServersStats возвращает статистику всех серверов
func (uc *VPNUseCase) GetAllServersStats(ctx context.Context) (map[int64]*domain.ServerStats, error) {
	return uc.serverService.GetAllServersStats(ctx)
}

// StartHealthMonitoring запускает мониторинг состояния серверов
func (uc *VPNUseCase) StartHealthMonitoring(ctx context.Context) error {
	return uc.serverService.StartHealthMonitoring(ctx)
}

// StopHealthMonitoring останавливает мониторинг состояния серверов
func (uc *VPNUseCase) StopHealthMonitoring() error {
	return uc.serverService.StopHealthMonitoring()
}

// IsHealthMonitoringActive проверяет, активен ли мониторинг
func (uc *VPNUseCase) IsHealthMonitoringActive() bool {
	return uc.serverService.IsHealthMonitoringActive()
}
