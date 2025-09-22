package server

import (
	"context"

	"3xui-bot/internal/domain"
)

// ServerService реализует domain.ServerService
type ServerService struct {
	serverManager domain.XUIClientManager
}

// NewServerService создает новый сервис серверов
func NewServerService(serverManager domain.XUIClientManager) *ServerService {
	return &ServerService{
		serverManager: serverManager,
	}
}

// GetAvailableServers возвращает список доступных серверов
func (s *ServerService) GetAvailableServers(ctx context.Context) ([]*domain.Server, error) {
	return s.serverManager.GetAvailableServers(ctx)
}

// GetServerHealth возвращает состояние сервера
func (s *ServerService) GetServerHealth(ctx context.Context, serverID int64) (*domain.ServerHealth, error) {
	return s.serverManager.GetServerHealth(ctx, serverID)
}

// GetAllServersHealth возвращает состояние всех серверов
func (s *ServerService) GetAllServersHealth(ctx context.Context) (map[int64]*domain.ServerHealth, error) {
	return s.serverManager.GetAllServersHealth(ctx)
}

// GetServerStats возвращает статистику сервера
func (s *ServerService) GetServerStats(ctx context.Context, serverID int64) (*domain.ServerStats, error) {
	return s.serverManager.GetServerStats(ctx, serverID)
}

// GetAllServersStats возвращает статистику всех серверов
func (s *ServerService) GetAllServersStats(ctx context.Context) (map[int64]*domain.ServerStats, error) {
	return s.serverManager.GetAllServersStats(ctx)
}

// StartHealthMonitoring запускает мониторинг состояния серверов
func (s *ServerService) StartHealthMonitoring(ctx context.Context) error {
	return s.serverManager.StartHealthMonitoring(ctx)
}

// StopHealthMonitoring останавливает мониторинг состояния серверов
func (s *ServerService) StopHealthMonitoring() error {
	return s.serverManager.StopHealthMonitoring()
}

// IsHealthMonitoringActive проверяет, активен ли мониторинг
func (s *ServerService) IsHealthMonitoringActive() bool {
	return s.serverManager.IsHealthMonitoringActive()
}
