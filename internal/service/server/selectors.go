package server

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"3xui-bot/internal/domain"
)

// LeastLoadSelector выбирает сервер с наименьшей нагрузкой
type LeastLoadSelector struct{}

// NewLeastLoadSelector создает новый селектор по наименьшей нагрузке
func NewLeastLoadSelector() *LeastLoadSelector {
	return &LeastLoadSelector{}
}

// SelectServer выбирает сервер с наименьшей нагрузкой
func (s *LeastLoadSelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// Сортируем серверы по нагрузке
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].GetLoadPercentage() < servers[j].GetLoadPercentage()
	})

	return servers[0], nil
}

// GetStrategy возвращает стратегию селектора
func (s *LeastLoadSelector) GetStrategy() domain.ServerSelectionStrategy {
	return domain.StrategyLeastLoad
}

// RoundRobinSelector выбирает серверы циклически
type RoundRobinSelector struct {
	index int
	mutex sync.Mutex
}

// NewRoundRobinSelector создает новый round-robin селектор
func NewRoundRobinSelector() *RoundRobinSelector {
	return &RoundRobinSelector{}
}

// SelectServer выбирает сервер циклически
func (s *RoundRobinSelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	server := servers[s.index%len(servers)]
	s.index++

	return server, nil
}

// GetStrategy возвращает стратегию селектора
func (s *RoundRobinSelector) GetStrategy() domain.ServerSelectionStrategy {
	return domain.StrategyRoundRobin
}

// RandomSelector выбирает сервер случайно
type RandomSelector struct{}

// NewRandomSelector создает новый случайный селектор
func NewRandomSelector() *RandomSelector {
	rand.Seed(time.Now().UnixNano())
	return &RandomSelector{}
}

// SelectServer выбирает сервер случайно
func (s *RandomSelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	index := rand.Intn(len(servers))
	return servers[index], nil
}

// GetStrategy возвращает стратегию селектора
func (s *RandomSelector) GetStrategy() domain.ServerSelectionStrategy {
	return domain.StrategyRandom
}

// GeographicSelector выбирает сервер по географическому положению
type GeographicSelector struct{}

// NewGeographicSelector создает новый географический селектор
func NewGeographicSelector() *GeographicSelector {
	return &GeographicSelector{}
}

// SelectServer выбирает сервер по географическому положению
func (s *GeographicSelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// Если указан регион, ищем серверы в этом регионе
	if criteria.Region != "" {
		var regionalServers []*domain.Server
		for _, server := range servers {
			// TODO: Добавить поле Region в domain.Server
			// Пока используем все серверы, так как поле Region еще не добавлено
			_ = server // избегаем ошибки неиспользуемой переменной
			// if server.Region == criteria.Region {
			//     regionalServers = append(regionalServers, server)
			// }
		}

		if len(regionalServers) > 0 {
			// Выбираем сервер с наименьшей нагрузкой в регионе
			sort.Slice(regionalServers, func(i, j int) bool {
				return regionalServers[i].GetLoadPercentage() < regionalServers[j].GetLoadPercentage()
			})
			return regionalServers[0], nil
		}
	}

	// Fallback на наименьшую нагрузку
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].GetLoadPercentage() < servers[j].GetLoadPercentage()
	})

	return servers[0], nil
}

// GetStrategy возвращает стратегию селектора
func (s *GeographicSelector) GetStrategy() domain.ServerSelectionStrategy {
	return domain.StrategyGeographic
}

// PrioritySelector выбирает сервер по приоритету
type PrioritySelector struct{}

// NewPrioritySelector создает новый селектор по приоритету
func NewPrioritySelector() *PrioritySelector {
	return &PrioritySelector{}
}

// SelectServer выбирает сервер по приоритету
func (s *PrioritySelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	// TODO: Добавить поле Priority в domain.Server
	// Сортируем серверы по приоритету (высший приоритет = меньшее число)
	// sort.Slice(servers, func(i, j int) bool {
	//     return servers[i].Priority < servers[j].Priority
	// })

	// Пока возвращаем сервер с наименьшей нагрузкой
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].GetLoadPercentage() < servers[j].GetLoadPercentage()
	})

	return servers[0], nil
}

// GetStrategy возвращает стратегию селектора
func (s *PrioritySelector) GetStrategy() domain.ServerSelectionStrategy {
	return domain.StrategyPriority
}
