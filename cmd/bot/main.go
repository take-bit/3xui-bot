package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"3xui-bot/internal/config"
	"3xui-bot/internal/controller/bot"
	"3xui-bot/internal/usecase"
	// TODO: Импортировать сервисы и репозитории
	// "3xui-bot/internal/service"
	// "3xui-bot/internal/repository"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем сигналы для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// TODO: Инициализировать сервисы и репозитории
	// userRepo := repository.NewUserRepository(db)
	// subscriptionRepo := repository.NewSubscriptionRepository(db)
	// paymentRepo := repository.NewPaymentRepository(db)
	// promocodeRepo := repository.NewPromocodeRepository(db)
	// referralRepo := repository.NewReferralRepository(db)
	// notificationRepo := repository.NewNotificationRepository(db)
	// serverRepo := repository.NewServerRepository(db)

	// TODO: Создать сервисы
	// userService := service.NewUserService(userRepo)
	// subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	// paymentService := service.NewPaymentService(paymentRepo)
	// promocodeService := service.NewPromocodeService(promocodeRepo)
	// referralService := service.NewReferralService(referralRepo)
	// notificationService := service.NewNotificationService(notificationRepo)
	// vpnService := service.NewVPNService(userRepo, subscriptionRepo, serverRepo, serverService)
	// serverService := service.NewServerService(serverRepo, serverManager)

	// TODO: Создать Use Case Manager
	// useCaseManager := usecase.NewUseCaseManager(
	//     vpnService, paymentService, userService,
	//     subscriptionService, promocodeService,
	//     referralService, notificationService, serverService,
	// )

	// Временно создаем заглушку для Use Case Manager
	useCaseManager := &usecase.UseCaseManager{}

	// Создаем Telegram бота
	telegramBot, err := bot.NewBot(cfg, useCaseManager)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Инициализируем Use Case Manager
	err = useCaseManager.Initialize(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize use case manager: %v", err)
	}

	// Запускаем бота
	log.Println("Starting Telegram bot...")
	err = telegramBot.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	// Останавливаем бота
	log.Println("Stopping Telegram bot...")
	err = telegramBot.Stop(ctx)
	if err != nil {
		log.Printf("Failed to stop bot: %v", err)
	}

	log.Println("Bot stopped successfully")
}
