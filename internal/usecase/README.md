# 🎯 Use Case Layer

Use Case слой координирует бизнес-логику между различными сервисами, обеспечивая правильную архитектуру и разделение ответственности.

## 📋 Принципы

### **Use Case отвечает за:**
- ✅ **Координацию** работы между сервисами
- ✅ **Оркестрацию** бизнес-процессов
- ✅ **Валидацию** на уровне приложения
- ✅ **Уведомления** пользователей
- ✅ **Логирование** операций

### **Use Case НЕ отвечает за:**
- ❌ **Прямую работу с БД** (это делают сервисы)
- ❌ **API вызовы** (это делают репозитории)
- ❌ **Бизнес-правила** (это делают сервисы)

## 🏗️ Архитектура

```
┌─────────────────┐
│   Telegram Bot  │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│ UseCaseManager  │ ← Central coordinator
└─────────┬───────┘
          │
          ▼
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   VPN UseCase   │    │ Payment UseCase  │    │  User UseCase    │
│                 │    │                  │    │                  │
│ • VPN Logic     │    │ • Payment Logic  │    │ • User Logic     │
│ • Validation    │    │ • Webhooks       │    │ • Registration   │
└─────────────────┘    └──────────────────┘    └──────────────────┘
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   Service       │    │   Service        │    │   Service        │
│   Layer         │    │   Layer          │    │   Layer          │
└─────────────────┘    └──────────────────┘    └──────────────────┘
```

## 🚀 Use Cases

### **1. VPNUseCase**
Управление VPN подключениями:
- Создание VPN подключений
- Получение информации о подключениях
- Удаление VPN подключений
- Продление подписок
- Мониторинг серверов

### **2. PaymentUseCase**
Обработка платежей:
- Создание платежей
- Обработка webhook'ов
- Завершение платежей
- Обработка возвратов
- Статистика платежей

### **3. UserUseCase**
Управление пользователями:
- Регистрация пользователей
- Получение профилей
- Блокировка/разблокировка
- Обновление языка
- Статистика пользователей

### **4. SubscriptionUseCase**
Управление подписками:
- Создание пробных подписок
- Продление подписок
- Получение активных подписок
- Истечение подписок
- История подписок

### **5. PromocodeUseCase**
Работа с промокодами:
- Создание промокодов
- Применение промокодов
- Валидация промокодов
- Деактивация промокодов
- Массовое создание

### **6. ReferralUseCase**
Реферальная программа:
- Создание реферальных связей
- Обработка вознаграждений
- Статистика рефералов
- Выплата вознаграждений
- Реферальные ссылки

### **7. NotificationUseCase**
Уведомления:
- Отправка уведомлений
- Создание черновиков
- Массовые рассылки
- Планирование уведомлений
- Статистика уведомлений

### **8. ServerUseCase**
Управление серверами:
- Мониторинг состояния
- Получение статистики
- Балансировка нагрузки
- Технические работы
- Метрики производительности

## 🔧 UseCaseManager

Центральный координатор всех Use Cases:

```go
// Создание менеджера
manager := usecase.NewUseCaseManager(
    vpnService,
    paymentService,
    userService,
    subscriptionService,
    promocodeService,
    referralService,
    notificationService,
    serverService,
)

// Инициализация
err := manager.Initialize(ctx)
if err != nil {
    log.Fatal(err)
}

// Использование
user, err := manager.ProcessUserRegistration(ctx, telegramID, username, firstName, lastName, languageCode)
```

## 📊 Основные операции

### **Регистрация пользователя:**
```go
user, err := manager.ProcessUserRegistration(ctx, telegramID, username, firstName, lastName, languageCode)
if err != nil {
    return fmt.Errorf("failed to register user: %w", err)
}
```

### **Обработка платежа:**
```go
err := manager.ProcessPaymentWebhook(ctx, externalID, domain.PaymentStatusCompleted)
if err != nil {
    return fmt.Errorf("failed to process payment: %w", err)
}
```

### **Применение промокода:**
```go
result, err := manager.ProcessPromocodeApplication(ctx, userID, "WELCOME10")
if err != nil {
    return fmt.Errorf("failed to apply promocode: %w", err)
}
```

### **Получение дашборда пользователя:**
```go
dashboard, err := manager.GetUserDashboard(ctx, userID)
if err != nil {
    return fmt.Errorf("failed to get user dashboard: %w", err)
}
```

## 🛠️ Расширение функциональности

### **Добавление нового Use Case:**

```go
// 1. Создаем новый Use Case
type NewFeatureUseCase struct {
    newFeatureService domain.NewFeatureService
    notificationService domain.NotificationService
}

func NewNewFeatureUseCase(
    newFeatureService domain.NewFeatureService,
    notificationService domain.NotificationService,
) *NewFeatureUseCase {
    return &NewFeatureUseCase{
        newFeatureService: newFeatureService,
        notificationService: notificationService,
    }
}

// 2. Добавляем в UseCaseManager
type UseCaseManager struct {
    // ... existing use cases
    NewFeatureUseCase *NewFeatureUseCase
}

// 3. Инициализируем в конструкторе
func NewUseCaseManager(...) *UseCaseManager {
    // ... existing initialization
    newFeatureUseCase := NewNewFeatureUseCase(newFeatureService, notificationService)
    
    return &UseCaseManager{
        // ... existing use cases
        NewFeatureUseCase: newFeatureUseCase,
    }
}
```

## 📚 Преимущества Use Case слоя

### **1. Разделение ответственности**
- **Use Case** - координация
- **Service** - бизнес-логика
- **Repository** - работа с данными

### **2. Тестируемость**
```go
// Легко мокать зависимости
mockVPNService := &MockVPNService{}
mockServerService := &MockServerService{}

vpnUseCase := usecase.NewVPNUseCase(
    mockVPNService,
    mockServerService,
    mockUserService,
    mockSubscriptionService,
    mockNotificationService,
)
```

### **3. Переиспользование**
- Use Case можно использовать в разных местах
- Telegram Bot, HTTP API, CLI утилиты

### **4. Чистая архитектура**
- Соблюдение принципов SOLID
- Инверсия зависимостей
- Единственная ответственность

## 🔄 Жизненный цикл

### **Инициализация:**
```go
manager := usecase.NewUseCaseManager(...)
err := manager.Initialize(ctx)
```

### **Использование:**
```go
// Обработка команд бота
user, err := manager.ProcessUserRegistration(ctx, ...)
payment, err := manager.ProcessPaymentWebhook(ctx, ...)
```

### **Завершение:**
```go
err := manager.Shutdown()
```

## 📈 Мониторинг и логирование

### **Логирование операций:**
```go
func (uc *VPNUseCase) CreateVPNConnection(ctx context.Context, userID int64, region string) (*domain.VPNConnection, error) {
    // Логируем начало операции
    log.Printf("Creating VPN connection for user %d in region %s", userID, region)
    
    // ... логика создания
    
    // Логируем результат
    log.Printf("VPN connection created successfully for user %d", userID)
    return connection, nil
}
```

### **Метрики производительности:**
```go
func (uc *PaymentUseCase) ProcessPaymentWebhook(ctx context.Context, externalID string, status domain.PaymentStatus) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        metrics.RecordPaymentProcessingTime(duration)
    }()
    
    // ... логика обработки
}
```

## 🚀 Следующие шаги

1. **Интеграция с Telegram Bot** - Подключить Use Case к боту
2. **HTTP API** - Создать REST API для внешних интеграций
3. **Мониторинг** - Добавить метрики и алерты
4. **Тестирование** - Создать unit и integration тесты
5. **Документация** - API документация и примеры использования

## 📚 Дополнительные ресурсы

- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Use Case Pattern](https://martinfowler.com/bliki/UseCase.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)