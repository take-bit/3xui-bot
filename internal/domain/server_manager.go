package domain

import (
	"context"
	"time"
)

// ServerSelectionStrategy определяет стратегию выбора сервера
type ServerSelectionStrategy string

const (
	StrategyLeastLoad  ServerSelectionStrategy = "least_load"  // наименьшая нагрузка
	StrategyRoundRobin ServerSelectionStrategy = "round_robin" // циклический выбор
	StrategyRandom     ServerSelectionStrategy = "random"      // случайный выбор
	StrategyGeographic ServerSelectionStrategy = "geographic"  // по географическому положению
	StrategyLatency    ServerSelectionStrategy = "latency"     // по задержке
	StrategyPriority   ServerSelectionStrategy = "priority"    // по приоритету
)

// ServerHealth представляет состояние сервера
type ServerHealth struct {
	ServerID      int64         `json:"server_id"`
	IsHealthy     bool          `json:"is_healthy"`
	LastCheck     time.Time     `json:"last_check"`
	ResponseTime  time.Duration `json:"response_time"`
	ErrorCount    int           `json:"error_count"`
	LastError     string        `json:"last_error,omitempty"`
	CurrentLoad   float64       `json:"current_load"`   // текущая нагрузка (0.0-1.0)
	ActiveClients int           `json:"active_clients"` // количество активных клиентов
}

// ServerStats представляет статистику сервера
type ServerStats struct {
	ServerID         int64         `json:"server_id"`
	TotalClients     int           `json:"total_clients"`
	ActiveClients    int           `json:"active_clients"`
	MaxClients       int           `json:"max_clients"`
	LoadPercentage   float64       `json:"load_percentage"`
	Uptime           time.Duration `json:"uptime"`
	LastActivity     time.Time     `json:"last_activity"`
	BytesTransferred int64         `json:"bytes_transferred"`
	ConnectionsCount int64         `json:"connections_count"`
}

// ServerSelectionCriteria критерии для выбора сервера
type ServerSelectionCriteria struct {
	UserID      int64   `json:"user_id,omitempty"`
	Region      string  `json:"region,omitempty"`
	PreferredID int64   `json:"preferred_id,omitempty"`
	ExcludeIDs  []int64 `json:"exclude_ids,omitempty"`
	MinPriority int     `json:"min_priority,omitempty"`
	MaxLoad     float64 `json:"max_load,omitempty"` // максимальная допустимая нагрузка
}

// XUIClientManager интерфейс для управления множественными клиентами 3X-UI
type XUIClientManager interface {
	// Управление серверами
	GetAvailableServers(ctx context.Context) ([]*Server, error)
	GetServerByID(ctx context.Context, id int64) (*Server, error)
	GetServerHealth(ctx context.Context, id int64) (*ServerHealth, error)
	GetAllServersHealth(ctx context.Context) (map[int64]*ServerHealth, error)

	// Выбор сервера
	SelectServer(ctx context.Context, criteria ServerSelectionCriteria) (*Server, error)
	SelectBestServer(ctx context.Context, userID int64, region string) (*Server, error)

	// Статистика
	GetServerStats(ctx context.Context, id int64) (*ServerStats, error)
	GetAllServersStats(ctx context.Context) (map[int64]*ServerStats, error)

	// Мониторинг
	StartHealthMonitoring(ctx context.Context) error
	StopHealthMonitoring() error
	IsHealthMonitoringActive() bool

	// Управление клиентами
	CreateClient(ctx context.Context, serverID int64, userID int64, uuid string, totalGB int64, expiryTime int64) error
	UpdateClient(ctx context.Context, serverID int64, userID int64, uuid string, totalGB int64, expiryTime int64) error
	DeleteClient(ctx context.Context, serverID int64, userID int64) error
	GetClient(ctx context.Context, serverID int64, userID int64) (*XUIClientInfo, error)

	// Получение информации о сервере
	GetInbounds(ctx context.Context, serverID int64) ([]XUIServerInfo, error)
	GetClients(ctx context.Context, serverID int64, inboundID int) ([]XUIClientInfo, error)
}

// ServerSelector интерфейс для выбора сервера
type ServerSelector interface {
	SelectServer(ctx context.Context, servers []*Server, criteria ServerSelectionCriteria) (*Server, error)
	GetStrategy() ServerSelectionStrategy
}

// HealthChecker интерфейс для проверки состояния сервера
type HealthChecker interface {
	CheckHealth(ctx context.Context, server *Server) (*ServerHealth, error)
	CheckAllServers(ctx context.Context, servers []*Server) (map[int64]*ServerHealth, error)
}

// LoadBalancer интерфейс для балансировки нагрузки
type LoadBalancer interface {
	GetServerLoad(ctx context.Context, serverID int64) (float64, error)
	GetLeastLoadedServer(ctx context.Context, servers []*Server) (*Server, error)
	IsServerOverloaded(ctx context.Context, serverID int64, threshold float64) (bool, error)
}

// GeographicRouter интерфейс для географической маршрутизации
type GeographicRouter interface {
	GetUserRegion(userID int64) (string, error)
	GetNearestServers(region string, servers []*Server) []*Server
	GetServerLatency(ctx context.Context, serverID int64) (time.Duration, error)
}
