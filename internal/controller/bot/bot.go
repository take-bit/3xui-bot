package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"3xui-bot/internal/config"
	"3xui-bot/internal/controller/bot/handlers"
	"3xui-bot/internal/controller/bot/middleware"
	"3xui-bot/internal/interfaces"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot представляет Telegram бота
type Bot struct {
	api            *tgbotapi.BotAPI
	config         *config.Config
	useCaseManager *usecase.UseCaseManager
	handlers       map[string]Handler
	middlewares    []Middleware
}

// Handler представляет обработчик команды или сообщения
type Handler = interfaces.HandlerInterface

// Middleware представляет middleware для обработки обновлений
type Middleware interface {
	Process(ctx context.Context, update tgbotapi.Update, next func(ctx context.Context, update tgbotapi.Update) error) error
}

// NewBot создает новый экземпляр Telegram бота
func NewBot(cfg *config.Config, useCaseManager *usecase.UseCaseManager) (*Bot, error) {
	// Создаем бота с токеном
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// Настраиваем режим отладки
	bot.Debug = cfg.LogLevel == "debug"

	// Создаем экземпляр бота
	telegramBot := &Bot{
		api:            bot,
		config:         cfg,
		useCaseManager: useCaseManager,
		handlers:       make(map[string]Handler),
		middlewares:    make([]Middleware, 0),
	}

	// Регистрируем обработчики
	telegramBot.registerHandlers()

	// Регистрируем middleware
	telegramBot.registerMiddlewares()

	return telegramBot, nil
}

// Start запускает бота
func (b *Bot) Start(ctx context.Context) error {
	log.Printf("Starting bot @%s", b.api.Self.UserName)

	// Настраиваем webhook или polling
	if b.config.BotURL != "" {
		return b.startWebhook(ctx)
	}
	return b.startPolling(ctx)
}

// Stop останавливает бота
func (b *Bot) Stop(ctx context.Context) error {
	log.Println("Stopping bot...")

	// Останавливаем Use Case Manager
	err := b.useCaseManager.Shutdown()
	if err != nil {
		log.Printf("Failed to shutdown use case manager: %v", err)
	}

	return nil
}

// startWebhook запускает бота в режиме webhook
func (b *Bot) startWebhook(ctx context.Context) error {
	// Устанавливаем webhook
	webhook, err := tgbotapi.NewWebhook(b.config.BotURL)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	_, err = b.api.Request(webhook)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	// Получаем информацию о webhook
	info, err := b.api.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Webhook last error: %s", info.LastErrorMessage)
	}

	log.Printf("Bot started with webhook: %s", b.config.BotURL)
	return nil
}

// startPolling запускает бота в режиме polling
func (b *Bot) startPolling(ctx context.Context) error {
	// Настраиваем polling
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot started with polling")

	// Обрабатываем обновления
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			go b.handleUpdate(ctx, update)
		}
	}
}

// handleUpdate обрабатывает обновление от Telegram
func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Применяем middleware
	err := b.processMiddlewares(ctx, update, func(ctx context.Context, update tgbotapi.Update) error {
		return b.processUpdate(ctx, update)
	})

	if err != nil {
		log.Printf("Failed to process update: %v", err)
	}
}

// processMiddlewares применяет все middleware
func (b *Bot) processMiddlewares(ctx context.Context, update tgbotapi.Update, next func(ctx context.Context, update tgbotapi.Update) error) error {
	if len(b.middlewares) == 0 {
		return next(ctx, update)
	}

	return b.middlewares[0].Process(ctx, update, func(ctx context.Context, update tgbotapi.Update) error {
		return b.processMiddlewares(ctx, update, next)
	})
}

// processUpdate обрабатывает обновление
func (b *Bot) processUpdate(ctx context.Context, update tgbotapi.Update) error {
	// Определяем тип обновления
	if update.Message != nil {
		return b.handleMessage(ctx, update)
	}

	if update.CallbackQuery != nil {
		return b.handleCallbackQuery(ctx, update)
	}

	return nil
}

// handleMessage обрабатывает сообщение
func (b *Bot) handleMessage(ctx context.Context, update tgbotapi.Update) error {
	message := update.Message

	// Определяем команду
	command := ""
	if message.IsCommand() {
		command = message.Command()
	} else if message.Text != "" {
		// Обрабатываем текстовые сообщения
		command = "text"
	}

	// Ищем обработчик
	handler, exists := b.handlers[command]
	if !exists {
		// Используем обработчик по умолчанию
		handler = b.handlers["default"]
	}

	if handler != nil {
		return handler.Handle(ctx, update)
	}

	return nil
}

// handleCallbackQuery обрабатывает callback query
func (b *Bot) handleCallbackQuery(ctx context.Context, update tgbotapi.Update) error {
	// Ищем обработчик для callback query
	handler, exists := b.handlers["callback"]
	if !exists {
		return nil
	}

	return handler.Handle(ctx, update)
}

// registerHandlers регистрирует все обработчики
func (b *Bot) registerHandlers() {
	// Создаем базовый обработчик
	baseHandler := NewBaseHandler(b.useCaseManager, b)

	// Создаем обработчики
	startHandler := handlers.NewStartHandler(b.useCaseManager)
	startHandler.BaseHandlerInterface = baseHandler

	helpHandler := handlers.NewHelpHandler(b.useCaseManager)
	helpHandler.BaseHandlerInterface = baseHandler

	profileHandler := handlers.NewProfileHandler(b.useCaseManager)
	profileHandler.BaseHandlerInterface = baseHandler

	subscriptionHandler := handlers.NewSubscriptionHandler(b.useCaseManager)
	subscriptionHandler.BaseHandlerInterface = baseHandler

	vpnHandler := handlers.NewVPNHandler(b.useCaseManager)
	vpnHandler.BaseHandlerInterface = baseHandler

	paymentHandler := handlers.NewPaymentHandler(b.useCaseManager)
	paymentHandler.BaseHandlerInterface = baseHandler

	promocodeHandler := handlers.NewPromocodeHandler(b.useCaseManager)
	promocodeHandler.BaseHandlerInterface = baseHandler

	referralHandler := handlers.NewReferralHandler(b.useCaseManager)
	referralHandler.BaseHandlerInterface = baseHandler

	settingsHandler := handlers.NewSettingsHandler(b.useCaseManager)
	settingsHandler.BaseHandlerInterface = baseHandler

	textHandler := handlers.NewTextHandler(b.useCaseManager)
	textHandler.BaseHandlerInterface = baseHandler

	callbackHandler := handlers.NewCallbackHandler(b.useCaseManager)
	callbackHandler.BaseHandlerInterface = baseHandler

	defaultHandler := handlers.NewDefaultHandler(b.useCaseManager)
	defaultHandler.BaseHandlerInterface = baseHandler

	// Регистрируем обработчики
	b.handlers["start"] = startHandler
	b.handlers["help"] = helpHandler
	b.handlers["profile"] = profileHandler
	b.handlers["subscription"] = subscriptionHandler
	b.handlers["vpn"] = vpnHandler
	b.handlers["payment"] = paymentHandler
	b.handlers["promocode"] = promocodeHandler
	b.handlers["referral"] = referralHandler
	b.handlers["settings"] = settingsHandler
	b.handlers["text"] = textHandler
	b.handlers["callback"] = callbackHandler
	b.handlers["default"] = defaultHandler
}

// registerMiddlewares регистрирует все middleware
func (b *Bot) registerMiddlewares() {
	// Создаем middleware
	loggingMiddleware := middleware.NewLoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(b.useCaseManager)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware()

	// Регистрируем middleware
	b.middlewares = append(b.middlewares, loggingMiddleware)
	b.middlewares = append(b.middlewares, authMiddleware)
	b.middlewares = append(b.middlewares, rateLimitMiddleware)
}

// SendMessage отправляет сообщение пользователю
func (b *Bot) SendMessage(ctx context.Context, chatID int64, text string, replyMarkup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if replyMarkup != nil {
		if keyboard, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = keyboard
		}
	}

	_, err := b.api.Send(msg)
	return err
}

// SendPhoto отправляет фото пользователю
func (b *Bot) SendPhoto(ctx context.Context, chatID int64, photo tgbotapi.FileBytes, caption string, replyMarkup interface{}) error {
	msg := tgbotapi.NewPhoto(chatID, photo)
	msg.Caption = caption
	msg.ParseMode = tgbotapi.ModeHTML

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	_, err := b.api.Send(msg)
	return err
}

// AnswerCallbackQuery отвечает на callback query
func (b *Bot) AnswerCallbackQuery(ctx context.Context, callbackQueryID string, text string, showAlert bool) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	callback.ShowAlert = showAlert

	_, err := b.api.Request(callback)
	return err
}

// EditMessageText редактирует текст сообщения
func (b *Bot) EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if replyMarkup != nil {
		if keyboard, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = &keyboard
		}
	}

	_, err := b.api.Send(msg)
	return err
}

// DeleteMessage удаляет сообщение
func (b *Bot) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := b.api.Send(deleteMsg)
	return err
}
