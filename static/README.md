# Статические файлы

Эта директория содержит статические файлы для Telegram бота.

## Структура

```
static/
├── README.md              # Этот файл
├── images/                # Изображения
│   ├── README.md         # Инструкции по изображениям
│   ├── referral_ranking.png  # PNG изображение для реферального рейтинга
│   ├── referral_ranking.svg  # SVG шаблон (если нужен)
│   └── create_simple_png.py  # Скрипт создания простого PNG
└── assets/               # Другие статические ресурсы
```

## Использование

### Изображения

Все изображения хранятся в `static/images/` и используются в коде:

```go
photoPath := "static/images/referral_ranking.png"
err := notifier.SendPhotoFromFile(ctx, chatID, photoPath, caption, keyboard)
```

### Добавление новых изображений

1. Поместите файл в `static/images/`
2. Обновите код, указав правильный путь
3. Проверьте, что файл существует перед отправкой

### Рекомендации

- Используйте PNG или JPG форматы
- Оптимизируйте размер файлов (рекомендуется < 1MB)
- Используйте понятные имена файлов
- Для продакшена рассмотрите использование file_id вместо файлов

## Создание изображений

### Простой PNG файл

```bash
python3 static/images/create_simple_png.py
```

### SVG в PNG (если есть ImageMagick)

```bash
convert static/images/referral_ranking.svg static/images/referral_ranking.png
```

### Онлайн конвертеры

- https://convertio.co/svg-png/
- https://cloudconvert.com/svg-to-png
