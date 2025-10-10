# 🚀 MVP Roadmap - План реализации для полного функционирования бота

## 📊 Текущий статус

### ✅ Что уже есть:
- Clean Architecture с разделением на слои
- PostgreSQL репозитории (user, subscription, payment, vpn, referral, notification)
- Marzban API клиент (базовая реализация)
- Domain модели и Use Cases
- База данных с полной схемой
- Базовые Telegram handlers

### ❌ Что критично не хватает:
1. **VPN Service** - связь между локальной БД и Marzban API
2. **Payment Service** - обработка платежей
3. **Бизнес-логика** - автоматизация процессов
4. **Данные** - планы подписки в БД

---

## 🎯 MVP - Минимальный функционал (Этап 1)

### Задача 1: Добавить планы подписки в БД
**Время: 5 минут**

```sql
-- Добавить в migrations/002_seed_plans.sql
INSERT INTO plans (id, name, description, price, days, is_active) VALUES
('trial_7', 'Пробный период', '7 дней бесплатного VPN', 0.00, 7, true),
('monthly', 'Месячная подписка', '30 дней безлимитного VPN', 299.00, 30, true),
('quarterly', 'Квартальная подписка', '90 дней VPN со скидкой 15%', 749.00, 90, true),
('yearly', 'Годовая подписка', '365 дней VPN со скидкой 30%', 2499.00, 365, true);
```

**Применить:**
```bash
psql -d 3xui_bot -f migrations/002_seed_plans.sql
```

---

### Задача 2: Создать VPN Service
**Время: 2-3 часа**

**Файл:** `internal/service/vpn_service.go`

```go
type VPNService struct {
    localRepo   usecase.VPNRepository
    marzbanRepo marzban.MarzbanRepository
    subRepo     usecase.SubscriptionRepository
}

// Основные методы:
// - CreateVPNForSubscription(ctx, userID, planID) - создать VPN при активации подписки
// - GetUserVPNWithStats(ctx, userID) - получить VPN с данными из Marzban
// - SyncVPNStatus(ctx) - синхронизировать статусы
// - DeactivateExpiredVPNs(ctx) - деактивировать истекшие
```

**Ключевая логика:**
1. Создание VPN в Marzban при активации подписки
2. Сохранение связи в локальной БД
3. Получение статистики из Marzban по запросу
4. Автоматическое управление жизненным циклом

---

### Задача 3: Создать Payment Service
**Время: 4-6 часов**

**Файл:** `internal/service/payment_service.go`

```go
type PaymentService struct {
    paymentRepo usecase.PaymentRepository
    subRepo     usecase.SubscriptionRepository
    vpnService  *VPNService
    provider    PaymentProvider // ЮKassa/Stripe/etc
}

// Методы:
// - CreatePayment(ctx, userID, planID) - создать платеж
// - ProcessWebhook(ctx, data) - обработать webhook от платежной системы
// - CompletePayment(ctx, paymentID) - завершить платеж и создать подписку
// - RefundPayment(ctx, paymentID) - вернуть деньги
```

**Интеграция с платежным провайдером:**
- ЮKassa (рекомендуется для РФ)
- Stripe (для международных платежей)
- CryptoCloud (для крипты)

---

### Задача 4: Обработчики для бизнес-логики
**Время: 2-3 часа**

#### 4.1 Payment Handler
**Файл:** `internal/controller/bot/handlers/payment.go`

```go
// Обработка выбора плана и создания платежа
func (h *PaymentHandler) HandleSelectPlan(ctx, planID)
// Обработка успешной оплаты
func (h *PaymentHandler) HandlePaymentSuccess(ctx, paymentID)
// Обработка отмены платежа
func (h *PaymentHandler) HandlePaymentCancel(ctx, paymentID)
```

#### 4.2 VPN Handler
**Файл:** `internal/controller/bot/handlers/vpn.go`

```go
// Показать VPN конфигурации пользователя
func (h *VPNHandler) HandleShowVPNConfigs(ctx, userID)
// Получить файл конфигурации
func (h *VPNHandler) HandleDownloadConfig(ctx, vpnID)
// Показать статистику использования
func (h *VPNHandler) HandleShowStats(ctx, vpnID)
```

#### 4.3 Webhook Handler
**Файл:** `internal/controller/bot/handlers/webhook.go`

```go
// HTTP endpoint для приема webhook от платежной системы
func (h *WebhookHandler) HandlePaymentWebhook(w, r)
```

---

### Задача 5: Связать все вместе
**Время: 1-2 часа**

**Flow: Покупка подписки → Создание VPN**

```
1. Пользователь выбирает план → HandleSelectPlan()
2. Создается платеж → paymentService.CreatePayment()
3. Пользователь оплачивает → внешняя система
4. Webhook → HandlePaymentWebhook()
5. Обновление статуса → paymentService.CompletePayment()
6. Создание подписки → subscriptionUC.CreateSubscription()
7. Создание VPN → vpnService.CreateVPNForSubscription()
8. Уведомление пользователю → "VPN активирован!"
```

**Обновить:** `cmd/bot/main.go`
```go
// Добавить сервисы
vpnService := service.NewVPNService(vpnRepo, marzbanRepo, subRepo)
paymentService := service.NewPaymentService(paymentRepo, subRepo, vpnService, provider)

// Добавить handlers
paymentHandler := handlers.NewPaymentHandler(api, paymentService)
vpnHandler := handlers.NewVPNHandler(api, vpnService)
```

---

## 🔄 Этап 2: Автоматизация и фоновые задачи

### Задача 6: Scheduled Jobs
**Время: 2-3 часа**

**Файл:** `internal/scheduler/jobs.go`

```go
// Проверка истекших подписок (каждый час)
func CheckExpiredSubscriptions(ctx)

// Отправка уведомлений за 3 дня до окончания
func SendExpirationNotifications(ctx)

// Деактивация просроченных VPN (каждые 6 часов)
func DeactivateExpiredVPNs(ctx)
```

**Использовать:**
- `github.com/robfig/cron/v3` для планирования
- Запускать в отдельной горутине

---

### Задача 7: Улучшенная обработка ошибок
**Время: 1-2 часа**

**Файл:** `internal/pkg/logger/logger.go`

```go
// Структурированное логирование
type Logger struct {
    *zap.Logger
}

// Логирование с контекстом
func (l *Logger) ErrorWithContext(ctx, msg, err)
func (l *Logger) InfoWithContext(ctx, msg, fields)
```

---

## 🎨 Этап 3: UX улучшения

### Задача 8: Улучшить UI бота
**Время: 2-3 часа**

- Красивые карточки планов с emoji
- Inline кнопки для быстрых действий
- Прогресс-бары использования трафика
- QR коды для быстрого подключения

### Задача 9: Реферальная система
**Время: 3-4 часа**

- Обработка реферальных ссылок
- Бонусы за приглашенных друзей
- Статистика по рефералам

---

## 📈 Этап 4: Масштабирование

### Задача 10: Мониторинг
- Prometheus метрики
- Grafana дашборды
- Alertmanager для уведомлений

### Задача 11: Admin панель
- Web интерфейс для администрирования
- Управление планами
- Статистика

---

## 🔧 Порядок реализации для MVP

### Неделя 1: Ядро функционала
**Дни 1-2:**
- ✅ Добавить планы в БД
- ✅ Создать VPN Service
- ✅ Базовая интеграция с Marzban

**Дни 3-4:**
- ✅ Создать Payment Service
- ✅ Интеграция с ЮKassa
- ✅ Webhook handler

**День 5:**
- ✅ Связать все компоненты
- ✅ Тестирование flow

### Неделя 2: Автоматизация
**Дни 6-7:**
- ✅ Scheduled jobs
- ✅ Обработка истечений

**Дни 8-9:**
- ✅ Улучшение UI
- ✅ Error handling

**День 10:**
- ✅ Финальное тестирование
- ✅ Deploy на production

---

## 📝 Чек-лист готовности к запуску

### Must Have (критично):
- [ ] Планы добавлены в БД
- [ ] VPN Service работает
- [ ] Payment Service работает
- [ ] Webhook обрабатывается
- [ ] Flow покупки → VPN работает
- [ ] Обработка ошибок настроена

### Should Have (важно):
- [ ] Scheduled jobs запущены
- [ ] Логирование настроено
- [ ] UI улучшен
- [ ] Тесты написаны

### Nice to Have (желательно):
- [ ] Реферальная система
- [ ] Admin панель
- [ ] Мониторинг
- [ ] Документация

---

## 🚀 Быстрый старт для разработчика

```bash
# 1. Добавить планы
psql -d 3xui_bot -f migrations/002_seed_plans.sql

# 2. Создать VPN Service
touch internal/service/vpn_service.go

# 3. Создать Payment Service  
touch internal/service/payment_service.go

# 4. Создать handlers
touch internal/controller/bot/handlers/payment.go
touch internal/controller/bot/handlers/vpn.go
touch internal/controller/bot/handlers/webhook.go

# 5. Обновить main.go

# 6. Тестировать
make dev
```

---

## 📚 Полезные ссылки

- [ЮKassa API](https://yookassa.ru/developers/api)
- [Marzban API](https://github.com/Gozargah/Marzban)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

