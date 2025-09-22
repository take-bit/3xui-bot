# Client Package

Пакет `client` содержит клиенты для взаимодействия с внешними API и сервисами.

## XUI Client

`XUIClient` - клиент для работы с 3X-UI API панелью.

### Основные функции:

- **Аутентификация** - логин через username/password или токен
- **Управление клиентами** - добавление, обновление, удаление VPN клиентов
- **Получение данных** - список inbound'ов и клиентов
- **HTTP клиент** - с настраиваемым таймаутом и обработкой ошибок

### Конфигурация:

```go
import "3xui-bot/internal/repository/client"

config := client.XUIConfig{
    BaseURL:  "https://your-3xui-panel.com",
    Username: "admin",
    Password: "password",
    Token:    "", // Опционально, если есть токен
    Timeout:  30 * time.Second,
}

xuiClient := client.NewXUIClient(config)
```

### Использование:

```go
// Аутентификация
err := xuiClient.Login(ctx)
if err != nil {
    // Обработка ошибки
}

// Получение списка inbound'ов
inbounds, err := xuiClient.GetInbounds(ctx)
if err != nil {
    // Обработка ошибки
}

// Добавление клиента (использует User ID вместо email)
err = xuiClient.AddClient(ctx, inboundID, userID, uuid, totalGB, expiryTime)
if err != nil {
    // Обработка ошибки
}
```

### Идентификация клиентов:

XUI Client использует **Telegram User ID** для идентификации клиентов:

- **Email генерируется автоматически**: `user_{userID}@vpn.local`
- **Не требует ввода email** от пользователя
- **Уникальная идентификация** каждого пользователя
- **Совместимость с 3X-UI** - система требует email, но мы генерируем его из User ID

### Обработка ошибок:

- HTTP ошибки с указанием статус кода
- JSON парсинг с обработкой ошибок
- Проверка успешности ответа от API
- Оборачивание ошибок с контекстом

### Принципы дизайна:

1. **Реализация интерфейса** - клиент реализует `domain.XUIClient`
2. **Конфигурируемость** - настройка через структуру конфигурации
3. **Контекст** - все методы принимают `context.Context`
4. **Обработка ошибок** - все ошибки оборачиваются с контекстом
5. **HTTP клиент** - переиспользуемый HTTP клиент с таймаутом
