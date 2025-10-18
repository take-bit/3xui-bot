package app

import (
	"3xui-bot/internal/adapters/db/postgres"
	"context"
	"fmt"

	"3xui-bot/internal/adapters/bot/telegram"
	"3xui-bot/internal/adapters/db/postgres/notification"
	paymentAdapter "3xui-bot/internal/adapters/db/postgres/payment"
	"3xui-bot/internal/adapters/db/postgres/referral"
	"3xui-bot/internal/adapters/db/postgres/subscription"
	"3xui-bot/internal/adapters/db/postgres/user"
	"3xui-bot/internal/adapters/db/postgres/vpn"
	"3xui-bot/internal/adapters/marzban"
	"3xui-bot/internal/adapters/notify"
	"3xui-bot/internal/adapters/payment"
	"3xui-bot/internal/pkg/config"
	"3xui-bot/internal/pkg/logger"
	"3xui-bot/internal/ports"
	"3xui-bot/internal/scheduler"
	"3xui-bot/internal/usecase"

	transactorPgx "github.com/Thiht/transactor/pgx"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Container контейнер зависимостей приложения
type Container struct {
	Config *config.Config
	Logger *logger.Logger

	// Infrastructure
	DB         *pgxpool.Pool
	Bot        *tgbotapi.BotAPI
	DBGetter   transactorPgx.DBGetter
	UnitOfWork ports.UnitOfWork
	Clock      ports.Clock
	Marzban    ports.Marzban
	Notifier   ports.Notifier

	// Use Cases
	UserUC     *usecase.UserUseCase
	SubUC      *usecase.SubscriptionUseCase
	PaymentUC  *usecase.PaymentUseCase
	VPNUC      *usecase.VPNUseCase
	ReferralUC *usecase.ReferralUseCase
	NotifUC    *usecase.NotificationUseCase

	// Adapters
	Router    *telegram.Router
	Scheduler *scheduler.Scheduler
}

// NewContainer создает и инициализирует контейнер зависимостей
func NewContainer(ctx context.Context, configPath string) (*Container, error) {
	c := &Container{}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	c.Config = cfg

	// 2. Создаем логгер
	c.Logger = logger.New()
	c.Logger.Info("Configuration loaded successfully")

	// 3. Подключаемся к БД
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSLMode,
	)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	c.DB = pool

	transactor, dbGetter := transactorPgx.NewTransactorFromPool(pool)
	c.UnitOfWork = postgres.NewUoW(transactor)

	c.DBGetter = dbGetter
	c.Logger.Info("Database connected successfully")

	// 4. Создаем Telegram Bot API
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}
	bot.Debug = cfg.Bot.Debug
	c.Bot = bot
	c.Logger.Info("Telegram Bot API initialized: @%s", bot.Self.UserName)

	// 5. Создаем Clock
	c.Clock = &ports.SystemClock{}

	// 6. Создаем Marzban client
	c.Marzban = marzban.NewMarzbanRepository(
		cfg.Marzban.BaseURL,
		cfg.Marzban.Username,
		cfg.Marzban.Password,
	)

	// 7. Создаем Notifier
	c.Notifier = notify.NewTelegramNotifier(bot)

	// 8. Создаем репозитории
	userRepo := user.NewUser(c.DBGetter)
	subRepo := subscription.NewSubscription(c.DBGetter)
	planRepo := subscription.NewPlan(c.DBGetter)
	paymentRepo := paymentAdapter.NewPayment(c.DBGetter)
	vpnRepo := vpn.NewVPNConnection(c.DBGetter)
	referralRepo := referral.NewReferral(c.DBGetter)
	referralLinkRepo := referral.NewReferralLink(c.DBGetter)
	notifRepo := notification.NewNotification(c.DBGetter)

	// 9. Создаем базовые use cases
	c.UserUC = usecase.NewUserUseCase(userRepo, c.Clock)
	c.SubUC = usecase.NewSubscriptionUseCase(subRepo, planRepo)
	c.ReferralUC = usecase.NewReferralUseCase(referralRepo, referralLinkRepo)

	// 10. VPN UseCase
	c.VPNUC = usecase.NewVPNUseCase(vpnRepo, c.Marzban, subRepo, planRepo)

	// 11. Notification UseCase (с bot API)
	c.NotifUC = usecase.NewNotificationUseCase(notifRepo, userRepo, c.Notifier)

	// 12. Mock Payment Provider (TODO: заменить на реальный)
	paymentProvider := payment.NewMockProvider()

	// 13. Payment UseCase (оркестратор)
	c.PaymentUC = usecase.NewPaymentUseCase(
		paymentRepo,
		c.SubUC,
		c.VPNUC,
		c.NotifUC,
		paymentProvider,
	)

	// 14. Создаем роутер
	c.Router = telegram.NewRouter(
		bot,
		c.Notifier,
		c.UserUC,
		c.SubUC,
		c.PaymentUC,
		c.VPNUC,
		c.ReferralUC,
		c.NotifUC,
	)

	// 15. Создаем планировщик
	c.Scheduler = scheduler.NewScheduler(subRepo, c.VPNUC, c.NotifUC, userRepo)

	c.Logger.Info("All components initialized successfully")

	return c, nil
}

// Close закрывает все ресурсы
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
		c.Logger.Info("Database connection closed")
	}
}
