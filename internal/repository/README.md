# Repository Layer

Слой репозиториев содержит реализации для работы с базой данных и внешними сервисами.

## Структура

- **postgres/** - репозитории для работы с PostgreSQL
- **client/** - клиенты для внешних API (3X-UI, платежи)
- **repository.go** - интерфейсы и утилиты для работы с БД

## PostgreSQL Repositories

Пакет `postgres` содержит реализации репозиториев для работы с PostgreSQL:

- `user_repository.go` - работа с пользователями
- `subscription_repository.go` - работа с подписками
- `plan_repository.go` - работа с планами
- `server_repository.go` - работа с серверами
- `payment_repository.go` - работа с платежами
- `promocode_repository.go` - работа с промокодами
- `referral_repository.go` - работа с рефералами
- `notification_repository.go` - работа с уведомлениями
- `migrations.go` - SQL миграции для создания схемы БД

## External API Clients

Пакет `client` содержит клиенты для внешних API:

- `xui_client.go` - клиент для работы с 3X-UI API

## DBGetter Interface

`DBGetter` - интерфейс для получения соединений с базой данных:

```go
type DBGetter interface {
    GetDB(ctx context.Context) pgx.DB
}
```

### Реализации:

- `PoolDBGetter` - для работы с пулом соединений
- `TransactionDBGetter` - для работы в рамках транзакции

## Примеры использования

### Создание репозиториев

```go
import (
    "3xui-bot/internal/repository/postgres"
    "3xui-bot/internal/repository/client"
)

// Создание репозитория пользователей
userRepo := postgres.NewUserRepository(dbGetter)

// Создание XUI клиента
xuiConfig := client.XUIConfig{
    BaseURL:  "https://your-3xui-panel.com",
    Username: "admin",
    Password: "password",
    Timeout:  30 * time.Second,
}
xuiClient := client.NewXUIClient(xuiConfig)
```

### Работа с транзакциями

```go
// Создание транзакции
tx, err := pool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Создание DBGetter для транзакции
txDBGetter := &TransactionDBGetter{tx: tx}

// Использование репозиториев в транзакции
userRepo := postgres.NewUserRepository(txDBGetter)
subscriptionRepo := postgres.NewSubscriptionRepository(txDBGetter)

// Выполнение операций
// ...

// Подтверждение транзакции
return tx.Commit(ctx)
```

## Принципы

1. **Интерфейсы** - все репозитории реализуют интерфейсы из `domain` пакета
2. **DBGetter** - использование интерфейса для абстракции соединений с БД
3. **Обработка ошибок** - использование `errors.Is()` и `errors.As()`
4. **Доменные ошибки** - возврат специфических ошибок из `domain` пакета
5. **Контекст** - все методы принимают `context.Context`
6. **Транзакции** - поддержка работы в рамках транзакций