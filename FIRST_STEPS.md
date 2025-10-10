# 🚀 Первые шаги к запуску бота

## Шаг 1: Добавить планы подписки (⏱️ 2 минуты)

```bash
# Применить миграцию с планами
psql -d 3xui_bot -f migrations/002_seed_plans.sql
```

**Что это добавит:**
- 🎁 Пробный период (7 дней, бесплатно)
- 📅 Месячная подписка (30 дней, 299₽)
- 📆 Квартальная (90 дней, 749₽)
- 🎯 Годовая (365 дней, 2499₽)

**Проверить:**
```bash
psql -d 3xui_bot -c "SELECT id, name, price, days FROM plans;"
```

---

## Шаг 2: Создать VPN Service (следующий приоритет)

### Что нужно реализовать:

**Файл:** `internal/service/vpn_service.go`

<details>
<summary>📄 Код структуры (кликните для раскрытия)</summary>

```go
package service

import (
    "context"
    "fmt"
    
    "3xui-bot/internal/domain"
    "3xui-bot/internal/repository/marzban"
    "3xui-bot/internal/usecase"
)

type VPNService struct {
    localRepo   usecase.VPNRepository
    marzbanRepo marzban.MarzbanRepository
    subRepo     usecase.SubscriptionRepository
}

func NewVPNService(
    localRepo usecase.VPNRepository,
    marzbanRepo marzban.MarzbanRepository,
    subRepo usecase.SubscriptionRepository,
) *VPNService {
    return &VPNService{
        localRepo:   localRepo,
        marzbanRepo: marzbanRepo,
        subRepo:     subRepo,
    }
}

// CreateVPNForSubscription создает VPN при активации подписки
func (s *VPNService) CreateVPNForSubscription(ctx context.Context, userID int64, planID string) (*domain.VPNConnection, error) {
    // 1. Получить план
    // 2. Создать пользователя в Marzban
    // 3. Сохранить в локальной БД
    // 4. Вернуть VPN connection
    return nil, nil
}

// GetUserVPNWithStats получает VPN с данными из Marzban
func (s *VPNService) GetUserVPNWithStats(ctx context.Context, userID int64) ([]*domain.VPNConnection, error) {
    // 1. Получить локальные VPN
    // 2. Обогатить данными из Marzban
    // 3. Вернуть полные данные
    return nil, nil
}

// DeactivateExpiredVPNs деактивирует истекшие VPN
func (s *VPNService) DeactivateExpiredVPNs(ctx context.Context) error {
    // 1. Найти истекшие подписки
    // 2. Деактивировать в Marzban
    // 3. Обновить локальную БД
    return nil
}
```
</details>

---

## Шаг 3: Создать Payment Service

**Файл:** `internal/service/payment_service.go`

<details>
<summary>📄 Код структуры (кликните для раскрытия)</summary>

```go
package service

import (
    "context"
    
    "3xui-bot/internal/domain"
    "3xui-bot/internal/usecase"
)

type PaymentProvider interface {
    CreatePayment(amount float64, description string) (paymentURL string, err error)
    CheckPaymentStatus(paymentID string) (status string, err error)
}

type PaymentService struct {
    paymentRepo usecase.PaymentRepository
    subRepo     usecase.SubscriptionRepository
    vpnService  *VPNService
    provider    PaymentProvider
}

func NewPaymentService(
    paymentRepo usecase.PaymentRepository,
    subRepo usecase.SubscriptionRepository,
    vpnService *VPNService,
    provider PaymentProvider,
) *PaymentService {
    return &PaymentService{
        paymentRepo: paymentRepo,
        subRepo:     subRepo,
        vpnService:  vpnService,
        provider:    provider,
    }
}

// CreatePayment создает платеж для выбранного плана
func (s *PaymentService) CreatePayment(ctx context.Context, userID int64, planID string) (*domain.Payment, string, error) {
    // 1. Создать payment в БД
    // 2. Создать платеж в провайдере
    // 3. Вернуть payment и URL для оплаты
    return nil, "", nil
}

// ProcessWebhook обрабатывает webhook от платежной системы
func (s *PaymentService) ProcessWebhook(ctx context.Context, paymentID string) error {
    // 1. Обновить статус платежа
    // 2. Если успешно - создать подписку
    // 3. Создать VPN через VPNService
    return nil
}
```
</details>

---

## Шаг 4: Создать Handlers

### Payment Handler
**Файл:** `internal/controller/bot/handlers/payment.go`

```go
package handlers

// HandleSelectPlan обрабатывает выбор плана
func (h *PaymentHandler) HandleSelectPlan(ctx context.Context, userID int64, planID string) error {
    // Создать платеж и отправить ссылку пользователю
}

// HandlePaymentSuccess обрабатывает успешную оплату
func (h *PaymentHandler) HandlePaymentSuccess(ctx context.Context, paymentID string) error {
    // Активировать подписку и VPN
}
```

### VPN Handler
**Файл:** `internal/controller/bot/handlers/vpn.go`

```go
package handlers

// HandleShowVPNConfigs показывает VPN конфигурации
func (h *VPNHandler) HandleShowVPNConfigs(ctx context.Context, userID int64) error {
    // Получить и показать список VPN
}

// HandleDownloadConfig отправляет файл конфигурации
func (h *VPNHandler) HandleDownloadConfig(ctx context.Context, vpnID string) error {
    // Отправить .ovpn файл
}
```

---

## Шаг 5: Интеграция с платежной системой

### ЮKassa (рекомендуется для РФ)

```bash
# Установить SDK
go get github.com/rvinnie/yookassa-sdk-go/yookassa
```

**Файл:** `internal/provider/yookassa/provider.go`

```go
package yookassa

import (
    "github.com/rvinnie/yookassa-sdk-go/yookassa"
    "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type YooKassaProvider struct {
    client *yookassa.Client
}

func NewYooKassaProvider(shopID, secretKey string) *YooKassaProvider {
    client := yookassa.NewClient(shopID, secretKey)
    return &YooKassaProvider{client: client}
}

func (p *YooKassaProvider) CreatePayment(amount float64, description string) (string, error) {
    payment, err := p.client.Payment().Create(&payment.PaymentRequest{
        Amount: &payment.Amount{
            Value:    fmt.Sprintf("%.2f", amount),
            Currency: "RUB",
        },
        Description: description,
        Confirmation: &payment.Confirmation{
            Type:      "redirect",
            ReturnURL: "https://t.me/your_bot",
        },
    })
    
    if err != nil {
        return "", err
    }
    
    return payment.Confirmation.ConfirmationURL, nil
}
```

### Настройка в .env

```bash
# ЮKassa
YOOKASSA_SHOP_ID=your_shop_id
YOOKASSA_SECRET_KEY=your_secret_key
```

---

## Шаг 6: Обновить main.go

```go
// Создаем провайдера платежей
yookassaProvider := yookassa.NewYooKassaProvider(
    os.Getenv("YOOKASSA_SHOP_ID"),
    os.Getenv("YOOKASSA_SECRET_KEY"),
)

// Создаем VPN Service
vpnService := service.NewVPNService(vpnRepo, marzbanRepo, subRepo)

// Создаем Payment Service
paymentService := service.NewPaymentService(
    paymentRepo,
    subRepo,
    vpnService,
    yookassaProvider,
)

// Создаем handlers
paymentHandler := handlers.NewPaymentHandler(api, paymentService)
vpnHandler := handlers.NewVPNHandler(api, vpnService)
```

---

## 🧪 Тестирование

### 1. Проверить планы
```bash
psql -d 3xui_bot -c "SELECT * FROM plans;"
```

### 2. Запустить бота
```bash
make dev
```

### 3. Тестовый flow
1. `/start` - начать работу с ботом
2. Выбрать план
3. Создать тестовый платеж
4. Обработать webhook (локально через ngrok)
5. Проверить создание подписки и VPN

---

## 📋 Чек-лист первого запуска

- [ ] Планы добавлены в БД
- [ ] VPN Service создан
- [ ] Payment Service создан
- [ ] Handlers созданы
- [ ] Платежный провайдер настроен
- [ ] main.go обновлен
- [ ] .env заполнен
- [ ] Бот запускается без ошибок
- [ ] Тестовая покупка работает

---

## 🆘 Помощь

**Если что-то не работает:**

1. Проверить логи: `make dev`
2. Проверить БД: `psql -d 3xui_bot`
3. Проверить .env переменные
4. Посмотреть `MVP_ROADMAP.md` для деталей

**Полезные команды:**
```bash
make build          # Собрать проект
make test           # Запустить тесты
make migrate-up     # Применить миграции
make db-create      # Создать БД
```

