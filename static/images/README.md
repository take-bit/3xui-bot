# Изображения для бота

Эта директория содержит изображения, которые используются в Telegram боте.

## Текущие файлы

- `referral_ranking.png` - изображение для реферального рейтинга (400x300, темно-синий фон)
- `referral_ranking.svg` - SVG шаблон (если нужен)
- `create_simple_png.py` - скрипт для создания простого PNG файла

## Как добавить фото в бот

### 1. Загрузка изображения в Telegram (рекомендуется для продакшена)

1. Отправьте изображение боту в личные сообщения
2. Получите `file_id` из логов или API
3. Используйте `file_id` в коде:

```go
photoFileID := "AgACAgIAAxkBAAIBY2Y..." // file_id из базы данных
err := notifier.SendPhoto(ctx, chatID, photoFileID, caption, keyboard)
```

### 2. Использование локальных файлов (для разработки)

Поместите изображения в эту директорию и используйте путь к файлу:

```go
photoPath := "static/images/referral_ranking.png"
err := notifier.SendPhotoFromFile(ctx, chatID, photoPath, caption, keyboard)
```

### 3. Использование io.Reader (для динамических изображений)

```go
photoReader := bytes.NewReader(imageBytes)
err := notifier.SendPhotoFromReader(ctx, chatID, photoReader, caption, keyboard)
```

## Доступные методы в Notifier

- `SendPhoto()` - отправка по file_id
- `SendPhotoFromFile()` - отправка из файла  
- `SendPhotoFromReader()` - отправка из io.Reader
- `EditMessagePhoto()` - редактирование подписи к фото

## Рекомендации

- Используйте изображения размером не более 10MB
- Поддерживаемые форматы: JPG, PNG, GIF
- Для лучшего отображения используйте соотношение сторон 16:9 или 1:1
- Оптимизируйте изображения для быстрой загрузки
- Всегда предусматривайте fallback на текстовое сообщение

## Создание изображений

### Простой PNG файл

```bash
python3 create_simple_png.py
```

### SVG в PNG (если есть ImageMagick)

```bash
convert referral_ranking.svg referral_ranking.png
```

### Онлайн конвертеры

- https://convertio.co/svg-png/
- https://cloudconvert.com/svg-to-png
