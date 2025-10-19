package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/handlers"
	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router маршрутизатор для Telegram бота
type Router struct {
	bot      *tgbotapi.BotAPI
	notifier ports.Notifier

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
	notifier ports.Notifier,
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
) *Router {
	r := &Router{
		bot:        bot,
		notifier:   notifier,
		userUC:     userUC,
		subUC:      subUC,
		paymentUC:  paymentUC,
		vpnUC:      vpnUC,
		referralUC: referralUC,
		notifUC:    notifUC,
	}

	// Инициализируем handlers
	r.startHandler = handlers.NewStartHandler(bot, notifier, userUC, subUC)
	r.callbackHandler = handlers.NewCallbackHandler(r) // Передаем Router как controller
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

	// Обработка pre-checkout query (для Stars)
	if update.PreCheckoutQuery != nil {
		return r.handlePreCheckout(ctx, update.PreCheckoutQuery)
	}

	// Обработка успешного платежа (для Stars)
	if update.Message != nil && update.Message.SuccessfulPayment != nil {
		return r.handleSuccessfulPayment(ctx, update.Message)
	}

	// Обработка обычных сообщений (не команд)
	if update.Message != nil && update.Message.Text != "" {
		return r.handleUnknownMessage(ctx, update.Message)
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
		return r.handleUnknownCommand(ctx, message)
	}
}

// handleCallback обрабатывает callback queries
func (r *Router) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	// Ранний ACK (отправляем сразу, чтобы убрать "часики" в Telegram)
	if botPort, ok := r.notifier.(ports.BotPort); ok {
		_ = botPort.AnswerCallback(ctx, callback.ID, "", false)
	}

	// Создаем Update из callback
	update := tgbotapi.Update{
		CallbackQuery: callback,
	}
	return r.callbackHandler.Handle(ctx, update)
}

// handleStart обрабатывает команду /start
func (r *Router) handleStart(ctx context.Context, message *tgbotapi.Message) error {
	return r.startHandler.Handle(ctx, message)
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

// handlePreCheckout обрабатывает pre-checkout query для Stars
func (r *Router) handlePreCheckout(ctx context.Context, query *tgbotapi.PreCheckoutQuery) error {
	slog.Info("Pre-checkout query received",
		"query_id", query.ID,
		"user_id", query.From.ID,
		"currency", query.Currency,
		"total_amount", query.TotalAmount,
		"payload", query.InvoicePayload)

	// Всегда подтверждаем платеж
	config := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
	}

	_, err := r.bot.Request(config)
	if err != nil {
		slog.Error("Failed to answer pre-checkout query", "error", err)
	}
	return err
}

// handleSuccessfulPayment обрабатывает успешный платеж Stars
func (r *Router) handleSuccessfulPayment(ctx context.Context, message *tgbotapi.Message) error {
	payment := message.SuccessfulPayment
	userID := message.From.ID
	chatID := message.Chat.ID

	slog.Info("Successful payment received",
		"user_id", userID,
		"currency", payment.Currency,
		"total_amount", payment.TotalAmount,
		"payload", payment.InvoicePayload,
		"telegram_payment_charge_id", payment.TelegramPaymentChargeID)

	// Парсим payload чтобы получить plan_id
	// Формат: plan_{planID}_user_{userID}
	var planID string
	if _, err := fmt.Sscanf(payment.InvoicePayload, "plan_%s", &planID); err == nil {
		// Убираем суффикс "_user_..."
		if idx := strings.Index(planID, "_user_"); idx > 0 {
			planID = planID[:idx]
		}
	} else {
		slog.Error("Failed to parse payload", "payload", payment.InvoicePayload)
		r.notifier.Send(ctx, chatID, "❌ Ошибка обработки платежа. Обратитесь в поддержку.", nil)
		return err
	}

	// Получаем план
	plan, err := r.subUC.GetPlanByID(ctx, planID)
	if err != nil {
		slog.Error("Failed to get plan", "plan_id", planID, "error", err)
		r.notifier.Send(ctx, chatID, "❌ План не найден. Обратитесь в поддержку.", nil)
		return err
	}

	slog.Info("Creating subscription after successful Stars payment",
		"user_id", userID,
		"plan_id", planID,
		"charge_id", payment.TelegramPaymentChargeID)

	// Создаем подписку
	subscription, err := r.subUC.CreateSubscription(ctx, usecase.CreateSubscriptionDTO{
		UserID:    userID,
		PlanID:    plan.ID,
		Name:      fmt.Sprintf("%s (Stars)", plan.Name),
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, plan.Days),
		IsActive:  true,
	})
	if err != nil {
		slog.Error("Failed to create subscription after Stars payment",
			"error", err,
			"user_id", userID,
			"plan_id", planID)
		r.notifier.Send(ctx, chatID, "❌ Не удалось создать подписку. Деньги будут возвращены. Обратитесь в поддержку.", nil)
		return err
	}

	slog.Info("Subscription created successfully",
		"subscription_id", subscription.ID,
		"user_id", userID,
		"plan_id", planID,
		"end_date", subscription.EndDate.Format("2006-01-02 15:04:05"))

	// Создаем VPN конфигурацию
	vpnConnection, err := r.vpnUC.CreateVPNForSubscription(ctx, userID, subscription.ID)
	if err != nil {
		slog.Error("Failed to create VPN", "error", err)
		text := fmt.Sprintf(`💎 Оплата Stars - Успешно! ✅

📦 План: %s
💎 Оплачено: %d Stars

⚠️ Подписка создана, но не удалось создать VPN конфигурацию.
Обратитесь в поддержку.`, plan.Name, payment.TotalAmount)

		r.notifier.Send(ctx, chatID, text, ui.GetMainMenuWithProfileKeyboard(true))
		return err
	}

	slog.Info("VPN created successfully",
		"vpn_id", vpnConnection.ID,
		"marzban_username", vpnConnection.MarzbanUsername,
		"subscription_id", subscription.ID)

	// Отправляем сообщение об успехе
	text := fmt.Sprintf(`🎉 Оплата Stars завершена успешно!

📦 План: %s
💎 Оплачено: %d Stars (%.0f₽)
⏰ Длительность: %d дней
📅 Действует до: %s

✅ Подписка активирована
🔑 VPN ключ создан: %s

Перейдите в "💳 Мои подписки" для получения конфигурации и настройки VPN.`,
		plan.Name,
		payment.TotalAmount,
		plan.Price,
		plan.Days,
		subscription.EndDate.Format("02.01.2006 15:04"),
		vpnConnection.Name)

	keyboard := ui.GetMainMenuWithProfileKeyboard(true)

	slog.Info("Sending success message to user", "user_id", userID)

	return r.notifier.SendWithParseMode(ctx, chatID, text, "HTML", keyboard)
}

// handleUnknownCommand обрабатывает неизвестные команды
func (r *Router) handleUnknownCommand(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	chatID := message.Chat.ID
	command := message.Command()

	slog.Info("Unknown command received",
		"user_id", userID,
		"command", command,
		"chat_id", chatID)

	text := ui.GetUnknownCommandText()
	keyboard := ui.GetUnknownCommandKeyboard()

	return r.notifier.Send(ctx, chatID, text, keyboard)
}

// handleUnknownMessage обрабатывает неизвестные текстовые сообщения
func (r *Router) handleUnknownMessage(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	chatID := message.Chat.ID
	messageText := message.Text

	// Сначала проверяем, не находится ли пользователь в процессе переименования подписки
	handled, err := r.callbackHandler.HandleTextMessage(ctx, userID, chatID, messageText)
	if err != nil {
		slog.Error("Error handling text message", "error", err, "user_id", userID)
		return err
	}

	if handled {
		return nil // Сообщение обработано (переименование подписки)
	}

	// Если не обработано, отправляем стандартное сообщение для неизвестных команд
	slog.Info("Unknown message received",
		"user_id", userID,
		"message", messageText,
		"chat_id", chatID)

	text := ui.GetUnknownCommandText()
	keyboard := ui.GetUnknownCommandKeyboard()

	return r.notifier.Send(ctx, chatID, text, keyboard)
}

// ============================================================================
// МЕТОДЫ ДЛЯ ИНТЕРФЕЙСА КОНТРОЛЛЕРА (для CallbackHandler)
// ============================================================================

func (r *Router) Bot() *tgbotapi.BotAPI {
	return r.bot
}

func (r *Router) UserUC() *usecase.UserUseCase {
	return r.userUC
}

func (r *Router) SubUC() *usecase.SubscriptionUseCase {
	return r.subUC
}

func (r *Router) VpnUC() *usecase.VPNUseCase {
	return r.vpnUC
}

func (r *Router) PaymentUC() *usecase.PaymentUseCase {
	return r.paymentUC
}

func (r *Router) ReferralUC() *usecase.ReferralUseCase {
	return r.referralUC
}

func (r *Router) EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if replyMarkup != nil {
		if keyboard, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
			edit.ReplyMarkup = &keyboard
		}
	}
	_, err := r.bot.Send(edit)
	return err
}

func (r *Router) AnswerCallbackQuery(ctx context.Context, callbackQueryID, text string, showAlert bool) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	callback.ShowAlert = showAlert
	_, err := r.bot.Request(callback)
	return err
}

func (r *Router) SendMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := r.bot.Send(msg)
	return err
}

func (r *Router) LogError(err error, context string) {
	slog.Error("Error in context", "context", context, "error", err)
}
