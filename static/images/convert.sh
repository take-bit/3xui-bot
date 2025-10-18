#!/bin/bash

# Скрипт для конвертации SVG в PNG
# Требует ImageMagick: brew install imagemagick (macOS) или apt-get install imagemagick (Ubuntu)

echo "Конвертируем SVG в PNG..."

if command -v convert &> /dev/null; then
    convert referral_ranking.svg referral_ranking.png
    echo "✅ referral_ranking.png создан"
else
    echo "❌ ImageMagick не найден. Установите его:"
    echo "   macOS: brew install imagemagick"
    echo "   Ubuntu: sudo apt-get install imagemagick"
    echo ""
    echo "Или используйте онлайн конвертер: https://convertio.co/svg-png/"
fi
