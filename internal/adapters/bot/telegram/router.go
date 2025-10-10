package telegram

import (
	"context"
	"log"

	"3xui-bot/internal/adapters/bot/telegram/handlers"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router маршрутизатор для Telegram бота
type Router struct {
	bot *tgbotapi.BotAPI

	// Use cases
	userUC     *usecase.UserUseCase
	subUC      *usecase.SubscriptionUseCase
	paymentUC  *usecase.PaymentUseCase
	vpnUC      *usecase.VPNUseCase
	referralUC *usecase.ReferralUseCase
	notifUC    *usecase.NotificationUseCase

	// Handlers
	startHandler    *handlers.StartHandler
	callbackHandler *handlers.CallbackHandler
	paymentHandler  *handlers.PaymentHandler
	vpnHandler      *handlers.VPNHandler
}

// NewRouter создает новый роутер
func NewRouter(
	bot *tgbotapi.BotAPI,
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
) *Router {
	r := &Router{
		bot:        bot,
		userUC:     userUC,
		subUC:      subUC,
		paymentUC:  paymentUC,
		vpnUC:      vpnUC,
		referralUC: referralUC,
		notifUC:    notifUC,
	}

	// Инициализируем handlers
	r.paymentHandler = handlers.NewPaymentHandler(bot, paymentUC)
	r.vpnHandler = handlers.NewVPNHandler(bot, vpnUC)

	return r
}

// HandleUpdate обрабатывает обновление от Telegram
func (r *Router) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	// Обработка команд
	if update.Message != nil && update.Message.IsCommand() {
		return r.handleCommand(ctx, update.Message)
	}

	// Обработка callback queries
	if update.CallbackQuery != nil {
		return r.handleCallback(ctx, update.CallbackQuery)
	}

	return nil
}

// handleCommand обрабатывает команды
func (r *Router) handleCommand(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		return r.handleStart(ctx, message)
	case "help":
		return r.handleHelp(ctx, message)
	case "vpn":
		return r.vpnHandler.HandleShowVPNs(ctx, message.From.ID, message.Chat.ID)
	default:
		return nil
	}
}

// handleCallback обрабатывает callback queries
func (r *Router) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	// Здесь должна быть логика маршрутизации callback
	// TODO: Реализовать полную маршрутизацию
	log.Printf("Callback: %s from user %d", callback.Data, callback.From.ID)
	return nil
}

// handleStart обрабатывает команду /start
func (r *Router) handleStart(ctx context.Context, message *tgbotapi.Message) error {
	// Создаем или получаем пользователя
	// TODO: Реализовать через handlers.StartHandler

	msg := tgbotapi.NewMessage(message.Chat.ID, "Добро пожаловать в VPN бот!")
	_, err := r.bot.Send(msg)
	return err
}

// handleHelp обрабатывает команду /help
func (r *Router) handleHelp(ctx context.Context, message *tgbotapi.Message) error {
	helpText := "📖 *Помощь*\n\n" +
		"/start - Начать работу\n" +
		"/vpn - Мои VPN подключения\n" +
		"/help - Эта справка"

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	_, err := r.bot.Send(msg)
	return err
}
