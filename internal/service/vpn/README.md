# VPN Service Package

Пакет `vpn` содержит сервис для управления VPN подключениями пользователей.

## VPNService

`VPNService` - основной сервис для работы с VPN подключениями.

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
import "3xui-bot/internal/service/vpn"

// Создание VPN сервиса
vpnService := vpn.NewVPNService(
    userRepo,
    subscriptionRepo,
    serverRepo,
    xuiClient,
)

// Создание подключения (использует Telegram User ID)
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

### Идентификация клиентов

VPN сервис использует **Telegram User ID** для идентификации клиентов в 3X-UI системе:

- **Email генерируется автоматически**: `user_{userID}@vpn.local`
- **Не требует ввода email** от пользователя Telegram
- **Уникальная идентификация** каждого пользователя
- **Совместимость с 3X-UI** - система требует email, но мы генерируем его из User ID

## VPNConnection

`VPNConnection` находится в `domain` слое (`internal/domain/vpn_connection.go`) и представляет VPN подключение пользователя:

```go
type VPNConnection struct {
    UserID       int64     `json:"user_id"`
    ServerID     int64     `json:"server_id"`
    XUIInboundID int       `json:"xui_inbound_id"`
    XUIClientID  string    `json:"xui_client_id"`
    UUID         string    `json:"uuid"`
    Email        string    `json:"email"`
    ConfigURL    string    `json:"config_url"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
}
```

### Методы VPNConnection:

- **IsActive()** - проверяет, активно ли подключение
- **IsExpired()** - проверяет, истекло ли подключение
- **GetRemainingDays()** - возвращает количество оставшихся дней
- **GetStatus()** - возвращает статус подключения ("active", "expired", "inactive")

## Принципы дизайна:

1. **Единственная ответственность** - сервис отвечает только за VPN логику
2. **Зависимость от интерфейсов** - зависит от интерфейсов, а не от конкретных реализаций
3. **Обработка ошибок** - все ошибки обрабатываются и оборачиваются с контекстом
4. **Контекст** - все методы принимают `context.Context` для отмены операций
5. **Чистая архитектура** - не знает о деталях реализации репозиториев и клиентов
6. **Разделение слоев** - доменные модели в `domain` слое, сервисы только для бизнес-логики
