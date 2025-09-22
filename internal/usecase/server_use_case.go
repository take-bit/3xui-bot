package usecase

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/domain"
)

// ServerUseCase представляет use case для работы с серверами
type ServerUseCase struct {
	serverService       domain.ServerService
	notificationService domain.NotificationService
}

// NewServerUseCase создает новый Server use case
func NewServerUseCase(
	serverService domain.ServerService,
	notificationService domain.NotificationService,
) *ServerUseCase {
	return &ServerUseCase{
		serverService:       serverService,
		notificationService: notificationService,
	}
}

// GetAvailableServers возвращает список доступных серверов
func (uc *ServerUseCase) GetAvailableServers(ctx context.Context) ([]*domain.Server, error) {
	return uc.serverService.GetAvailableServers(ctx)
}

// GetServerHealth возвращает состояние сервера
func (uc *ServerUseCase) GetServerHealth(ctx context.Context, serverID int64) (*domain.ServerHealth, error) {
	return uc.serverService.GetServerHealth(ctx, serverID)
}

// GetAllServersHealth возвращает состояние всех серверов
func (uc *ServerUseCase) GetAllServersHealth(ctx context.Context) (map[int64]*domain.ServerHealth, error) {
	return uc.serverService.GetAllServersHealth(ctx)
}

// GetServerStats возвращает статистику сервера
func (uc *ServerUseCase) GetServerStats(ctx context.Context, serverID int64) (*domain.ServerStats, error) {
	return uc.serverService.GetServerStats(ctx, serverID)
}

// GetAllServersStats возвращает статистику всех серверов
func (uc *ServerUseCase) GetAllServersStats(ctx context.Context) (map[int64]*domain.ServerStats, error) {
	return uc.serverService.GetAllServersStats(ctx)
}

// StartHealthMonitoring запускает мониторинг состояния серверов
func (uc *ServerUseCase) StartHealthMonitoring(ctx context.Context) error {
	err := uc.serverService.StartHealthMonitoring(ctx)
	if err != nil {
		return fmt.Errorf("failed to start health monitoring: %w", err)
	}

	// Отправляем уведомление о запуске мониторинга
	if uc.notificationService != nil {
		message := "🔍 Мониторинг серверов запущен\n\n📊 Отслеживание состояния серверов в реальном времени"
		_ = uc.notificationService.SendToAll(ctx, "Мониторинг серверов", message, false)
	}

	return nil
}

// StopHealthMonitoring останавливает мониторинг состояния серверов
func (uc *ServerUseCase) StopHealthMonitoring() error {
	err := uc.serverService.StopHealthMonitoring()
	if err != nil {
		return fmt.Errorf("failed to stop health monitoring: %w", err)
	}

	return nil
}

// IsHealthMonitoringActive проверяет, активен ли мониторинг
func (uc *ServerUseCase) IsHealthMonitoringActive() bool {
	return uc.serverService.IsHealthMonitoringActive()
}

// GetServerOverview возвращает обзор состояния серверов
func (uc *ServerUseCase) GetServerOverview(ctx context.Context) (*ServerOverview, error) {
	// 1. Получаем состояние всех серверов
	healthMap, err := uc.serverService.GetAllServersHealth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers health: %w", err)
	}

	// 2. Получаем статистику всех серверов
	statsMap, err := uc.serverService.GetAllServersStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers stats: %w", err)
	}

	// 3. Получаем список доступных серверов
	servers, err := uc.serverService.GetAvailableServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available servers: %w", err)
	}

	// 4. Формируем обзор
	overview := &ServerOverview{
		TotalServers:     len(servers),
		HealthyServers:   0,
		UnhealthyServers: 0,
		TotalClients:     0,
		ActiveClients:    0,
		AverageLoad:      0.0,
		LastUpdate:       time.Now(),
		Servers:          make([]ServerInfo, 0, len(servers)),
	}

	var totalLoad float64
	healthyCount := 0

	for _, server := range servers {
		health, healthExists := healthMap[server.ID]
		stats, statsExists := statsMap[server.ID]

		serverInfo := ServerInfo{
			Server: server,
			Health: health,
			Stats:  stats,
		}

		if healthExists && health.IsHealthy {
			healthyCount++
			overview.HealthyServers++
		} else {
			overview.UnhealthyServers++
		}

		if statsExists {
			overview.TotalClients += stats.TotalClients
			overview.ActiveClients += stats.ActiveClients
			totalLoad += stats.LoadPercentage
		}

		overview.Servers = append(overview.Servers, serverInfo)
	}

	// Вычисляем среднюю нагрузку
	if len(servers) > 0 {
		overview.AverageLoad = totalLoad / float64(len(servers))
	}

	return overview, nil
}

// CheckServerStatus проверяет статус сервера и отправляет уведомления при изменениях
func (uc *ServerUseCase) CheckServerStatus(ctx context.Context, serverID int64) error {
	// 1. Получаем текущее состояние сервера
	health, err := uc.serverService.GetServerHealth(ctx, serverID)
	if err != nil {
		return fmt.Errorf("failed to get server health: %w", err)
	}

	// 2. Отправляем уведомление о статусе сервера
	if uc.notificationService != nil {
		var message string
		if health.IsHealthy {
			message = fmt.Sprintf("✅ Сервер %d работает нормально\n\n📊 Нагрузка: %.1f%%\n👥 Активных клиентов: %d",
				serverID, health.CurrentLoad*100, health.ActiveClients)
		} else {
			message = fmt.Sprintf("❌ Сервер %d недоступен\n\n📋 Ошибка: %s\n⏰ Время проверки: %s",
				serverID, health.LastError, health.LastCheck.Format("02.01.2006 15:04"))
		}

		_ = uc.notificationService.SendToAll(ctx, "Статус сервера", message, false)
	}

	return nil
}

// GetServerLoadBalancingInfo возвращает информацию о балансировке нагрузки
func (uc *ServerUseCase) GetServerLoadBalancingInfo(ctx context.Context) (*LoadBalancingInfo, error) {
	// 1. Получаем обзор серверов
	overview, err := uc.GetServerOverview(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server overview: %w", err)
	}

	// 2. Анализируем нагрузку
	info := &LoadBalancingInfo{
		TotalServers:       overview.TotalServers,
		HealthyServers:     overview.HealthyServers,
		AverageLoad:        overview.AverageLoad,
		OverloadedServers:  0,
		UnderloadedServers: 0,
		Recommendations:    make([]string, 0),
	}

	// 3. Анализируем каждый сервер
	for _, serverInfo := range overview.Servers {
		if serverInfo.Stats != nil {
			load := serverInfo.Stats.LoadPercentage
			if load > 80 {
				info.OverloadedServers++
				info.Recommendations = append(info.Recommendations,
					fmt.Sprintf("Сервер %d перегружен (%.1f%%)", serverInfo.Server.ID, load))
			} else if load < 20 {
				info.UnderloadedServers++
				info.Recommendations = append(info.Recommendations,
					fmt.Sprintf("Сервер %d недогружен (%.1f%%)", serverInfo.Server.ID, load))
			}
		}
	}

	// 4. Добавляем общие рекомендации
	if info.OverloadedServers > 0 {
		info.Recommendations = append(info.Recommendations, "Рассмотрите возможность добавления новых серверов")
	}

	if info.UnderloadedServers > info.HealthyServers/2 {
		info.Recommendations = append(info.Recommendations, "Рассмотрите возможность отключения недогруженных серверов")
	}

	return info, nil
}

// SendServerMaintenanceNotification отправляет уведомление о технических работах на сервере
func (uc *ServerUseCase) SendServerMaintenanceNotification(ctx context.Context, serverID int64, startTime, endTime time.Time) error {
	if uc.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}

	message := fmt.Sprintf("🔧 Технические работы на сервере %d\n\n⏰ Начало: %s\n⏰ Окончание: %s\n\n⚠️ В это время VPN подключения на данном сервере могут быть недоступны",
		serverID, startTime.Format("02.01.2006 15:04"), endTime.Format("02.01.2006 15:04"))

	return uc.notificationService.SendToAll(ctx, "Технические работы", message, false)
}

// GetServerPerformanceMetrics возвращает метрики производительности серверов
func (uc *ServerUseCase) GetServerPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error) {
	// 1. Получаем статистику всех серверов
	statsMap, err := uc.serverService.GetAllServersStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers stats: %w", err)
	}

	// 2. Вычисляем метрики
	metrics := &PerformanceMetrics{
		TotalServers:  len(statsMap),
		TotalClients:  0,
		ActiveClients: 0,
		TotalTraffic:  0,
		AverageUptime: 0,
		AverageLoad:   0,
		LastUpdate:    time.Now(),
	}

	var totalUptime, totalLoad float64
	serverCount := 0

	for _, stats := range statsMap {
		metrics.TotalClients += stats.TotalClients
		metrics.ActiveClients += stats.ActiveClients
		metrics.TotalTraffic += stats.BytesTransferred
		totalUptime += float64(stats.Uptime.Seconds())
		totalLoad += stats.LoadPercentage
		serverCount++
	}

	if serverCount > 0 {
		metrics.AverageUptime = time.Duration(totalUptime/float64(serverCount)) * time.Second
		metrics.AverageLoad = totalLoad / float64(serverCount)
	}

	return metrics, nil
}

// ServerOverview представляет обзор состояния серверов
type ServerOverview struct {
	TotalServers     int          `json:"total_servers"`
	HealthyServers   int          `json:"healthy_servers"`
	UnhealthyServers int          `json:"unhealthy_servers"`
	TotalClients     int          `json:"total_clients"`
	ActiveClients    int          `json:"active_clients"`
	AverageLoad      float64      `json:"average_load"`
	LastUpdate       time.Time    `json:"last_update"`
	Servers          []ServerInfo `json:"servers"`
}

// ServerInfo представляет информацию о сервере
type ServerInfo struct {
	Server *domain.Server       `json:"server"`
	Health *domain.ServerHealth `json:"health,omitempty"`
	Stats  *domain.ServerStats  `json:"stats,omitempty"`
}

// LoadBalancingInfo представляет информацию о балансировке нагрузки
type LoadBalancingInfo struct {
	TotalServers       int      `json:"total_servers"`
	HealthyServers     int      `json:"healthy_servers"`
	AverageLoad        float64  `json:"average_load"`
	OverloadedServers  int      `json:"overloaded_servers"`
	UnderloadedServers int      `json:"underloaded_servers"`
	Recommendations    []string `json:"recommendations"`
}

// PerformanceMetrics представляет метрики производительности
type PerformanceMetrics struct {
	TotalServers  int           `json:"total_servers"`
	TotalClients  int           `json:"total_clients"`
	ActiveClients int           `json:"active_clients"`
	TotalTraffic  int64         `json:"total_traffic"`
	AverageUptime time.Duration `json:"average_uptime"`
	AverageLoad   float64       `json:"average_load"`
	LastUpdate    time.Time     `json:"last_update"`
}
