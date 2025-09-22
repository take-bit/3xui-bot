package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"3xui-bot/internal/config"
	"3xui-bot/internal/domain"
	"3xui-bot/internal/repository/client"
)

// ServerManager реализует интерфейс domain.XUIClientManager
type ServerManager struct {
	config              *config.Config
	serverRepo          domain.ServerRepository
	xuiClients          map[int64]domain.XUIClient
	selectors           map[domain.ServerSelectionStrategy]domain.ServerSelector
	healthMonitorWG     sync.WaitGroup
	healthMonitorCtx    context.Context
	healthMonitorCancel context.CancelFunc
	isMonitoringActive  bool
	mu                  sync.RWMutex
}

// NewServerManager создает новый экземпляр ServerManager
func NewServerManager(cfg *config.Config, serverRepo domain.ServerRepository) (*ServerManager, error) {
	sm := &ServerManager{
		config:     cfg,
		serverRepo: serverRepo,
		xuiClients: make(map[int64]domain.XUIClient),
		selectors:  make(map[domain.ServerSelectionStrategy]domain.ServerSelector),
	}

	// Инициализируем XUI клиенты для каждого сервера из конфигурации
	for _, serverCfg := range cfg.XUIServers {
		xuiClientConfig := client.XUIConfig{
			BaseURL:  serverCfg.Host,
			Username: serverCfg.Username,
			Password: serverCfg.Password,
			Token:    serverCfg.Token,
			Timeout:  10 * time.Second, // TODO: Вынести в конфиг
		}
		sm.xuiClients[serverCfg.ID] = client.NewXUIClient(xuiClientConfig)
	}

	// Инициализируем селекторы
	sm.selectors[domain.StrategyLeastLoad] = NewLeastLoadSelector()
	sm.selectors[domain.StrategyRoundRobin] = NewRoundRobinSelector()
	sm.selectors[domain.StrategyRandom] = NewRandomSelector()
	sm.selectors[domain.StrategyGeographic] = NewGeographicSelector()
	sm.selectors[domain.StrategyPriority] = NewPrioritySelector()

	return sm, nil
}

// GetClient возвращает XUIClient для указанного serverID
func (sm *ServerManager) GetClient(serverID int64) (domain.XUIClient, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	client, ok := sm.xuiClients[serverID]
	if !ok {
		return nil, fmt.Errorf("XUI client not found for server ID: %d", serverID)
	}
	return client, nil
}

// SelectBestServer выбирает лучший сервер на основе заданной стратегии
func (sm *ServerManager) SelectBestServer(ctx context.Context, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Получаем все активные и здоровые серверы из репозитория
	servers, err := sm.serverRepo.GetAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available servers from repository: %w", err)
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no available servers to select from")
	}

	// Выбираем селектор на основе стратегии из конфигурации
	selector, ok := sm.selectors[domain.ServerSelectionStrategy(sm.config.ServerManagement.SelectionStrategy)]
	if !ok {
		// Fallback на стратегию наименьшей нагрузки, если указанная не найдена
		log.Printf("Warning: Unknown server selection strategy '%s'. Falling back to 'least_load'.", sm.config.ServerManagement.SelectionStrategy)
		selector = sm.selectors[domain.StrategyLeastLoad]
	}

	selectedServer, err := selector.SelectServer(ctx, servers, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to select server using strategy %s: %w", selector.GetStrategy(), err)
	}

	return selectedServer, nil
}

// CreateClient добавляет клиента на указанный сервер
func (sm *ServerManager) CreateClient(ctx context.Context, serverID int64, inboundID int, userID int64, uuid string, totalGB int64, expiryTime int64) error {
	xuiClient, err := sm.GetClient(serverID)
	if err != nil {
		return err
	}
	return xuiClient.AddClient(ctx, inboundID, userID, uuid, totalGB, expiryTime)
}

// DeleteClient удаляет клиента с указанного сервера
func (sm *ServerManager) DeleteClient(ctx context.Context, serverID int64, inboundID int, clientID string) error {
	xuiClient, err := sm.GetClient(serverID)
	if err != nil {
		return err
	}
	return xuiClient.DeleteClient(ctx, inboundID, clientID)
}

// UpdateClient обновляет клиента на указанном сервере
func (sm *ServerManager) UpdateClient(ctx context.Context, serverID int64, inboundID int, clientID string, userID int64, totalGB int64, expiryTime int64) error {
	xuiClient, err := sm.GetClient(serverID)
	if err != nil {
		return err
	}
	return xuiClient.UpdateClient(ctx, inboundID, clientID, userID, totalGB, expiryTime)
}

// GetAvailableServers возвращает список доступных серверов
func (sm *ServerManager) GetAvailableServers(ctx context.Context) ([]*domain.Server, error) {
	return sm.serverRepo.GetAvailable(ctx)
}

// CheckServerHealth проверяет состояние конкретного сервера
func (sm *ServerManager) CheckServerHealth(ctx context.Context, serverID int64) (*domain.ServerHealth, error) {
	sm.mu.RLock()
	serverCfg, found := findServerConfig(sm.config.XUIServers, serverID)
	sm.mu.RUnlock()

	if !found {
		return nil, fmt.Errorf("server configuration not found for ID: %d", serverID)
	}

	xuiClient, err := sm.GetClient(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get XUI client: %w", err)
	}

	health := &domain.ServerHealth{
		ServerID:      serverID,
		IsHealthy:     false,
		LastCheck:     time.Now(),
		ResponseTime:  0,
		ErrorCount:    0,
		LastError:     "",
		CurrentLoad:   0.0,
		ActiveClients: 0,
	}

	startTime := time.Now()
	err = xuiClient.Login(ctx)
	health.ResponseTime = time.Since(startTime)

	if err != nil {
		health.IsHealthy = false
		health.LastError = err.Error()
		health.ErrorCount++
		log.Printf("Server %s (ID: %d) is down: %v", serverCfg.Name, serverID, err)
		return health, fmt.Errorf("server %s (ID: %d) is down: %w", serverCfg.Name, serverID, err)
	}

	inbounds, err := xuiClient.GetInbounds(ctx)
	if err != nil {
		health.IsHealthy = false
		health.LastError = fmt.Sprintf("failed to get inbounds: %v", err)
		health.ErrorCount++
		log.Printf("Server %s (ID: %d) is partially down (inbounds error): %v", serverCfg.Name, serverID, err)
		return health, fmt.Errorf("server %s (ID: %d) is partially down (inbounds error): %w", serverCfg.Name, serverID, err)
	}

	health.IsHealthy = true
	health.ErrorCount = 0
	health.LastError = ""

	// Подсчитываем активных клиентов
	activeClients := 0
	for _, inbound := range inbounds {
		clients, err := xuiClient.GetClients(ctx, inbound.ID)
		if err != nil {
			continue // Пропускаем ошибки
		}
		for _, client := range clients {
			if client.Enable {
				activeClients++
			}
		}
	}
	health.ActiveClients = activeClients

	// Вычисляем нагрузку (простая логика)
	server, err := sm.serverRepo.GetByID(ctx, serverID)
	if err == nil && server.MaxClients > 0 {
		health.CurrentLoad = float64(activeClients) / float64(server.MaxClients)
	}

	return health, nil
}

// GetServerHealth возвращает состояние сервера
func (sm *ServerManager) GetServerHealth(ctx context.Context, serverID int64) (*domain.ServerHealth, error) {
	return sm.CheckServerHealth(ctx, serverID)
}

// GetAllServersHealth проверяет состояние всех серверов
func (sm *ServerManager) GetAllServersHealth(ctx context.Context) (map[int64]*domain.ServerHealth, error) {
	healthMap := make(map[int64]*domain.ServerHealth)
	var wg sync.WaitGroup
	var mu sync.Mutex // Для защиты healthMap

	for _, serverCfg := range sm.config.XUIServers {
		if !serverCfg.Enabled {
			continue
		}

		wg.Add(1)
		go func(sCfg config.XUIServerConfig) {
			defer wg.Done()
			health, err := sm.CheckServerHealth(ctx, sCfg.ID)
			mu.Lock()
			healthMap[sCfg.ID] = health
			mu.Unlock()
			if err != nil {
				log.Printf("Error checking health for server %s (ID: %d): %v", sCfg.Name, sCfg.ID, err)
			}
		}(serverCfg)
	}
	wg.Wait()
	return healthMap, nil
}

// GetServerStats возвращает статистику сервера
func (sm *ServerManager) GetServerStats(ctx context.Context, serverID int64) (*domain.ServerStats, error) {
	// Получаем информацию о сервере
	server, err := sm.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	// Получаем XUI клиент
	xuiClient, err := sm.GetClient(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get XUI client: %w", err)
	}

	// Аутентифицируемся
	if err := xuiClient.Login(ctx); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	// Получаем список inbound'ов
	inbounds, err := xuiClient.GetInbounds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbounds: %w", err)
	}

	// Подсчитываем клиентов
	totalClients := 0
	activeClients := 0
	var totalTraffic, usedTraffic int64

	for _, inbound := range inbounds {
		clients, err := xuiClient.GetClients(ctx, inbound.ID)
		if err != nil {
			continue // Пропускаем ошибки
		}

		totalClients += len(clients)
		for _, client := range clients {
			if client.Enable {
				activeClients++
			}
			totalTraffic += client.TotalGB * 1024 * 1024 * 1024 // Конвертируем GB в байты
			usedTraffic += client.UsedGB * 1024 * 1024 * 1024
		}
	}

	// Вычисляем процент нагрузки
	loadPercentage := 0.0
	if server.MaxClients > 0 {
		loadPercentage = float64(totalClients) / float64(server.MaxClients) * 100
	}

	stats := &domain.ServerStats{
		ServerID:         serverID,
		TotalClients:     totalClients,
		ActiveClients:    activeClients,
		MaxClients:       server.MaxClients,
		LoadPercentage:   loadPercentage,
		Uptime:           time.Since(server.CreatedAt),
		LastActivity:     time.Now(), // TODO: Получать реальное время последней активности
		BytesTransferred: usedTraffic,
		ConnectionsCount: int64(activeClients),
	}

	return stats, nil
}

// GetAllServersStats возвращает статистику всех серверов
func (sm *ServerManager) GetAllServersStats(ctx context.Context) (map[int64]*domain.ServerStats, error) {
	statsMap := make(map[int64]*domain.ServerStats)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, serverCfg := range sm.config.XUIServers {
		if !serverCfg.Enabled {
			continue
		}

		wg.Add(1)
		go func(sCfg config.XUIServerConfig) {
			defer wg.Done()
			stats, err := sm.GetServerStats(ctx, sCfg.ID)
			mu.Lock()
			statsMap[sCfg.ID] = stats
			mu.Unlock()
			if err != nil {
				log.Printf("Error getting stats for server %s (ID: %d): %v", sCfg.Name, sCfg.ID, err)
			}
		}(serverCfg)
	}
	wg.Wait()
	return statsMap, nil
}

// StartHealthMonitoring запускает периодический мониторинг состояния серверов
func (sm *ServerManager) StartHealthMonitoring(ctx context.Context) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.isMonitoringActive {
		return fmt.Errorf("health monitoring is already active")
	}

	sm.healthMonitorCtx, sm.healthMonitorCancel = context.WithCancel(ctx)
	sm.isMonitoringActive = true

	go func() {
		ticker := time.NewTicker(sm.config.ServerManagement.HealthCheckInterval)
		defer ticker.Stop()

		log.Println("Starting server health monitoring...")
		// Выполняем первую проверку сразу
		sm.GetAllServersHealth(sm.healthMonitorCtx)

		for {
			select {
			case <-ticker.C:
				sm.GetAllServersHealth(sm.healthMonitorCtx)
			case <-sm.healthMonitorCtx.Done():
				log.Println("Server health monitoring stopped.")
				return
			}
		}
	}()
	return nil
}

// StopHealthMonitoring останавливает мониторинг состояния серверов
func (sm *ServerManager) StopHealthMonitoring() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.isMonitoringActive && sm.healthMonitorCancel != nil {
		sm.healthMonitorCancel()
		sm.isMonitoringActive = false
	}
}

// IsHealthMonitoringActive проверяет, активен ли мониторинг здоровья
func (sm *ServerManager) IsHealthMonitoringActive() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.isMonitoringActive
}

func findServerConfig(serverConfigs []config.XUIServerConfig, serverID int64) (config.XUIServerConfig, bool) {
	for _, cfg := range serverConfigs {
		if cfg.ID == serverID {
			return cfg, true
		}
	}
	return config.XUIServerConfig{}, false
}
