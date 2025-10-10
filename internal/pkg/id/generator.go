package id

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Generate генерирует уникальный ID
func Generate() string {
	return uuid.New().String()
}

// GenerateWithPrefix генерирует ID с префиксом
func GenerateWithPrefix(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, uuid.New().String())
}

// GenerateShort генерирует короткий ID (первые 8 символов UUID)
func GenerateShort() string {
	return uuid.New().String()[:8]
}

// GenerateWithTimestamp генерирует ID с префиксом и timestamp
func GenerateWithTimestamp(prefix string) string {
	return fmt.Sprintf("%s_%d_%s", prefix, time.Now().Unix(), uuid.New().String()[:8])
}

// GenerateNumeric генерирует числовой ID на основе timestamp
func GenerateNumeric() int64 {
	return time.Now().UnixNano()
}
