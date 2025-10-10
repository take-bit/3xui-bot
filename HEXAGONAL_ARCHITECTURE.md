# 🏗️ Hexagonal Architecture (Ports & Adapters)

## Обзор

Проект реструктурирован согласно принципам **Hexagonal Architecture** (также известной как Ports & Adapters) и **Clean Architecture**.

## 📐 Структура слоев

```
┌─────────────────────────────────────────────────────┐
│                   ADAPTERS                           │
│  (Telegram, PostgreSQL, Marzban, HTTP)              │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│                    PORTS                             │
│  (Interfaces: Repo, Marzban, Notifier, Clock)       │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│                  USE CASES                           │
│  (Business Logic, Orchestration)                    │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│                    CORE                              │
│  (Domain Entities, Business Rules)                  │
└─────────────────────────────────────────────────────┘
```

## 📂 Структура папок

```
internal/
├── core/           # Доменные сущности (entities)
│   ├── user.go
│   ├── subscription.go
│   ├── payment.go
│   ├── vpn.go
│   ├── referral.go
│   └── notification.go
│
├── ports/          # Интерфейсы (contracts)
│   ├── repo.go        # Репозитории
│   ├── marzban.go     # Marzban VPN Manager
│   ├── notifier.go    # Отправка уведомлений
│   └── clock.go       # Время (для тестируемости)
│
├── usecase/        # Бизнес-логика
│   ├── user.go
│   ├── subscription.go
│   ├── payment.go     # Оркестратор
│   ├── vpn.go
│   ├── referral.go
│   ├── notification.go
│   └── dto.go
│
├── adapters/       # Реализации интерфейсов
│   ├── db/postgres/      # PostgreSQL репозитории
│   ├── marzban/          # Marzban API клиент
│   ├── notify/           # Telegram notifier
│   ├── payment/          # Payment providers
│   └── bot/telegram/     # Telegram bot
│
├── app/            # Application layer
│   ├── container.go  # DI container
│   └── run.go        # App runner
│
├── scheduler/      # Фоновые задачи
└── pkg/            # Утилиты
```

## 🎯 Принципы

### 1. Dependency Rule (Правило зависимостей)

**Зависимости направлены ВНУТРЬ:**

```
Adapters → Ports → UseCase → Core
```

- **Core** не знает ни о чем (только stdlib)
- **UseCase** знает о Core и Ports
- **Adapters** знают о Ports и реализуют их
- **Ports** определяют контракты

### 2. Use Case Orchestration

**UseCase может вызывать другой UseCase:**

```go
type PaymentUseCase struct {
    paymentRepo    ports.PaymentRepo
    subscriptionUC *SubscriptionUseCase  // ← UseCase композиция
    vpnUC          *VPNUseCase
    notifUC        *NotificationUseCase
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, paymentID string) error {
    // 1. Обновляем платеж
    uc.paymentRepo.UpdateStatus(...)
    
    // 2. Создаем подписку (через другой UC)
    subscription, _ := uc.subscriptionUC.CreateSubscription(...)
    
    // 3. Создаем VPN (через другой UC)
    vpn, _ := uc.vpnUC.CreateVPNForSubscription(...)
    
    // 4. Отправляем уведомление (через другой UC)
    uc.notifUC.SendNotification(...)
    
    return nil
}
```

### 3. Ports (Interfaces)

Все внешние зависимости через интерфейсы:

- `ports.UserRepo` - работа с БД пользователей
- `ports.Marzban` - Marzban API
- `ports.Notifier` - отправка сообщений
- `ports.Clock` - работа со временем

### 4. Adapters (Реализации)

Каждый адаптер реализует один или несколько портов:

- `adapters/db/postgres/*` → `ports.*Repo`
- `adapters/marzban/client.go` → `ports.Marzban`
- `adapters/notify/tg_notifier.go` → `ports.Notifier`

## 🔄 Flow обработки запроса

### Пример: Покупка подписки

```
1. Telegram Update
       ↓
2. Router.HandleUpdate()
       ↓
3. PaymentHandler.HandleSelectPlan()
       ↓
4. PaymentUseCase.CreatePaymentForPlan()
   ├→ SubscriptionUseCase.GetPlan()
   │   └→ ports.PlanRepo.GetPlanByID()
   │       └→ adapters/db/postgres/subscription
   │
   └→ ports.PaymentRepo.CreatePayment()
       └→ adapters/db/postgres/payment

5. User pays externally
       ↓
6. PaymentUseCase.ProcessPaymentSuccess()
   ├→ ports.PaymentRepo.UpdateStatus()
   ├→ SubscriptionUseCase.CreateSubscription()
   ├→ VPNUseCase.CreateVPNForSubscription()
   │   ├→ ports.VPNRepo.Create()
   │   └→ ports.Marzban.CreateUser()
   │
   └→ NotificationUseCase.SendNotification()
       └→ ports.Notifier.Send()
           └→ Telegram API
```

## 💡 Преимущества архитектуры

### 1. Тестируемость
- UseCase легко тестировать с моками ports
- Не нужна реальная БД/API для тестов

### 2. Гибкость
- Легко заменить PostgreSQL на MongoDB
- Легко добавить другой VPN провайдер
- Легко добавить HTTP API

### 3. Независимость
- Core не зависит от фреймворков
- UseCase не зависит от деталей реализации
- Бизнес-логика изолирована

### 4. Расширяемость
- Новые адаптеры без изменения core
- Новые use cases без изменения портов
- Новые функции без ломающих изменений

## 🔧 Dependency Injection

Все зависимости собираются в `app/container.go`:

```go
// 1. Конфиг → Logger → DB
config → logger → pgxpool

// 2. Репозитории (реализуют ports)
UserRepo, SubscriptionRepo, etc.

// 3. Внешние клиенты
Marzban, Notifier, Clock

// 4. Use Cases (базовые)
UserUC, SubscriptionUC, ReferralUC

// 5. Use Cases (зависимые)
VPNUC, NotificationUC

// 6. Use Cases (оркестраторы)
PaymentUC (зависит от SubUC, VPNUC, NotifUC)

// 7. Адаптеры
Router, Scheduler
```

## 📝 Создание нового адаптера

### Пример: Добавление Email уведомлений

1. **Обновить порт:**
```go
// internal/ports/notifier.go
type EmailNotifier interface {
    SendEmail(ctx context.Context, to, subject, body string) error
}
```

2. **Создать адаптер:**
```go
// internal/adapters/email/smtp.go
type SMTPNotifier struct {
    client *smtp.Client
}

func (s *SMTPNotifier) SendEmail(ctx, to, subject, body string) error {
    // Реализация
}
```

3. **Использовать в UseCase:**
```go
type NotificationUseCase struct {
    telegramNotifier ports.Notifier
    emailNotifier    ports.EmailNotifier  // ← Новая зависимость
}
```

## 🧪 Тестирование

### Пример: Тест PaymentUseCase

```go
func TestProcessPaymentSuccess(t *testing.T) {
    // Создаем моки для всех портов
    mockPaymentRepo := &MockPaymentRepo{}
    mockSubUC := &MockSubscriptionUseCase{}
    mockVPNUC := &MockVPNUseCase{}
    mockNotifUC := &MockNotificationUseCase{}
    mockProvider := &MockPaymentProvider{}
    
    // Создаем UseCase с моками
    uc := usecase.NewPaymentUseCase(
        mockPaymentRepo,
        mockSubUC,
        mockVPNUC,
        mockNotifUC,
        mockProvider,
    )
    
    // Тестируем
    err := uc.ProcessPaymentSuccess(ctx, "payment_id", "plan_id")
    assert.NoError(t, err)
}
```

## 📚 Дополнительные ресурсы

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

