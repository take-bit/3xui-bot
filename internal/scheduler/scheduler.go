package scheduler

import (
	"context"
	"log/slog"
	"time"

	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"
)

// Scheduler управляет фоновыми задачами
type Scheduler struct {
	subRepo  ports.SubscriptionRepo
	vpnUC    *usecase.VPNUseCase
	notifUC  *usecase.NotificationUseCase
	userRepo ports.UserRepo
}

// NewScheduler создает новый планировщик
func NewScheduler(
	subRepo ports.SubscriptionRepo,
	vpnUC *usecase.VPNUseCase,
	notifUC *usecase.NotificationUseCase,
	userRepo ports.UserRepo,
) *Scheduler {
	return &Scheduler{
		subRepo:  subRepo,
		vpnUC:    vpnUC,
		notifUC:  notifUC,
		userRepo: userRepo,
	}
}

// Start запускает все фоновые задачи
func (s *Scheduler) Start(ctx context.Context) {
	slog.Info("Starting scheduler...")

	// Проверка истекших подписок каждый час
	go s.runPeriodically(ctx, time.Hour, s.CheckExpiredSubscriptions)

	// Отправка уведомлений за 3 дня до окончания (каждые 6 часов)
	go s.runPeriodically(ctx, 6*time.Hour, s.SendExpirationNotifications)

	// Деактивация просроченных VPN (каждые 6 часов)
	go s.runPeriodically(ctx, 6*time.Hour, s.DeactivateExpiredVPNs)

	slog.Info("Scheduler started successfully")
}

// runPeriodically запускает функцию периодически
func (s *Scheduler) runPeriodically(ctx context.Context, interval time.Duration, fn func(context.Context) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Запускаем сразу при старте
	if err := fn(ctx); err != nil {
		slog.Error("Error in scheduled job", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping scheduled job...")
			return
		case <-ticker.C:
			if err := fn(ctx); err != nil {
				slog.Error("Error in scheduled job", "error", err)
			}
		}
	}
}

// CheckExpiredSubscriptions проверяет и деактивирует истекшие подписки
func (s *Scheduler) CheckExpiredSubscriptions(ctx context.Context) error {
	slog.Info("Checking expired subscriptions...")

	// TODO: Реализовать получение всех активных подписок
	// и деактивацию истекших

	// Примерная логика:
	// 1. Получить все активные подписки
	// 2. Проверить дату окончания
	// 3. Если истекла - деактивировать
	// 4. Деактивировать связанные VPN

	slog.Info("Expired subscriptions check completed")
	return nil
}

// SendExpirationNotifications отправляет уведомления об истечении подписки
func (s *Scheduler) SendExpirationNotifications(ctx context.Context) error {
	slog.Info("Sending expiration notifications...")

	// TODO: Реализовать логику
	// 1. Найти подписки, истекающие через 3 дня
	// 2. Проверить, не отправляли ли уже уведомление
	// 3. Отправить уведомление через NotificationUseCase

	// Примерная логика:
	/*
		threeDaysLater := time.Now().AddDate(0, 0, 3)
		subscriptions := // получить подписки, истекающие threeDaysLater

		for _, sub := range subscriptions {
			notification := &domain.Notification{
				UserID: sub.UserID,
				Type:   "subscription_expiring",
				Title:  "Подписка заканчивается",
				Message: fmt.Sprintf("Ваша подписка заканчивается %s", sub.EndDate.Format("02.01.2006")),
			}
			s.notifUC.CreateNotification(ctx, notification)
		}
	*/

	slog.Info("Expiration notifications sent")
	return nil
}

// DeactivateExpiredVPNs деактивирует истекшие VPN подключения
func (s *Scheduler) DeactivateExpiredVPNs(ctx context.Context) error {
	slog.Info("Deactivating expired VPNs...")

	// Используем VPNUseCase для деактивации
	if err := s.vpnUC.DeactivateExpiredVPNs(ctx); err != nil {
		return err
	}

	slog.Info("Expired VPNs deactivated")
	return nil
}

// CleanOldData очищает старые данные (опционально)
func (s *Scheduler) CleanOldData(ctx context.Context) error {
	slog.Info("Cleaning old data...")

	// TODO: Реализовать очистку:
	// - Старых уведомлений (> 30 дней)
	// - Отмененных платежей (> 90 дней)
	// - Неактивных VPN (> 180 дней)

	slog.Info("Old data cleaned")
	return nil
}
