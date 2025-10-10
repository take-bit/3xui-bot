# 🚀 Быстрый старт

## 1. Настройка окружения

```bash
# Скопировать пример конфигурации
cp env.sample .env

# Отредактировать .env и заполнить:
# - BOT_TOKEN (получить у @BotFather)
# - DB_PASSWORD
# - MARZBAN_* (если используется)
```

## 2. Установка зависимостей

```bash
make install
```

## 3. Настройка базы данных

### Вариант 1: Локальный PostgreSQL
```bash
# Создать БД и применить миграции
make setup
```

### Вариант 2: Docker PostgreSQL
```bash
# Запустить PostgreSQL в Docker
make docker-up

# Применить миграции
make migrate-up
```

## 4. Запуск бота

### Dev режим (с auto-reload)
```bash
make dev
```

### Production режим
```bash
# Сборка
make build

# Запуск
make run
```

## 5. Полезные команды

```bash
make help           # Список всех команд
make test           # Запустить тесты
make fmt            # Форматировать код
make lint           # Проверить линтером
make clean          # Очистить build артефакты
```

## 6. Структура проекта

```
3xui-bot/
├── cmd/bot/           # Точка входа
├── internal/
│   ├── domain/       # Доменные модели
│   ├── usecase/      # Бизнес-логика
│   ├── repository/   # Репозитории
│   ├── controller/   # Telegram handlers
│   └── config/       # Конфигурация
├── migrations/       # SQL миграции
└── bin/             # Собранные бинарники
```

## 7. Troubleshooting

### Ошибка подключения к БД
```bash
# Проверить PostgreSQL
psql -U postgres -l

# Пересоздать БД
make migrate-down
```

### Бот не отвечает
- Проверить BOT_TOKEN в .env
- Проверить логи: `./bin/bot` или `make dev`

### Проблемы с миграциями
```bash
# Откатить и применить заново
make migrate-down
make migrate-up
```

## 📚 Дополнительная документация

- `README.md` - Полное руководство
- `ARCHITECTURE.md` - Описание архитектуры
- `TODO.md` - План развития
- `REFACTORING_SUMMARY.md` - История изменений
