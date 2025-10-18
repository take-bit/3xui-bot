package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/pkg/id"
	"3xui-bot/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NotificationUseCase use case для работы с уведомлениями
type NotificationUseCase struct {
	notifRepo ports.NotificationRepo
	userRepo  ports.UserRepo
	notifier  ports.Notifier
}

// NewNotificationUseCase создает новый use case для уведомлений
func NewNotificationUseCase(
	notifRepo ports.NotificationRepo,
	userRepo ports.UserRepo,
	notifier ports.Notifier,
) *NotificationUseCase {
	return &NotificationUseCase{
		notifRepo: notifRepo,
		userRepo:  userRepo,
		notifier:  notifier,
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

	// Отправляем сообщение через notifier
	if err := uc.notifier.SendWithParseMode(ctx, user.TelegramID, message, "Markdown", nil); err != nil {
		slog.Error("Failed to send notification to user", "user_id", user.TelegramID, "error", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	slog.Info("Notification sent to user", "user_id", user.TelegramID, "title", notification.Title)

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

// SendNotificationWithPhoto отправляет уведомление с фото
func (uc *NotificationUseCase) SendNotificationWithPhoto(ctx context.Context, userID int64, photoPath, caption string, keyboard interface{}) error {
	// Получаем пользователя для проверки chat_id
	user, err := uc.userRepo.GetUserByTelegramID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Отправляем фото через notifier
	if err := uc.notifier.SendPhotoFromFile(ctx, user.TelegramID, photoPath, caption, keyboard); err != nil {
		return fmt.Errorf("failed to send photo notification: %w", err)
	}

	// Сохраняем уведомление в БД
	notification := &core.Notification{
		ID:        id.Generate(),
		UserID:    userID,
		Type:      "photo_notification",
		Title:     "Уведомление с фото",
		Message:   caption,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := uc.notifRepo.CreateNotification(ctx, notification); err != nil {
		slog.Warn("Failed to save photo notification to DB", "error", err)
		// Не возвращаем ошибку, так как уведомление уже отправлено
	}

	return nil
}

// SendReferralRankingPhoto отправляет фото реферального рейтинга
func (uc *NotificationUseCase) SendReferralRankingPhoto(ctx context.Context, userID int64) error {
	caption := `🏆 Рейтинг рефералов

Здесь можно увидеть топ людей, которые пригласили наибольшее количество рефералов в сервис.

Твоё место в рейтинге:
Ты еще не приглашал пользователей в проект.

🏆 Топ-5 пригласивших:
1. 57956***** - 156 чел.
2. 80000***** - 105 чел.
3. 52587***** - 12 чел.
4. 63999***** - 7 чел.
5. 10149***** - 6 чел.`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_referrals"),
		),
	)

	return uc.SendNotificationWithPhoto(ctx, userID, "static/images/bot_banner.png", caption, keyboard)
}
