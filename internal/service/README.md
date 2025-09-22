# Service Layer

Слой сервисов содержит бизнес-логику приложения. Сервисы координируют работу между репозиториями и доменными моделями.

## Структура

- **vpn/** - VPN сервис для управления подключениями

## VPN Service

`VPNService` в пакете `vpn` - основной сервис для управления VPN подключениями пользователей.

### Основные функции:

- **CreateConnection** - создание VPN подключения для пользователя
- **DeleteConnection** - удаление VPN подключения
- **GetConnectionInfo** - получение информации о подключении
- **CheckServerHealth** - проверка состояния сервера
- **UpdateConnectionExpiry** - обновление времени истечения подключения

### Архитектура:

```
VPNService
├── UserRepository (проверка пользователя)
├── SubscriptionRepository (проверка активной подписки)
├── ServerRepository (выбор доступного сервера)
└── XUIClient (взаимодействие с 3X-UI API)
```

### Процесс создания подключения:

1. **Проверка пользователя** - убеждаемся, что пользователь существует
2. **Проверка подписки** - проверяем наличие активной подписки
3. **Выбор сервера** - выбираем сервер с наименьшей нагрузкой
4. **Аутентификация** - логинимся в 3X-UI панель
5. **Создание клиента** - добавляем клиента в 3X-UI
6. **Генерация конфига** - создаем URL для скачивания конфигурации

### Обработка ошибок:

- Используется `errors.Is()` для проверки специфических ошибок
- Ошибки оборачиваются с контекстом через `fmt.Errorf("...: %w", err)`
- Возвращаются доменные ошибки (`domain.ErrUserNotFound`, etc.)

### Пример использования:

```go
import (
    "3xui-bot/internal/repository/client"
    "3xui-bot/internal/service/vpn"
)

// Создание XUI клиента
config := client.XUIConfig{
    BaseURL:  "https://your-3xui-panel.com",
    Username: "admin",
    Password: "password",
    Timeout:  30 * time.Second,
}
xuiClient := client.NewXUIClient(config)

// Создание VPN сервиса
vpnService := vpn.NewVPNService(
    userRepo,
    subscriptionRepo,
    serverRepo,
    xuiClient,
)

// Создание подключения
connection, err := vpnService.CreateConnection(ctx, userID)
if err != nil {
    if errors.Is(err, domain.ErrUserNotFound) {
        // Пользователь не найден
    }
    // Обработка других ошибок
}

// Получение информации о подключении
info, err := vpnService.GetConnectionInfo(ctx, userID)
if err != nil {
    // Обработка ошибок
}
```

## XUI Client

`XUIClient` в `internal/repository/client/xui_client.go` - клиент для работы с 3X-UI API.

> **Примечание:** XUI Client находится в `repository/client` пакете, так как это инфраструктурная зависимость для взаимодействия с внешними API.

### Основные методы:

- **Login** - аутентификация в 3X-UI панели
- **GetInbounds** - получение списка inbound'ов
- **GetClients** - получение списка клиентов
- **AddClient** - добавление нового клиента
- **UpdateClient** - обновление данных клиента
- **DeleteClient** - удаление клиента

### Конфигурация:

```go
import "3xui-bot/internal/repository/client"

config := client.XUIConfig{
    BaseURL:  "https://your-3xui-panel.com",
    Username: "admin",
    Password: "password",
    Timeout:  30 * time.Second,
}

xuiClient := client.NewXUIClient(config)
```

### Обработка ошибок:

- HTTP ошибки обрабатываются с указанием статус кода
- JSON парсинг с обработкой ошибок
- Проверка успешности ответа от API

## Принципы дизайна:

1. **Единственная ответственность** - каждый сервис отвечает за свою область
2. **Зависимость от интерфейсов** - сервисы зависят от интерфейсов, а не от конкретных реализаций
3. **Обработка ошибок** - все ошибки обрабатываются и оборачиваются с контекстом
4. **Контекст** - все методы принимают `context.Context` для отмены операций
5. **Чистая архитектура** - сервисы не знают о деталях реализации репозиториев и клиентов
6. **Разделение слоев** - клиенты для внешних API в пакете `internal/repository/client`
7. **Модульность** - каждый сервис в отдельном пакете (например, `vpn/`)
