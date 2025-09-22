package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"3xui-bot/internal/domain"
)

// PromocodeUseCase представляет use case для работы с промокодами
type PromocodeUseCase struct {
	promocodeService    domain.PromocodeService
	subscriptionService domain.SubscriptionService
	userService         domain.UserService
	notificationService domain.NotificationService
}

// NewPromocodeUseCase создает новый Promocode use case
func NewPromocodeUseCase(
	promocodeService domain.PromocodeService,
	subscriptionService domain.SubscriptionService,
	userService domain.UserService,
	notificationService domain.NotificationService,
) *PromocodeUseCase {
	return &PromocodeUseCase{
		promocodeService:    promocodeService,
		subscriptionService: subscriptionService,
		userService:         userService,
		notificationService: notificationService,
	}
}

// CreatePromocode создает новый промокод
func (uc *PromocodeUseCase) CreatePromocode(ctx context.Context, code string, promocodeType domain.PromocodeType, value int, usageLimit int, expiresAt *time.Time) (*domain.Promocode, error) {
	// 1. Нормализуем код промокода
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, fmt.Errorf("promocode code cannot be empty")
	}

	// 2. Создаем промокод
	promocode, err := uc.promocodeService.Create(ctx, code, promocodeType, value, usageLimit, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create promocode: %w", err)
	}

	return promocode, nil
}

// ApplyPromocode применяет промокод
func (uc *PromocodeUseCase) ApplyPromocode(ctx context.Context, userID int64, code string) (*PromocodeResult, error) {
	// 1. Нормализуем код промокода
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, fmt.Errorf("promocode code cannot be empty")
	}

	// 2. Получаем пользователя
	user, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 3. Валидируем промокод
	promocode, err := uc.promocodeService.Validate(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid promocode: %w", err)
	}

	// 4. Применяем промокод
	err = uc.promocodeService.Use(ctx, code, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to apply promocode: %w", err)
	}

	// 5. Обрабатываем результат в зависимости от типа промокода
	result := &PromocodeResult{
		Promocode: promocode,
		Success:   true,
		Message:   "",
		DaysAdded: 0,
		Discount:  0,
	}

	switch promocode.Type {
	case domain.PromocodeTypeExtraDays:
		// Продлеваем пробную подписку
		err = uc.subscriptionService.Extend(ctx, user.ID, promocode.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to extend trial subscription: %w", err)
		}
		result.DaysAdded = promocode.Value
		result.Message = fmt.Sprintf("🎁 Промокод применен!\n\n⏰ Добавлено %d дней пробного периода", promocode.Value)

	case domain.PromocodeTypeDiscount:
		// TODO: Реализовать применение скидки
		result.Discount = promocode.Value
		result.Message = fmt.Sprintf("💰 Промокод применен!\n\n💸 Скидка: %d%%", promocode.Value)

	default:
		return nil, fmt.Errorf("unknown promocode type: %s", promocode.Type)
	}

	// 6. Отправляем уведомление пользователю
	if uc.notificationService != nil {
		_ = uc.notificationService.SendToUser(ctx, user.TelegramID, "Промокод", result.Message, false)
	}

	return result, nil
}

// ValidatePromocode валидирует промокод без применения
func (uc *PromocodeUseCase) ValidatePromocode(ctx context.Context, code string) (*PromocodeValidation, error) {
	// 1. Нормализуем код промокода
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, fmt.Errorf("promocode code cannot be empty")
	}

	// 2. Валидируем промокод
	promocode, err := uc.promocodeService.Validate(ctx, code)
	if err != nil {
		return &PromocodeValidation{
			Valid:     false,
			Message:   fmt.Sprintf("❌ Промокод недействителен: %v", err),
			Promocode: nil,
		}, nil
	}

	// 3. Формируем информацию о промокоде
	var message string
	switch promocode.Type {
	case domain.PromocodeTypeExtraDays:
		message = fmt.Sprintf("✅ Промокод действителен!\n\n⏰ Добавит %d дней пробного периода", promocode.Value)
	case domain.PromocodeTypeDiscount:
		message = fmt.Sprintf("✅ Промокод действителен!\n\n💸 Скидка: %d%%", promocode.Value)
	default:
		message = "✅ Промокод действителен!"
	}

	validation := &PromocodeValidation{
		Valid:     true,
		Message:   message,
		Promocode: promocode,
	}

	return validation, nil
}

// GetActivePromocodes возвращает список активных промокодов
func (uc *PromocodeUseCase) GetActivePromocodes(ctx context.Context) ([]*domain.Promocode, error) {
	return uc.promocodeService.GetActive(ctx)
}

// DeactivatePromocode деактивирует промокод
func (uc *PromocodeUseCase) DeactivatePromocode(ctx context.Context, code string) error {
	// 1. Нормализуем код промокода
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return fmt.Errorf("promocode code cannot be empty")
	}

	// 2. Деактивируем промокод
	err := uc.promocodeService.Deactivate(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to deactivate promocode: %w", err)
	}

	return nil
}

// GetPromocodeStats возвращает статистику промокодов
func (uc *PromocodeUseCase) GetPromocodeStats(ctx context.Context) (*PromocodeStats, error) {
	// TODO: Реализовать получение статистики промокодов
	return nil, fmt.Errorf("promocode stats not implemented yet")
}

// CreateBulkPromocodes создает несколько промокодов
func (uc *PromocodeUseCase) CreateBulkPromocodes(ctx context.Context, count int, prefix string, promocodeType domain.PromocodeType, value int, usageLimit int, expiresAt *time.Time) ([]*domain.Promocode, error) {
	var promocodes []*domain.Promocode

	for i := 0; i < count; i++ {
		// Генерируем уникальный код
		code := fmt.Sprintf("%s%04d", prefix, i+1)

		promocode, err := uc.promocodeService.Create(ctx, code, promocodeType, value, usageLimit, expiresAt)
		if err != nil {
			// Логируем ошибку и продолжаем
			fmt.Printf("Failed to create promocode %s: %v\n", code, err)
			continue
		}

		promocodes = append(promocodes, promocode)
	}

	return promocodes, nil
}

// PromocodeResult представляет результат применения промокода
type PromocodeResult struct {
	Promocode *domain.Promocode `json:"promocode"`
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	DaysAdded int               `json:"days_added,omitempty"`
	Discount  int               `json:"discount,omitempty"`
}

// PromocodeValidation представляет результат валидации промокода
type PromocodeValidation struct {
	Valid     bool              `json:"valid"`
	Message   string            `json:"message"`
	Promocode *domain.Promocode `json:"promocode,omitempty"`
}

// PromocodeStats представляет статистику промокодов
type PromocodeStats struct {
	TotalPromocodes  int     `json:"total_promocodes"`
	ActivePromocodes int     `json:"active_promocodes"`
	UsedPromocodes   int     `json:"used_promocodes"`
	TotalUsage       int     `json:"total_usage"`
	AverageUsage     float64 `json:"average_usage"`
	MostPopularType  string  `json:"most_popular_type"`
}
