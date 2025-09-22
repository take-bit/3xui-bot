# 🏗️ Server Management Package

Пакет для управления множественными серверами 3X-UI с поддержкой балансировки нагрузки, мониторинга состояния и автоматического failover.

## 📋 Компоненты

### 1. **ServerManager** (`server_manager.go`)

Центральный компонент для управления множественными серверами:

```go
type ServerManager struct {
    config         *domain.Config
    serverRepo     domain.ServerRepository
    xuiClients     map[int64]domain.XUIClient
    healthStatus   map[int64]*domain.ServerHealth
    selectors      map[domain.ServerSelectionStrategy]domain.ServerSelector
    // ...
}
```

**Основные функции:**
- Управление клиентами на множественных серверах
- Выбор оптимального сервера для пользователя
- Мониторинг состояния всех серверов
- Сбор статистики по серверам

### 2. **Server Selectors** (`selectors.go`)

Реализуют различные стратегии выбора сервера:

#### **LeastLoadSelector**
Выбирает сервер с наименьшей нагрузкой:
```go
selector := NewLeastLoadSelector()
server, err := selector.SelectServer(ctx, servers, criteria)
```

#### **RoundRobinSelector**
Циклический выбор серверов:
```go
selector := NewRoundRobinSelector()
server, err := selector.SelectServer(ctx, servers, criteria)
```

#### **RandomSelector**
Случайный выбор сервера:
```go
selector := NewRandomSelector()
server, err := selector.SelectServer(ctx, servers, criteria)
```

#### **GeographicSelector**
Выбор по географическому положению:
```go
selector := NewGeographicSelector()
criteria := domain.ServerSelectionCriteria{Region: "EU"}
server, err := selector.SelectServer(ctx, servers, criteria)
```

#### **PrioritySelector**
Выбор по приоритету сервера:
```go
selector := NewPrioritySelector()
server, err := selector.SelectServer(ctx, servers, criteria)
```

## 🚀 Использование

### Создание ServerManager

```go
import (
    "3xui-bot/internal/service/server"
    "3xui-bot/internal/domain"
)

// Создание менеджера серверов
serverManager := server.NewServerManager(config, serverRepo)
```

### Выбор сервера

```go
// Критерии выбора
criteria := domain.ServerSelectionCriteria{
    UserID:     12345,
    Region:     "EU",
    MaxLoad:    0.8,
    ExcludeIDs: []int64{4}, // исключить сервер 4
}

// Выбор сервера
server, err := serverManager.SelectServer(ctx, criteria)
if err != nil {
    return fmt.Errorf("failed to select server: %w", err)
}
```

### Создание клиента

```go
// Создание клиента на выбранном сервере
err := serverManager.CreateClient(ctx, server.ID, userID, uuid, totalGB, expiryTime)
if err != nil {
    return fmt.Errorf("failed to create client: %w", err)
}
```

### Мониторинг состояния

```go
// Запуск мониторинга
err := serverManager.StartHealthMonitoring(ctx)
if err != nil {
    return fmt.Errorf("failed to start monitoring: %w", err)
}

// Получение состояния сервера
health, err := serverManager.GetServerHealth(ctx, serverID)
if err != nil {
    return fmt.Errorf("failed to get server health: %w", err)
}

fmt.Printf("Server %d is healthy: %t\n", serverID, health.IsHealthy)
```

### Получение статистики

```go
// Статистика одного сервера
stats, err := serverManager.GetServerStats(ctx, serverID)
if err != nil {
    return fmt.Errorf("failed to get server stats: %w", err)
}

fmt.Printf("Server %d: %d/%d clients (%.1f%% load)\n", 
    stats.ServerID, stats.ActiveClients, stats.MaxClients, stats.LoadPercentage)

// Статистика всех серверов
allStats, err := serverManager.GetAllServersStats(ctx)
if err != nil {
    return fmt.Errorf("failed to get all servers stats: %w", err)
}

for serverID, stats := range allStats {
    fmt.Printf("Server %d: %d active clients\n", serverID, stats.ActiveClients)
}
```

## ⚙️ Конфигурация

### Настройка серверов

```json
{
  "xui_servers": [
    {
      "id": 1,
      "name": "EU Server 1",
      "host": "eu1.example.com",
      "port": 2053,
      "username": "admin",
      "password": "",
      "enabled": true,
      "priority": 1,
      "max_clients": 100,
      "region": "EU"
    }
  ],
  "server_management": {
    "selection_strategy": "least_load",
    "health_check_interval": "30s",
    "health_check_timeout": "10s",
    "max_retries": 3,
    "retry_delay": "5s",
    "load_balance_threshold": 0.8,
    "auto_failover": true,
    "geographic_routing": true
  }
}
```

### Стратегии выбора

- **`least_load`** - сервер с наименьшей нагрузкой (по умолчанию)
- **`round_robin`** - циклический выбор
- **`random`** - случайный выбор
- **`geographic`** - по географическому положению
- **`priority`** - по приоритету сервера

## 📊 Мониторинг

### Health Check

Система автоматически проверяет состояние серверов:

```go
type ServerHealth struct {
    ServerID      int64         `json:"server_id"`
    IsHealthy     bool          `json:"is_healthy"`
    LastCheck     time.Time     `json:"last_check"`
    ResponseTime  time.Duration `json:"response_time"`
    ErrorCount    int           `json:"error_count"`
    CurrentLoad   float64       `json:"current_load"`
    ActiveClients int           `json:"active_clients"`
}
```

### Статистика серверов

```go
type ServerStats struct {
    ServerID         int64         `json:"server_id"`
    TotalClients     int           `json:"total_clients"`
    ActiveClients    int           `json:"active_clients"`
    LoadPercentage   float64       `json:"load_percentage"`
    BytesTransferred int64         `json:"bytes_transferred"`
    ConnectionsCount int64         `json:"connections_count"`
}
```

## 🔄 Автоматический Failover

При недоступности сервера система:

1. **Обнаруживает** недоступность через health check
2. **Исключает** недоступный сервер из выбора
3. **Перенаправляет** новых пользователей на доступные серверы
4. **Восстанавливает** работу при восстановлении сервера

## 🌍 Географическая маршрутизация

Пользователи автоматически направляются на ближайшие серверы:

```go
// Выбор сервера в определенном регионе
criteria := domain.ServerSelectionCriteria{
    Region: "EU", // Европейские серверы
}

server, err := serverManager.SelectServer(ctx, criteria)
```

## 🛠️ Расширение функциональности

### Добавление нового селектора

```go
type CustomSelector struct{}

func (s *CustomSelector) SelectServer(ctx context.Context, servers []*domain.Server, criteria domain.ServerSelectionCriteria) (*domain.Server, error) {
    // Ваша логика выбора сервера
    return servers[0], nil
}

func (s *CustomSelector) GetStrategy() domain.ServerSelectionStrategy {
    return "custom"
}

// Регистрация селектора
serverManager.selectors["custom"] = &CustomSelector{}
```

### Кастомные критерии выбора

```go
criteria := domain.ServerSelectionCriteria{
    UserID:      12345,
    Region:      "EU",
    PreferredID: 1,           // предпочтительный сервер
    ExcludeIDs:  []int64{4},  // исключить серверы
    MinPriority: 2,           // минимальный приоритет
    MaxLoad:     0.8,         // максимальная нагрузка
}
```

## 🔒 Безопасность

### Рекомендации

1. **Используйте HTTPS** для всех соединений с 3X-UI
2. **Ограничьте доступ** к панелям 3X-UI по IP
3. **Регулярно обновляйте** пароли и токены
4. **Мониторьте** подозрительную активность
5. **Используйте VPN** для доступа к серверам

### Переменные окружения

```bash
# Секретные данные серверов
XUI_SERVER_1_PASSWORD=secure_password_1
XUI_SERVER_2_PASSWORD=secure_password_2
XUI_SERVER_3_PASSWORD=secure_password_3
```

## 📈 Производительность

### Оптимизации

1. **Кэширование** состояния серверов
2. **Параллельные** health checks
3. **Умная** балансировка нагрузки
4. **Минимальные** задержки при выборе сервера

### Метрики

- **Время отклика** серверов
- **Процент доступности** серверов
- **Распределение нагрузки** между серверами
- **Количество ошибок** подключения

## 🐛 Отладка

### Логирование

```go
// Включение детального логирования
log.SetLevel(log.DebugLevel)

// Логирование выбора сервера
log.Debugf("Selected server %d for user %d using strategy %s", 
    server.ID, userID, strategy)
```

### Проверка состояния

```go
// Проверка всех серверов
health, err := serverManager.GetAllServersHealth(ctx)
if err != nil {
    log.Errorf("Failed to get servers health: %v", err)
    return
}

for serverID, health := range health {
    log.Infof("Server %d: healthy=%t, load=%.1f%%, clients=%d", 
        serverID, health.IsHealthy, health.CurrentLoad*100, health.ActiveClients)
}
```

## 📚 Дополнительные ресурсы

- [3X-UI Documentation](https://github.com/MHSanaei/3x-ui)
- [Go Context Package](https://pkg.go.dev/context)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Load Balancing Strategies](https://en.wikipedia.org/wiki/Load_balancing_(computing))
