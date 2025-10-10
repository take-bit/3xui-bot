package usecase

import (
	"3xui-bot/internal/ports"
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
)

// VPNUseCase use case для работы с VPN
type VPNUseCase struct {
	vpnRepo     ports.VPNRepo
	marzbanRepo ports.Marzban
	subRepo     ports.SubscriptionRepo
	planRepo    ports.PlanRepo
}

// NewVPNUseCase создает новый use case для VPN
func NewVPNUseCase(
	vpnRepo ports.VPNRepo,
	marzbanRepo ports.Marzban,
	subRepo ports.SubscriptionRepo,
	planRepo ports.PlanRepo,
) *VPNUseCase {
	return &VPNUseCase{
		vpnRepo:     vpnRepo,
		marzbanRepo: marzbanRepo,
		subRepo:     subRepo,
		planRepo:    planRepo,
	}
}

// CreateVPNConnection создает новое VPN подключение
func (uc *VPNUseCase) CreateVPNConnection(ctx context.Context, connection *core.VPNConnection) error {
	return uc.vpnRepo.CreateVPNConnection(ctx, connection)
}

// GetUserVPNConnections получает все VPN подключения пользователя
func (uc *VPNUseCase) GetUserVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {
	return uc.vpnRepo.GetVPNConnectionsByTelegramUserID(ctx, telegramUserID)
}

// GetVPNConnection получает VPN подключение по ID
func (uc *VPNUseCase) GetVPNConnection(ctx context.Context, id string) (*core.VPNConnection, error) {
	return uc.vpnRepo.GetVPNConnectionByID(ctx, id)
}

// GetVPNConnectionByMarzbanUsername получает VPN подключение по Marzban username
func (uc *VPNUseCase) GetVPNConnectionByMarzbanUsername(ctx context.Context, username string) (*core.VPNConnection, error) {
	return uc.vpnRepo.GetVPNConnectionByMarzbanUsername(ctx, username)
}

// UpdateVPNConnectionName обновляет имя VPN подключения
func (uc *VPNUseCase) UpdateVPNConnectionName(ctx context.Context, id, name string) error {
	return uc.vpnRepo.UpdateVPNConnectionName(ctx, id, name)
}

// DeleteVPNConnection удаляет VPN подключение
func (uc *VPNUseCase) DeleteVPNConnection(ctx context.Context, id string) error {
	return uc.vpnRepo.DeleteVPNConnection(ctx, id)
}

// DeleteVPNConnectionByMarzbanUsername удаляет VPN подключение по Marzban username
func (uc *VPNUseCase) DeleteVPNConnectionByMarzbanUsername(ctx context.Context, username string) error {
	return uc.vpnRepo.DeleteVPNConnectionByMarzbanUsername(ctx, username)
}

// GetActiveVPNConnections получает все активные VPN подключения пользователя
func (uc *VPNUseCase) GetActiveVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {
	return uc.vpnRepo.GetActiveVPNConnections(ctx, telegramUserID)
}

// CreateVPNForSubscription создает VPN подключение для подписки (бизнес-логика)
func (uc *VPNUseCase) CreateVPNForSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.VPNConnection, error) {
	// Получаем подписку
	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Получаем план для определения лимитов
	plan, err := uc.planRepo.GetPlanByID(ctx, sub.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Генерируем уникальный username для Marzban
	marzbanUsername := fmt.Sprintf("user_%d_%s", userID, id.GenerateShort())

	// Создаем пользователя в Marzban
	dataLimit := int64(100 * 1024 * 1024 * 1024) // 100 GB по умолчанию
	marzbanUser := &core.MarzbanUserData{
		Username:  marzbanUsername,
		DataLimit: &dataLimit,
		Status:    "active",
		Note:      fmt.Sprintf("User %d - %s", userID, plan.Name),
	}

	// Создаем в Marzban
	_, err = uc.marzbanRepo.CreateUser(ctx, marzbanUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in Marzban: %w", err)
	}

	// Создаем запись в локальной БД
	vpnConn := &core.VPNConnection{
		ID:              id.Generate(),
		TelegramUserID:  userID,
		MarzbanUsername: marzbanUsername,
		Name:            fmt.Sprintf("VPN - %s", plan.Name),
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.vpnRepo.CreateVPNConnection(ctx, vpnConn); err != nil {
		// Откатываем создание в Marzban если не удалось сохранить локально
		_ = uc.marzbanRepo.DeleteUser(ctx, marzbanUsername)
		return nil, fmt.Errorf("failed to create VPN connection: %w", err)
	}

	return vpnConn, nil
}

// GetUserVPNWithStats получает VPN подключения пользователя с данными из Marzban
func (uc *VPNUseCase) GetUserVPNWithStats(ctx context.Context, userID int64) ([]*core.VPNConnection, error) {
	// Получаем локальные VPN подключения
	connections, err := uc.vpnRepo.GetVPNConnectionsByTelegramUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN connections: %w", err)
	}

	// Обогащаем данными из Marzban
	for _, conn := range connections {
		marzbanData, err := uc.marzbanRepo.GetUser(ctx, conn.MarzbanUsername)
		if err != nil {
			// Если не удалось получить из Marzban, помечаем как неактивный
			conn.IsActive = false
			continue
		}

		// Обновляем данные из Marzban
		conn.ExpireAt = marzbanData.ExpireAt()
		conn.DataLimitBytes = marzbanData.DataLimit
		conn.DataUsedBytes = marzbanData.DataUsed
		conn.Status = marzbanData.Status
		conn.ProtocolConfig = marzbanData.Proxies
	}

	return connections, nil
}

// GetVPNConnectionWithStats получает конкретное VPN подключение с данными из Marzban
func (uc *VPNUseCase) GetVPNConnectionWithStats(ctx context.Context, vpnID string) (*core.VPNConnection, error) {
	connection, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VPN connection: %w", err)
	}

	// Получаем данные из Marzban
	marzbanData, err := uc.marzbanRepo.GetUser(ctx, connection.MarzbanUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get Marzban user data: %w", err)
	}

	// Обновляем данные из Marzban
	connection.ExpireAt = marzbanData.ExpireAt()
	connection.DataLimitBytes = marzbanData.DataLimit
	connection.DataUsedBytes = marzbanData.DataUsed
	connection.Status = marzbanData.Status
	connection.ProtocolConfig = marzbanData.Proxies

	return connection, nil
}

// DeleteVPNConnectionFull удаляет VPN подключение полностью (из БД и Marzban)
func (uc *VPNUseCase) DeleteVPNConnectionFull(ctx context.Context, vpnID string) error {
	// Получаем подключение
	conn, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {
		return fmt.Errorf("failed to get VPN connection: %w", err)
	}

	// Удаляем из Marzban
	if err := uc.marzbanRepo.DeleteUser(ctx, conn.MarzbanUsername); err != nil {
		return fmt.Errorf("failed to delete user from Marzban: %w", err)
	}

	// Удаляем из локальной БД
	if err := uc.vpnRepo.DeleteVPNConnection(ctx, vpnID); err != nil {
		return fmt.Errorf("failed to delete VPN connection: %w", err)
	}

	return nil
}

// SyncVPNStatus синхронизирует статус VPN с Marzban
func (uc *VPNUseCase) SyncVPNStatus(ctx context.Context, vpnID string) error {
	conn, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {
		return fmt.Errorf("failed to get VPN connection: %w", err)
	}

	marzbanData, err := uc.marzbanRepo.GetUser(ctx, conn.MarzbanUsername)
	if err != nil {
		return fmt.Errorf("failed to get Marzban data: %w", err)
	}

	// Обновляем активность в зависимости от статуса в Marzban
	isActive := marzbanData.Status == "active"
	if conn.IsActive != isActive {
		// TODO: Обновить IsActive в БД если добавим такой метод
		_ = isActive
	}

	return nil
}

// DeactivateExpiredVPNs деактивирует истекшие VPN подключения
func (uc *VPNUseCase) DeactivateExpiredVPNs(ctx context.Context) error {
	// TODO: Реализовать получение всех активных VPN и проверку их истечения
	// Сейчас это делается через scheduled jobs
	return nil
}
