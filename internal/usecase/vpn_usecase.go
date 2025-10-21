package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
	"3xui-bot/internal/ports"
)

type VPNUseCase struct {
	vpnRepo     ports.VPNRepo
	marzbanRepo ports.Marzban
	subRepo     ports.SubscriptionRepo
	planRepo    ports.PlanRepo
}

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

func (uc *VPNUseCase) CreateVPNConnection(ctx context.Context, connection *core.VPNConnection) error {

	return uc.vpnRepo.CreateVPNConnection(ctx, connection)
}

func (uc *VPNUseCase) GetUserVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {

	return uc.vpnRepo.GetVPNConnectionsByTelegramUserID(ctx, telegramUserID)
}

func (uc *VPNUseCase) GetVPNConnection(ctx context.Context, id string) (*core.VPNConnection, error) {

	return uc.vpnRepo.GetVPNConnectionByID(ctx, id)
}

func (uc *VPNUseCase) GetVPNConnectionByMarzbanUsername(ctx context.Context, username string) (*core.VPNConnection, error) {

	return uc.vpnRepo.GetVPNConnectionByMarzbanUsername(ctx, username)
}

func (uc *VPNUseCase) GetVPNConnectionsBySubscription(ctx context.Context, subscriptionID string) ([]*core.VPNConnection, error) {

	return uc.vpnRepo.GetVPNConnectionsBySubscriptionID(ctx, subscriptionID)
}

func (uc *VPNUseCase) UpdateVPNConnectionName(ctx context.Context, id, name string) error {

	return uc.vpnRepo.UpdateVPNConnectionName(ctx, id, name)
}

func (uc *VPNUseCase) DeleteVPNConnection(ctx context.Context, id string) error {

	return uc.vpnRepo.DeleteVPNConnection(ctx, id)
}

func (uc *VPNUseCase) DeleteVPNConnectionByMarzbanUsername(ctx context.Context, username string) error {

	return uc.vpnRepo.DeleteVPNConnectionByMarzbanUsername(ctx, username)
}

func (uc *VPNUseCase) GetActiveVPNConnections(ctx context.Context, telegramUserID int64) ([]*core.VPNConnection, error) {

	return uc.vpnRepo.GetActiveVPNConnections(ctx, telegramUserID)
}

func (uc *VPNUseCase) CreateVPNForSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.VPNConnection, error) {
	slog.Info("Creating VPN for subscription", "user_id", userID, "subscription_id", subscriptionID)

	sub, err := uc.subRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		slog.Error("Failed to get subscription", "subscription_id", subscriptionID, "error", err)

		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	plan, err := uc.planRepo.GetPlanByID(ctx, sub.PlanID)
	if err != nil {

		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	inbounds, err := uc.marzbanRepo.GetInbounds(ctx)
	if err != nil {
		slog.Warn("Failed to get inbounds from Marzban, will try without specific inbounds", "error", err)
		inbounds = nil
	}

	userInbounds := uc.buildUserInbounds(inbounds)
	slog.Debug("Built user inbounds", "inbounds", userInbounds)

	marzbanUsername := fmt.Sprintf("user_%d_%s", userID, id.GenerateShort())

	var expireTimestamp *int64
	if !sub.EndDate.IsZero() {
		timestamp := sub.EndDate.Unix()
		expireTimestamp = &timestamp
	}

	dataLimit := int64(100 * 1024 * 1024 * 1024)

	marzbanUser := &core.MarzbanUserData{
		Username:  marzbanUsername,
		DataLimit: &dataLimit,
		Expire:    expireTimestamp,
		Status:    "active",
		Note:      fmt.Sprintf("User %d - %s", userID, plan.Name),
		Proxies: map[string]interface{}{
			"vless": map[string]interface{}{},
		},
		Inbounds: userInbounds,
	}

	_, err = uc.marzbanRepo.CreateUser(ctx, marzbanUser)
	if err != nil {
		slog.Error("Failed to create user in Marzban", "username", marzbanUsername, "error", err)

		return nil, fmt.Errorf("failed to create user in Marzban: %w", err)
	}

	slog.Debug("Marzban user created", "username", marzbanUsername)

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
		_ = uc.marzbanRepo.DeleteUser(ctx, marzbanUsername)

		return nil, fmt.Errorf("failed to create VPN connection: %w", err)
	}

	return vpnConn, nil
}

func (uc *VPNUseCase) GetUserVPNWithStats(ctx context.Context, userID int64) ([]*core.VPNConnection, error) {
	connections, err := uc.vpnRepo.GetVPNConnectionsByTelegramUserID(ctx, userID)
	if err != nil {

		return nil, fmt.Errorf("failed to get VPN connections: %w", err)
	}

	for _, conn := range connections {
		marzbanData, err := uc.marzbanRepo.GetUser(ctx, conn.MarzbanUsername)
		if err != nil {
			conn.IsActive = false
			continue
		}

		conn.ExpireAt = marzbanData.ExpireAt()
		conn.DataLimitBytes = marzbanData.DataLimit
		conn.DataUsedBytes = marzbanData.DataUsed
		conn.Status = marzbanData.Status
		conn.ProtocolConfig = marzbanData.Proxies
	}

	return connections, nil
}

func (uc *VPNUseCase) GetVPNConnectionWithStats(ctx context.Context, vpnID string) (*core.VPNConnection, error) {
	connection, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {

		return nil, fmt.Errorf("failed to get VPN connection: %w", err)
	}

	marzbanData, err := uc.marzbanRepo.GetUser(ctx, connection.MarzbanUsername)
	if err != nil {

		return nil, fmt.Errorf("failed to get Marzban user data: %w", err)
	}

	connection.ExpireAt = marzbanData.ExpireAt()
	connection.DataLimitBytes = marzbanData.DataLimit
	connection.DataUsedBytes = marzbanData.DataUsed
	connection.Status = marzbanData.Status
	connection.ProtocolConfig = marzbanData.Proxies

	return connection, nil
}

func (uc *VPNUseCase) DeleteVPNConnectionFull(ctx context.Context, vpnID string) error {
	conn, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {

		return fmt.Errorf("failed to get VPN connection: %w", err)
	}

	if err := uc.marzbanRepo.DeleteUser(ctx, conn.MarzbanUsername); err != nil {

		return fmt.Errorf("failed to delete user from Marzban: %w", err)
	}

	if err := uc.vpnRepo.DeleteVPNConnection(ctx, vpnID); err != nil {

		return fmt.Errorf("failed to delete VPN connection: %w", err)
	}

	return nil
}

func (uc *VPNUseCase) SyncVPNStatus(ctx context.Context, vpnID string) error {
	conn, err := uc.vpnRepo.GetVPNConnectionByID(ctx, vpnID)
	if err != nil {

		return fmt.Errorf("failed to get VPN connection: %w", err)
	}

	marzbanData, err := uc.marzbanRepo.GetUser(ctx, conn.MarzbanUsername)
	if err != nil {

		return fmt.Errorf("failed to get Marzban data: %w", err)
	}

	isActive := marzbanData.Status == "active"
	if conn.IsActive != isActive {
		_ = isActive
	}

	return nil
}

func (uc *VPNUseCase) buildUserInbounds(inbounds []map[string]interface{}) map[string][]string {
	if len(inbounds) == 0 {

		return make(map[string][]string)
	}

	inboundsByProtocol := make(map[string][]string)

	for _, inbound := range inbounds {
		tag, hasTag := inbound["tag"].(string)
		protocol, hasProtocol := inbound["protocol"].(string)

		if hasTag && hasProtocol {
			inboundsByProtocol[protocol] = append(inboundsByProtocol[protocol], tag)
		} else if hasTag {
			inboundsByProtocol["vless"] = append(inboundsByProtocol["vless"], tag)
		}
	}

	slog.Debug("Inbounds grouped by protocol", "groups", inboundsByProtocol)

	return inboundsByProtocol
}

func (uc *VPNUseCase) DeactivateExpiredVPNs(ctx context.Context) error {

	return nil
}
