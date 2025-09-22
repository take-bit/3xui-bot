# 3X-UI Bot

Telegram бот для управления VPN сервисом на базе 3X-UI панели, построенный с использованием Clean Architecture.

## 🚀 Возможности

- **VPN Управление**: Создание, удаление и управление VPN подключениями
- **Подписки**: Гибкая система тарифных планов с различными сроками
- **Платежи**: Поддержка множественных платежных систем (Cryptomus, Heleket, YooKassa, YooMoney, Telegram Stars)
- **Реферальная программа**: Система приглашений с вознаграждениями
- **Промокоды**: Система скидок и бонусов
- **Пробный период**: Бесплатный доступ для новых пользователей
- **Многосерверность**: Автоматическое распределение нагрузки между серверами
- **Мониторинг**: Отслеживание состояния серверов и здоровья системы

## 🏗️ Архитектура

Проект построен на основе Clean Architecture с четким разделением слоев:

```
├── cmd/                    # Точки входа приложения
├── internal/
│   ├── config/            # Конфигурация
│   ├── domain/            # Доменные модели и интерфейсы
│   ├── interfaces/        # Интерфейсы для внешних слоев
│   ├── repository/        # Слой доступа к данным
│   ├── service/           # Бизнес-логика
│   ├── usecase/           # Сценарии использования
│   └── telegram/          # Telegram Bot
└── migrations/            # Миграции базы данных
```

### Слои архитектуры

1. **Domain Layer** (`internal/domain/`) - Доменные модели, интерфейсы и бизнес-правила
2. **Repository Layer** (`internal/repository/`) - Абстракция доступа к данным
3. **Service Layer** (`internal/service/`) - Бизнес-логика и доменные сервисы
4. **Use Case Layer** (`internal/usecase/`) - Сценарии использования приложения
5. **Presentation Layer** (`internal/telegram/`) - Telegram Bot интерфейс

## 🛠️ Технологии

- **Go 1.21+** - Основной язык программирования
- **PostgreSQL** - База данных
- **pgx** - Драйвер PostgreSQL
- **Telegram Bot API** - Интерфейс Telegram
- **3X-UI API** - Управление VPN серверами
- **Docker** - Контейнеризация

## 📦 Установка

### Требования

- Go 1.21 или выше
- PostgreSQL 13+
- 3X-UI панель
- Telegram Bot Token

### Быстрый старт

1. **Клонирование репозитория**
   ```bash
   git clone https://github.com/your-username/3xui-bot.git
   cd 3xui-bot
   ```

2. **Установка зависимостей**
   ```bash
   make setup
   ```

3. **Настройка конфигурации**
   ```bash
   cp config.example.json config.json
   cp plans.example.json plans.json
   # Отредактируйте config.json и plans.json
   ```

4. **Настройка переменных окружения**
   ```bash
   export BOT_TOKEN=your_bot_token_here
   export DATABASE_URL=postgres://user:password@localhost:5432/3xui_bot?sslmode=disable
   ```

5. **Запуск**
   ```bash
   make run
   ```

### Docker

```bash
# Сборка образа
make docker-build

# Запуск с docker-compose
make docker-compose-up
```

## ⚙️ Конфигурация

### Основные параметры

- `BOT_TOKEN` - Токен Telegram бота (обязательно)
- `BOT_URL` - URL для webhook (опционально)
- `DATABASE_URL` - Строка подключения к PostgreSQL
- `LOG_LEVEL` - Уровень логирования (debug, info, warn, error)

### Платежные системы

Поддерживаются следующие платежные системы:

- **Cryptomus** - Криптовалютные платежи
- **Heleket** - Криптовалютные платежи
- **YooKassa** - Банковские карты и электронные кошельки
- **YooMoney** - Электронные кошельки
- **Telegram Stars** - Платежи через Telegram

### Серверы 3X-UI

Настройте ваши 3X-UI серверы в `config.json`:

```json
{
  "xui_servers": [
    {
      "id": 1,
      "name": "Server 1",
      "host": "192.168.1.100",
      "port": 2053,
      "username": "admin",
      "password": "password",
      "enabled": true,
      "region": "US"
    }
  ]
}
```

## 🎯 Использование

### Команды бота

- `/start` - Начать работу с ботом
- `/help` - Показать справку
- `/profile` - Мой профиль
- `/subscription` - Управление подпиской
- `/vpn` - VPN подключение
- `/promocode` - Применить промокод
- `/referral` - Реферальная программа
- `/settings` - Настройки

### Планы подписки

- **1 месяц** - 100₽
- **3 месяца** - 250₽ (экономия 50₽)
- **6 месяцев** - 450₽ (экономия 150₽)
- **1 год** - 800₽ (экономия 400₽)

## 🔧 Разработка

### Структура проекта

```
├── cmd/bot/               # Главный файл приложения
├── internal/
│   ├── config/           # Конфигурация
│   ├── domain/           # Доменные модели
│   ├── repository/       # Репозитории
│   ├── service/          # Сервисы
│   ├── usecase/          # Use Cases
│   └── telegram/         # Telegram Bot
├── migrations/           # Миграции БД
├── config.example.json   # Пример конфигурации
├── plans.example.json    # Пример планов
└── Makefile             # Команды сборки
```

### Команды разработки

```bash
# Установка инструментов
make install-tools

# Форматирование кода
make format

# Линтинг
make lint

# Тесты
make test

# Тесты с покрытием
make test-coverage

# Сборка
make build

# Запуск в режиме разработки
make run-dev
```

### Добавление новых функций

1. **Доменная модель** - Добавьте в `internal/domain/`
2. **Репозиторий** - Реализуйте в `internal/repository/`
3. **Сервис** - Создайте в `internal/service/`
4. **Use Case** - Добавьте в `internal/usecase/`
5. **Handler** - Создайте в `internal/telegram/handlers/`

## 📊 Мониторинг

### Логирование

Бот поддерживает структурированное логирование в JSON формате:

```json
{
  "level": "info",
  "time": "2024-01-01T12:00:00Z",
  "message": "User registered",
  "user_id": 123456789,
  "telegram_id": 987654321
}
```

### Метрики

- Количество активных пользователей
- Статистика платежей
- Нагрузка на серверы
- Производительность системы

## 🔒 Безопасность

- **Rate Limiting** - Ограничение частоты запросов
- **JWT токены** - Аутентификация API
- **Шифрование** - Защита чувствительных данных
- **Валидация** - Проверка входных данных
- **Логирование** - Аудит всех операций

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🆘 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [Issues](https://github.com/your-username/3xui-bot/issues)
2. Создайте новый Issue
3. Обратитесь в Telegram: @your_support

## 🙏 Благодарности

- [3X-UI](https://github.com/MHSanaei/3x-ui) - VPN панель
- [Telegram Bot API](https://core.telegram.org/bots/api) - Telegram API
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Архитектурный подход
