package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RateLimitMiddleware представляет middleware для ограничения частоты запросов
type RateLimitMiddleware struct {
	requests map[int64][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimitMiddleware создает новый middleware для ограничения частоты запросов
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		requests: make(map[int64][]time.Time),
		limit:    10,          // максимум 10 запросов
		window:   time.Minute, // в минуту
	}
}

// Process обрабатывает обновление с проверкой лимита запросов
func (m *RateLimitMiddleware) Process(ctx context.Context, update tgbotapi.Update, next func(ctx context.Context, update tgbotapi.Update) error) error {
	userID := m.getUserID(update)
	if userID == 0 {
		return next(ctx, update)
	}

	// Проверяем лимит запросов
	if !m.checkRateLimit(userID) {
		// Превышен лимит запросов
		chatID := m.getChatID(update)
		message := "⚠️ Слишком много запросов. Подождите минуту и попробуйте снова."

		// TODO: Отправить сообщение пользователю
		_ = chatID
		_ = message

		return fmt.Errorf("rate limit exceeded for user %d", userID)
	}

	// Продолжаем обработку
	return next(ctx, update)
}

// checkRateLimit проверяет лимит запросов для пользователя
func (m *RateLimitMiddleware) checkRateLimit(userID int64) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()

	// Получаем список запросов пользователя
	requests, exists := m.requests[userID]
	if !exists {
		requests = make([]time.Time, 0)
	}

	// Удаляем старые запросы (старше окна)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range requests {
		if now.Sub(reqTime) < m.window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Проверяем лимит
	if len(validRequests) >= m.limit {
		return false
	}

	// Добавляем новый запрос
	validRequests = append(validRequests, now)
	m.requests[userID] = validRequests

	return true
}

// getUserID получает ID пользователя из обновления
func (m *RateLimitMiddleware) getUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

// getChatID получает ID чата из обновления
func (m *RateLimitMiddleware) getChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return 0
}
