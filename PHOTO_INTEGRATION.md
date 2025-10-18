# Интеграция фото в Telegram бот

## Обзор

Добавлена полная поддержка отправки изображений в Telegram бот с тремя способами отправки:

1. **По file_id** - для уже загруженных в Telegram изображений
2. **Из файла** - для локальных файлов
3. **Из io.Reader** - для изображений из памяти/потока

## Реализованные компоненты

### 1. Notifier интерфейс (`internal/ports/notifier.go`)

```go
type Notifier interface {
    // Основные методы
    Send(ctx context.Context, chatID int64, text string, markup interface{}) error
    SendWithParseMode(ctx context.Context, chatID int64, text string, parseMode string, markup interface{}) error
    EditMessage(ctx context.Context, chatID int64, messageID int, text string, markup interface{}) error
    DeleteMessage(ctx context.Context, chatID int64, messageID int) error
    
    // Методы для работы с фото
    SendPhoto(ctx context.Context, chatID int64, photoFileID string, caption string, markup interface{}) error
    SendPhotoFromReader(ctx context.Context, chatID int64, photoReader io.Reader, caption string, markup interface{}) error
    SendPhotoFromFile(ctx context.Context, chatID int64, photoPath string, caption string, markup interface{}) error
    EditMessagePhoto(ctx context.Context, chatID int64, messageID int, photoFileID string, caption string, markup interface{}) error
}
```

### 2. Telegram Notifier (`internal/adapters/notify/tg_notifier.go`)

Реализация всех методов интерфейса Notifier с поддержкой:
- Отправки фото по file_id
- Отправки фото из файла
- Отправки фото из io.Reader
- Редактирования подписи к фото
- Удаления сообщений

### 3. Пример использования в хендлере

```go
func (h *CallbackHandler) handleReferralRanking(ctx context.Context, chatID int64, messageID int) error {
    caption := ui.GetReferralRankingText()
    keyboard := ui.GetReferralRankingKeyboard()
    
    // Отправка фото из файла
    photoPath := "assets/images/referral_ranking.png"
    err := h.notifier.SendPhotoFromFile(ctx, chatID, photoPath, caption, keyboard)
    
    if err != nil {
        // Fallback на обычное сообщение
        return h.editMessageWithID(chatID, messageID, caption, keyboard)
    }
    
    // Удаляем предыдущее сообщение
    h.notifier.DeleteMessage(ctx, chatID, messageID)
    return nil
}
```

### 4. Use Case методы

В `NotificationUseCase` добавлены методы:
- `SendNotificationWithPhoto()` - отправка уведомления с фото
- `SendReferralRankingPhoto()` - отправка фото реферального рейтинга

## Способы использования

### 1. Отправка по file_id (рекомендуется для продакшена)

```go
photoFileID := "AgACAgIAAxkBAAIBY2Y..." // file_id из Telegram
err := notifier.SendPhoto(ctx, chatID, photoFileID, caption, keyboard)
```

**Преимущества:**
- Быстрая отправка (фото уже в Telegram)
- Не расходует трафик на загрузку
- Надежно работает

**Как получить file_id:**
1. Отправить изображение боту в личные сообщения
2. Получить file_id из webhook или getUpdates
3. Сохранить в базе данных

### 2. Отправка из файла (для разработки)

```go
photoPath := "static/images/referral_ranking.png"
err := notifier.SendPhotoFromFile(ctx, chatID, photoPath, caption, keyboard)
```

**Преимущества:**
- Простота разработки
- Контроль над файлами

**Недостатки:**
- Медленная отправка (нужно загружать файл)
- Расходует трафик

### 3. Отправка из io.Reader (для динамических изображений)

```go
photoReader := bytes.NewReader(imageBytes)
err := notifier.SendPhotoFromReader(ctx, chatID, photoReader, caption, keyboard)
```

**Применение:**
- Генерация изображений на лету
- Обработка изображений из API
- Создание графиков/диаграмм

## Структура файлов

```
static/
├── README.md              # Общая документация по статическим файлам
└── images/
    ├── README.md              # Инструкции по использованию изображений
    ├── referral_ranking.svg   # SVG шаблон для реферального рейтинга
    ├── referral_ranking.png   # PNG версия (создается скриптом)
    ├── create_simple_png.py   # Скрипт создания простого PNG
    └── convert.sh             # Скрипт конвертации SVG в PNG (если есть ImageMagick)
```

## Примеры интеграции

### 1. Реферальный рейтинг с фото

Добавлена кнопка "🏆 Рейтинг" в реферальное меню, которая отправляет фото с рейтингом пользователей.

### 2. Уведомления с изображениями

Возможность отправлять уведомления с изображениями через NotificationUseCase.

### 3. Fallback механизм

Если не удается отправить фото, автоматически отправляется обычное текстовое сообщение.

## Рекомендации

1. **Для продакшена** используйте file_id - это самый быстрый и надежный способ
2. **Для разработки** можно использовать локальные файлы
3. **Всегда предусматривайте fallback** на текстовое сообщение
4. **Оптимизируйте изображения** - размер не должен превышать 10MB
5. **Используйте подходящие форматы** - JPG для фото, PNG для графики

## Тестирование

1. Создайте PNG файл: `python3 static/images/create_simple_png.py`
2. Запустите бота: `make run`
3. Нажмите "🏆 Рейтинг" в реферальном меню
4. Проверьте, что отправляется фото с подписью и кнопками

## Структура проекта

```
3xui-bot/
├── static/                    # Статические файлы
│   ├── README.md
│   └── images/
│       ├── README.md
│       ├── referral_ranking.png
│       └── create_simple_png.py
├── internal/
│   ├── ports/
│   │   └── notifier.go        # Интерфейс для отправки уведомлений
│   ├── adapters/
│   │   ├── notify/
│   │   │   └── tg_notifier.go # Реализация Telegram notifier
│   │   └── bot/telegram/
│   │       └── handlers/
│   │           └── callback.go # Обработчики с примерами фото
│   └── usecase/
│       └── notification_usecase.go # Use case с методами для фото
└── PHOTO_INTEGRATION.md       # Эта документация
```
