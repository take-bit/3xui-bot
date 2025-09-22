# Telegram Bot

Этот пакет содержит реализацию Telegram Bot для 3X-UI VPN сервиса, построенную на основе Clean Architecture.

## Структура

```
internal/telegram/
├── bot.go                    # Основной класс бота
├── handler.go                # Базовый обработчик
├── handler_interface.go      # Интерфейсы для обработчиков
├── handlers/                 # Обработчики команд
│   ├── start.go             # Команда /start
│   ├── help.go              # Команда /help
│   ├── profile.go           # Команда /profile
│   ├── vpn.go               # Команда /vpn
│   ├── subscription.go      # Команда /subscription
│   ├── payment.go           # Команда /payment
│   ├── promocode.go         # Команда /promocode
│   ├── referral.go          # Команда /referral
│   ├── settings.go          # Команда /settings
│   ├── text.go              # Обработка текстовых сообщений
│   ├── callback.go          # Обработка callback query
│   └── default.go           # Обработчик по умолчанию
└── middleware/              # Middleware
    ├── logging.go           # Логирование
    ├── auth.go              # Аутентификация
    └── rate_limit.go        # Ограничение частоты запросов
```

## Основные компоненты

### Bot

Основной класс бота, который:
- Управляет подключением к Telegram API
- Обрабатывает входящие обновления
- Применяет middleware
- Маршрутизирует команды к соответствующим обработчикам

### Handlers

Обработчики команд и сообщений:

- **StartHandler** - Обрабатывает команду `/start`, регистрирует пользователей
- **HelpHandler** - Показывает справку по боту
- **ProfileHandler** - Отображает профиль пользователя
- **VPNHandler** - Управление VPN подключениями
- **SubscriptionHandler** - Управление подписками
- **PaymentHandler** - Обработка платежей
- **PromocodeHandler** - Применение промокодов
- **ReferralHandler** - Реферальная программа
- **SettingsHandler** - Настройки бота
- **TextHandler** - Обработка текстовых сообщений
- **CallbackHandler** - Обработка callback query
- **DefaultHandler** - Обработчик неизвестных команд

### Middleware

- **LoggingMiddleware** - Логирует все входящие обновления
- **AuthMiddleware** - Проверяет аутентификацию пользователей
- **RateLimitMiddleware** - Ограничивает частоту запросов

## Использование

```go
// Создание бота
bot, err := telegram.NewBot(cfg, useCaseManager)
if err != nil {
    log.Fatal(err)
}

// Запуск бота
err = bot.Start(ctx)
if err != nil {
    log.Fatal(err)
}
```

## Конфигурация

Бот использует следующие переменные окружения:

- `BOT_TOKEN` - Токен Telegram бота (обязательно)
- `BOT_URL` - URL для webhook (опционально, если не указан, используется polling)
- `LOG_LEVEL` - Уровень логирования (по умолчанию: info)

## Архитектура

Бот построен на основе Clean Architecture:

1. **Handlers** - Слой представления, обрабатывает пользовательский ввод
2. **Use Cases** - Бизнес-логика, вызывается из handlers
3. **Services** - Доменные сервисы
4. **Repositories** - Доступ к данным

## Особенности

- **Graceful Shutdown** - Корректное завершение работы
- **Error Handling** - Обработка ошибок с уведомлениями пользователей
- **Rate Limiting** - Защита от спама
- **Logging** - Подробное логирование всех операций
- **Webhook/Polling** - Поддержка обоих режимов работы
- **Inline Keyboards** - Интерактивные клавиатуры
- **Callback Query** - Обработка нажатий на кнопки

## Расширение

Для добавления новых команд:

1. Создайте новый handler в пакете `handlers`
2. Реализуйте интерфейс `HandlerInterface`
3. Зарегистрируйте handler в `bot.go`
4. Добавьте обработку в `CallbackHandler` если необходимо

## Тестирование

Бот можно тестировать с помощью:
- Unit тестов для handlers
- Integration тестов с mock Telegram API
- End-to-end тестов с реальным ботом
