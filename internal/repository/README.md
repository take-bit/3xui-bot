# Repository Layer

Этот пакет содержит реализацию репозиториев для работы с базой данных PostgreSQL.

## Структура

### Репозитории
- `user_repository.go` - работа с пользователями
- `subscription_repository.go` - работа с подписками
- `plan_repository.go` - работа с планами
- `server_repository.go` - работа с серверами
- `payment_repository.go` - работа с платежами
- `promocode_repository.go` - работа с промокодами
- `referral_repository.go` - работа с рефералами
- `notification_repository.go` - работа с уведомлениями

### Базовые компоненты
- `repository.go` - интерфейс DBGetter для работы с базой данных
- `migrations.go` - миграции базы данных

## Технологии

- **PostgreSQL** - основная база данных
- **pgx/v5** - драйвер PostgreSQL для Go
- **transactor** - управление транзакциями

## Использование

```go
// Создание DBGetter (например, через transactor)
dbGetter := transactor.NewTransactorFromPool(pool)

// Создание репозиториев
userRepo := NewUserRepository(dbGetter)
subscriptionRepo := NewSubscriptionRepository(dbGetter)
paymentRepo := NewPaymentRepository(dbGetter)

// Выполнение миграций
err = Migrate(ctx, pool)
if err != nil {
    log.Fatal(err)
}

// Использование репозитория
user, err := userRepo.GetByTelegramID(ctx, 123456789)
if err != nil {
    log.Printf("Error: %v", err)
}
```

## Транзакции

Все репозитории поддерживают транзакции через `transactor`:

```go
// Выполнение в транзакции
err := transactor.Transaction(ctx, func(ctx context.Context) error {
    // Создание пользователя
    user := &domain.User{...}
    if err := userRepo.Create(ctx, user); err != nil {
        return err
    }
    
    // Создание подписки
    subscription := &domain.Subscription{...}
    if err := subscriptionRepo.Create(ctx, subscription); err != nil {
        return err
    }
    
    return nil
})
```

## Миграции

Миграции выполняются автоматически при запуске приложения:

```go
// Выполнение миграций
err = Migrate(ctx, pool)
```

## Индексы

Все таблицы имеют оптимизированные индексы для быстрого поиска:
- Уникальные индексы для внешних ключей
- Составные индексы для сложных запросов
- Частичные индексы для активных записей
