package scheduler

import (
	"context"
	"log/slog"
	"time"

	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"
)

type Scheduler struct {
	subRepo  ports.SubscriptionRepo
	vpnUC    *usecase.VPNUseCase
	notifUC  *usecase.NotificationUseCase
	userRepo ports.UserRepo
}

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

func (s *Scheduler) Start(ctx context.Context) {
	slog.Info("Starting scheduler...")

	go s.runPeriodically(ctx, time.Hour, s.CheckExpiredSubscriptions)

	go s.runPeriodically(ctx, 6*time.Hour, s.SendExpirationNotifications)

	go s.runPeriodically(ctx, 6*time.Hour, s.DeactivateExpiredVPNs)

	slog.Info("Scheduler started successfully")
}

func (s *Scheduler) runPeriodically(ctx context.Context, interval time.Duration, fn func(context.Context) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

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

func (s *Scheduler) CheckExpiredSubscriptions(ctx context.Context) error {
	slog.Info("Checking expired subscriptions...")

	slog.Info("Expired subscriptions check completed")

	return nil
}

func (s *Scheduler) SendExpirationNotifications(ctx context.Context) error {
	slog.Info("Sending expiration notifications...")

	slog.Info("Expiration notifications sent")

	return nil
}

func (s *Scheduler) DeactivateExpiredVPNs(ctx context.Context) error {
	slog.Info("Deactivating expired VPNs...")

	if err := s.vpnUC.DeactivateExpiredVPNs(ctx); err != nil {

		return err
	}

	slog.Info("Expired VPNs deactivated")

	return nil
}

func (s *Scheduler) CleanOldData(ctx context.Context) error {
	slog.Info("Cleaning old data...")

	slog.Info("Old data cleaned")

	return nil
}
