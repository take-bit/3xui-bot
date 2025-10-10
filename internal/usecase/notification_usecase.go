package usecase

import (
	"3xui-bot/internal/ports"
	"context"
	"fmt"
	"log"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NotificationUseCase use case для работы с уведомлениями
type NotificationUseCase struct {
	notifRepo ports.NotificationRepo
	userRepo  ports.UserRepo
	bot       *tgbotapi.BotAPI
}

// NewNotificationUseCase создает новый use case для уведомлений
func NewNotificationUseCase(
	notifRepo ports.NotificationRepo,
	userRepo ports.UserRepo,
	bot *tgbotapi.BotAPI,
) *NotificationUseCase {
	return &NotificationUseCase{
		notifRepo: notifRepo,
		userRepo:  userRepo,
		bot:       bot,
	}
}

// CreateNotification создает уведомление (для совместимости с DTO)
func (uc *NotificationUseCase) CreateNotification(ctx context.Context, dto CreateNotificationDTO) error {
	notification := &core.Notification{
		ID:        id.Generate(),
		UserID:    dto.UserID,
		Type:      dto.Type,
		Title:     dto.Title,
		Message:   dto.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	// Сохраняем в БД
	if err := uc.notifRepo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Отправляем через Telegram
	return uc.sendToTelegram(ctx, notification)
}

// SendNotification отправляет уведомление пользователю
func (uc *NotificationUseCase) SendNotification(ctx context.Context, dto SendNotificationDTO) error {
	newNotif := &core.Notification{
		ID:        id.Generate(),
		UserID:    dto.UserID,
		Type:      string(dto.Type),
		Title:     dto.Title,
		Message:   dto.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	// Сохраняем в БД
	if err := uc.notifRepo.CreateNotification(ctx, newNotif); err != nil {
		return err
	}

	// Отправляем через Telegram
	return uc.sendToTelegram(ctx, newNotif)
}

// sendToTelegram отправляет уведомление через Telegram Bot API
func (uc *NotificationUseCase) sendToTelegram(ctx context.Context, notification *core.Notification) error {
	// Получаем пользователя
	user, err := uc.userRepo.GetUserByID(ctx, notification.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем что пользователь не заблокирован
	if user.IsBlocked {
		return fmt.Errorf("user is blocked")
	}

	// Формируем сообщение
	message := fmt.Sprintf("📢 *%s*\n\n%s", notification.Title, notification.Message)

	// Отправляем сообщение
	msg := tgbotapi.NewMessage(user.TelegramID, message)
	msg.ParseMode = "Markdown"

	if _, err := uc.bot.Send(msg); err != nil {
		log.Printf("Failed to send notification to user %d: %v", user.TelegramID, err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Notification sent to user %d: %s", user.TelegramID, notification.Title)

	return nil
}

// SendBulkNotification отправляет уведомление нескольким пользователям
func (uc *NotificationUseCase) SendBulkNotification(ctx context.Context, userIDs []int64, notifType core.NotificationType, title, message string) error {
	for _, userID := range userIDs {
		dto := SendNotificationDTO{
			UserID:  userID,
			Type:    notifType,
			Title:   title,
			Message: message,
		}

		err := uc.SendNotification(ctx, dto)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserNotifications получает все уведомления пользователя
func (uc *NotificationUseCase) GetUserNotifications(ctx context.Context, userID int64) ([]*core.Notification, error) {
	return uc.notifRepo.GetNotificationsByUserID(ctx, userID)
}

// GetUnreadNotifications получает непрочитанные уведомления пользователя
func (uc *NotificationUseCase) GetUnreadNotifications(ctx context.Context, userID int64) ([]*core.Notification, error) {
	return uc.notifRepo.GetUnreadNotificationsByUserID(ctx, userID)
}

// MarkAsRead отмечает уведомление как прочитанное
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, userID int64, notificationID string) error {
	notif, err := uc.notifRepo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notif.UserID != userID {
		return ErrUnauthorized
	}

	notif.IsRead = true

	return uc.notifRepo.UpdateNotification(ctx, notif)
}

// MarkAllAsRead отмечает все уведомления пользователя как прочитанные
func (uc *NotificationUseCase) MarkAllAsRead(ctx context.Context, userID int64) error {
	notifications, err := uc.notifRepo.GetUnreadNotificationsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, notif := range notifications {
		notif.IsRead = true

		err = uc.notifRepo.UpdateNotification(ctx, notif)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteNotification удаляет уведомление
func (uc *NotificationUseCase) DeleteNotification(ctx context.Context, userID int64, notificationID string) error {
	notif, err := uc.notifRepo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notif.UserID != userID {
		return ErrUnauthorized
	}

	return uc.notifRepo.DeleteNotification(ctx, notificationID)
}
