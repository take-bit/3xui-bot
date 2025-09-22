package usecase

import (
	"context"
	"fmt"
	"time"

	"3xui-bot/internal/domain"
)

// NotificationUseCase представляет use case для работы с уведомлениями
type NotificationUseCase struct {
	notificationService domain.NotificationService
	userService         domain.UserService
}

// NewNotificationUseCase создает новый Notification use case
func NewNotificationUseCase(
	notificationService domain.NotificationService,
	userService domain.UserService,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationService: notificationService,
		userService:         userService,
	}
}

// SendToUser отправляет уведомление пользователю
func (uc *NotificationUseCase) SendToUser(ctx context.Context, userID int64, title, message string, isHTML bool) error {
	// 1. Проверяем, что пользователь существует
	_, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Отправляем уведомление
	err = uc.notificationService.SendToUser(ctx, userID, title, message, isHTML)
	if err != nil {
		return fmt.Errorf("failed to send notification to user: %w", err)
	}

	return nil
}

// SendToAll отправляет уведомление всем пользователям
func (uc *NotificationUseCase) SendToAll(ctx context.Context, title, message string, isHTML bool) error {
	// 1. Отправляем уведомление всем пользователям
	err := uc.notificationService.SendToAll(ctx, title, message, isHTML)
	if err != nil {
		return fmt.Errorf("failed to send notification to all users: %w", err)
	}

	return nil
}

// CreateDraft создает черновик уведомления
func (uc *NotificationUseCase) CreateDraft(ctx context.Context, notificationType domain.NotificationType, userID *int64, title, message string, isHTML bool) (*domain.Notification, error) {
	// 1. Создаем черновик
	notification, err := uc.notificationService.CreateDraft(ctx, notificationType, userID, title, message, isHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification draft: %w", err)
	}

	return notification, nil
}

// UpdateDraft обновляет черновик уведомления
func (uc *NotificationUseCase) UpdateDraft(ctx context.Context, id int64, title, message string, isHTML bool) error {
	// 1. Обновляем черновик
	err := uc.notificationService.UpdateDraft(ctx, id, title, message, isHTML)
	if err != nil {
		return fmt.Errorf("failed to update notification draft: %w", err)
	}

	return nil
}

// SendDraft отправляет черновик уведомления
func (uc *NotificationUseCase) SendDraft(ctx context.Context, id int64) error {
	// 1. Отправляем черновик
	err := uc.notificationService.SendDraft(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to send notification draft: %w", err)
	}

	return nil
}

// GetUserNotifications получает уведомления пользователя
func (uc *NotificationUseCase) GetUserNotifications(ctx context.Context, userID int64) ([]*domain.Notification, error) {
	// 1. Проверяем, что пользователь существует
	_, err := uc.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Получаем уведомления пользователя
	notifications, err := uc.notificationService.GetUserNotifications(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	return notifications, nil
}

// SendWelcomeMessage отправляет приветственное сообщение
func (uc *NotificationUseCase) SendWelcomeMessage(ctx context.Context, userID int64, firstName string) error {
	message := fmt.Sprintf("👋 Добро пожаловать, %s!\n\n🎁 Вам предоставлен пробный период на 3 дня\n\n🚀 Начните пользоваться VPN прямо сейчас!\n\n📱 Используйте команды в меню для управления подпиской", firstName)

	return uc.SendToUser(ctx, userID, "Добро пожаловать!", message, false)
}

// SendSubscriptionExpired отправляет уведомление об истечении подписки
func (uc *NotificationUseCase) SendSubscriptionExpired(ctx context.Context, userID int64) error {
	message := "⏰ Ваша подписка истекла\n\n🔄 Продлите подписку, чтобы продолжить пользоваться VPN\n\n💳 Доступные способы оплаты в меню"

	return uc.SendToUser(ctx, userID, "Подписка истекла", message, false)
}

// SendSubscriptionExpiringSoon отправляет уведомление о скором истечении подписки
func (uc *NotificationUseCase) SendSubscriptionExpiringSoon(ctx context.Context, userID int64, daysRemaining int) error {
	message := fmt.Sprintf("⚠️ Ваша подписка истекает через %d дней\n\n🔄 Продлите подписку, чтобы не прерывать доступ к VPN\n\n💳 Доступные способы оплаты в меню", daysRemaining)

	return uc.SendToUser(ctx, userID, "Подписка истекает", message, false)
}

// SendPaymentSuccess отправляет уведомление об успешном платеже
func (uc *NotificationUseCase) SendPaymentSuccess(ctx context.Context, userID int64, amount int, currency string) error {
	message := fmt.Sprintf("✅ Платеж успешно обработан!\n\n💰 Сумма: %d %s\n🎉 Ваша подписка активирована!\n\n🔗 VPN подключение готово к использованию", amount, currency)

	return uc.SendToUser(ctx, userID, "Платеж успешен", message, false)
}

// SendPaymentFailed отправляет уведомление о неудачном платеже
func (uc *NotificationUseCase) SendPaymentFailed(ctx context.Context, userID int64, reason string) error {
	message := fmt.Sprintf("❌ Платеж не прошел\n\n📋 Причина: %s\n\n🔄 Попробуйте еще раз или обратитесь в поддержку", reason)

	return uc.SendToUser(ctx, userID, "Платеж не прошел", message, false)
}

// SendVPNConnectionCreated отправляет уведомление о создании VPN подключения
func (uc *NotificationUseCase) SendVPNConnectionCreated(ctx context.Context, userID int64, configURL string, expiresAt time.Time) error {
	message := fmt.Sprintf("✅ VPN подключение создано!\n\n🔗 Конфигурация: %s\n📅 Действует до: %s\n\n📱 Скачайте конфигурацию и настройте VPN", configURL, expiresAt.Format("02.01.2006 15:04"))

	return uc.SendToUser(ctx, userID, "VPN подключение", message, false)
}

// SendVPNConnectionDeleted отправляет уведомление об удалении VPN подключения
func (uc *NotificationUseCase) SendVPNConnectionDeleted(ctx context.Context, userID int64) error {
	message := "❌ VPN подключение удалено\n\n🔄 Создайте новое подключение в меню"

	return uc.SendToUser(ctx, userID, "VPN подключение", message, false)
}

// SendPromocodeApplied отправляет уведомление о применении промокода
func (uc *NotificationUseCase) SendPromocodeApplied(ctx context.Context, userID int64, code string, daysAdded int) error {
	message := fmt.Sprintf("🎁 Промокод '%s' применен!\n\n⏰ Добавлено %d дней\n\n🎉 Продлите подписку, чтобы не потерять доступ", code, daysAdded)

	return uc.SendToUser(ctx, userID, "Промокод применен", message, false)
}

// SendReferralReward отправляет уведомление о реферальном вознаграждении
func (uc *NotificationUseCase) SendReferralReward(ctx context.Context, userID int64, amount int, referredUser string) error {
	message := fmt.Sprintf("💰 Реферальное вознаграждение!\n\n👤 За пользователя: %s\n💸 Сумма: %d руб.\n\n🎉 Продолжайте приглашать друзей!", referredUser, amount)

	return uc.SendToUser(ctx, userID, "Реферальное вознаграждение", message, false)
}

// SendMaintenanceNotification отправляет уведомление о технических работах
func (uc *NotificationUseCase) SendMaintenanceNotification(ctx context.Context, startTime, endTime time.Time) error {
	message := fmt.Sprintf("🔧 Технические работы\n\n⏰ Начало: %s\n⏰ Окончание: %s\n\n⚠️ В это время VPN может быть недоступен",
		startTime.Format("02.01.2006 15:04"), endTime.Format("02.01.2006 15:04"))

	return uc.SendToAll(ctx, "Технические работы", message, false)
}

// SendServerStatusUpdate отправляет уведомление об обновлении статуса сервера
func (uc *NotificationUseCase) SendServerStatusUpdate(ctx context.Context, serverName string, status string) error {
	var emoji string
	switch status {
	case "up":
		emoji = "✅"
	case "down":
		emoji = "❌"
	case "maintenance":
		emoji = "🔧"
	default:
		emoji = "ℹ️"
	}

	message := fmt.Sprintf("%s Сервер %s\n\n📊 Статус: %s\n\n🔗 VPN подключения работают в штатном режиме", emoji, serverName, status)

	return uc.SendToAll(ctx, "Статус сервера", message, false)
}

// SendCustomMessage отправляет кастомное сообщение
func (uc *NotificationUseCase) SendCustomMessage(ctx context.Context, userID int64, title, message string, isHTML bool) error {
	return uc.SendToUser(ctx, userID, title, message, isHTML)
}

// SendBulkMessage отправляет массовое сообщение
func (uc *NotificationUseCase) SendBulkMessage(ctx context.Context, title, message string, isHTML bool) error {
	return uc.SendToAll(ctx, title, message, isHTML)
}

// ScheduleNotification планирует отправку уведомления
func (uc *NotificationUseCase) ScheduleNotification(ctx context.Context, userID int64, title, message string, sendAt time.Time) error {
	// TODO: Реализовать планирование уведомлений
	return fmt.Errorf("scheduled notifications not implemented yet")
}

// GetNotificationStats возвращает статистику уведомлений
func (uc *NotificationUseCase) GetNotificationStats(ctx context.Context) (*NotificationStats, error) {
	// TODO: Реализовать получение статистики уведомлений
	return nil, fmt.Errorf("notification stats not implemented yet")
}

// NotificationStats представляет статистику уведомлений
type NotificationStats struct {
	TotalSent           int           `json:"total_sent"`
	SuccessfullySent    int           `json:"successfully_sent"`
	FailedSent          int           `json:"failed_sent"`
	SuccessRate         float64       `json:"success_rate"`
	AverageDeliveryTime time.Duration `json:"average_delivery_time"`
}
