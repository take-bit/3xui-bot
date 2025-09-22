# 🏗️ Архитектура множественных серверов 3X-UI

## 📋 Обзор

Система управления множественными серверами 3X-UI позволяет боту работать с несколькими серверами одновременно, обеспечивая:

- **Балансировку нагрузки** между серверами
- **Автоматический failover** при недоступности сервера
- **Географическую маршрутизацию** пользователей
- **Мониторинг состояния** всех серверов
- **Гибкие стратегии выбора** сервера

## 🏛️ Архитектура

### Основные компоненты

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   VPN Service   │────│  Server Manager  │────│  XUI Clients    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   User Repo     │    │  Server Selector │    │  Health Monitor │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 1. **Server Manager** (`internal/service/server/server_manager.go`)

Центральный компонент для управления множественными серверами:

- **Управление клиентами**: создание, обновление, удаление клиентов на любом сервере
- **Выбор сервера**: интеллектуальный выбор лучшего сервера для пользователя
- **Мониторинг**: отслеживание состояния всех серверов
- **Статистика**: сбор и предоставление статистики по серверам

### 2. **Server Selectors** (`internal/service/server/selectors.go`)

Реализуют различные стратегии выбора сервера:

- **Least Load**: сервер с наименьшей нагрузкой
- **Round Robin**: циклический выбор серверов
- **Random**: случайный выбор
- **Geographic**: выбор по географическому положению
- **Priority**: выбор по приоритету сервера

### 3. **Configuration** (`internal/domain/config.go`)

Расширенная конфигурация для поддержки множественных серверов:

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
    "auto_failover": true,
    "geographic_routing": true
  }
}
```

## 🔧 Использование

### Создание VPN подключения

```go
// Создание подключения с автоматическим выбором сервера
connection, err := vpnService.CreateConnection(ctx, userID)

// Создание подключения с указанием региона
connection, err := vpnService.CreateConnectionWithRegion(ctx, userID, "EU")
```

### Мониторинг серверов

```go
// Запуск мониторинга
err := vpnService.StartHealthMonitoring(ctx)

// Получение состояния всех серверов
health, err := vpnService.GetAllServersHealth(ctx)

// Получение статистики серверов
stats, err := vpnService.GetAllServersStats(ctx)
```

### Выбор сервера

```go
// Критерии выбора сервера
criteria := domain.ServerSelectionCriteria{
    UserID:      userID,
    Region:      "EU",
    MaxLoad:     0.8,
    ExcludeIDs:  []int64{4}, // исключить сервер 4
}

// Выбор сервера
server, err := serverManager.SelectServer(ctx, criteria)
```

## 📊 Стратегии выбора сервера

### 1. **Least Load** (по умолчанию)
Выбирает сервер с наименьшей нагрузкой:
```json
{
  "selection_strategy": "least_load"
}
```

### 2. **Round Robin**
Циклический выбор серверов для равномерного распределения:
```json
{
  "selection_strategy": "round_robin"
}
```

### 3. **Geographic**
Выбор сервера по географическому положению:
```json
{
  "selection_strategy": "geographic",
  "geographic_routing": true
}
```

### 4. **Priority**
Выбор по приоритету сервера:
```json
{
  "selection_strategy": "priority"
}
```

## 🏥 Мониторинг состояния

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

При недоступности сервера система автоматически:

1. **Обнаруживает** недоступность сервера
2. **Переключает** новых пользователей на доступные серверы
3. **Уведомляет** администратора о проблеме
4. **Восстанавливает** работу при восстановлении сервера

## 🌍 Географическая маршрутизация

Пользователи автоматически направляются на ближайшие серверы:

- **EU** → Европейские серверы
- **US** → Американские серверы  
- **ASIA** → Азиатские серверы

## ⚙️ Конфигурация

### Переменные окружения

```bash
# Основные настройки
BOT_TOKEN=your_bot_token
DATABASE_URL=postgres://user:pass@localhost:5432/db

# Серверы 3X-UI (секретные данные)
XUI_SERVER_1_HOST=eu1.example.com
XUI_SERVER_1_PASSWORD=password1
XUI_SERVER_2_HOST=us1.example.com
XUI_SERVER_2_PASSWORD=password2
```

### JSON конфигурация

```json
{
  "xui_servers": [
    {
      "id": 1,
      "name": "EU Server",
      "host": "eu1.example.com",
      "port": 2053,
      "username": "admin",
      "enabled": true,
      "priority": 1,
      "max_clients": 100,
      "region": "EU"
    }
  ],
  "server_management": {
    "selection_strategy": "least_load",
    "health_check_interval": "30s",
    "auto_failover": true,
    "load_balance_threshold": 0.8
  }
}
```

## 🚀 Преимущества

### 1. **Масштабируемость**
- Легко добавлять новые серверы
- Равномерное распределение нагрузки
- Автоматическое управление ресурсами

### 2. **Надежность**
- Автоматический failover
- Мониторинг состояния серверов
- Резервирование серверов

### 3. **Производительность**
- Выбор оптимального сервера
- Географическая маршрутизация
- Балансировка нагрузки

### 4. **Гибкость**
- Различные стратегии выбора
- Настраиваемые критерии
- Динамическое управление

## 🔧 Развертывание

### 1. Подготовка серверов

```bash
# На каждом сервере установить 3X-UI
bash <(curl -Ls https://raw.githubusercontent.com/MHSanaei/3x-ui/master/install.sh)
```

### 2. Настройка конфигурации

```bash
# Создать config.json
cp config.example.json config.json

# Заполнить секретные данные через переменные окружения
export XUI_SERVER_1_PASSWORD="password1"
export XUI_SERVER_2_PASSWORD="password2"
```

### 3. Запуск бота

```bash
# Запуск с мониторингом
go run cmd/bot/main.go -config config.json

# Проверка состояния серверов
curl http://localhost:2096/health
```

## 📈 Мониторинг и метрики

### Ключевые метрики

- **Доступность серверов**: процент времени работы
- **Нагрузка серверов**: количество активных клиентов
- **Время отклика**: задержка ответа серверов
- **Ошибки**: количество неудачных запросов

### Логирование

```json
{
  "level": "info",
  "msg": "Server health check completed",
  "server_id": 1,
  "is_healthy": true,
  "response_time": "150ms",
  "load_percentage": 45.5
}
```

## 🛠️ Разработка

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
```

### Расширение мониторинга

```go
// Добавление кастомных метрик
type CustomHealthChecker struct{}

func (c *CustomHealthChecker) CheckHealth(ctx context.Context, server *domain.Server) (*domain.ServerHealth, error) {
    // Ваша логика проверки здоровья
    return &domain.ServerHealth{
        ServerID:  server.ID,
        IsHealthy: true,
        // ... другие поля
    }, nil
}
```

## 🔒 Безопасность

### Рекомендации

1. **Используйте HTTPS** для всех соединений
2. **Ограничьте доступ** к панелям 3X-UI
3. **Регулярно обновляйте** пароли и токены
4. **Мониторьте** подозрительную активность
5. **Используйте VPN** для доступа к серверам

### Переменные окружения

```bash
# Никогда не храните секреты в коде!
XUI_SERVER_1_PASSWORD=secure_password_here
XUI_SERVER_2_PASSWORD=another_secure_password
DATABASE_URL=postgres://user:password@localhost:5432/db
```

## 📚 Дополнительные ресурсы

- [3X-UI Documentation](https://github.com/MHSanaei/3x-ui)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go Context Package](https://pkg.go.dev/context)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
