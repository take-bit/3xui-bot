package handlers

import (
	"context"
	"strings"

	"3xui-bot/internal/interfaces"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackHandler обрабатывает callback query
type CallbackHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewCallbackHandler создает новый обработчик callback query
func NewCallbackHandler(useCaseManager *usecase.UseCaseManager) *CallbackHandler {
	return &CallbackHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает callback query
func (h *CallbackHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	callbackQuery := update.CallbackQuery
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)
	messageID := h.GetMessageID(update)
	data := callbackQuery.Data

	// Отвечаем на callback query
	h.AnswerCallbackQuery(ctx, callbackQuery.ID, "", false)

	// Обрабатываем различные типы callback
	switch {
	case data == "main_menu":
		return h.showMainMenu(ctx, chatID, messageID)
	case data == "profile":
		return h.showProfile(ctx, chatID, int64(messageID), userID)
	case data == "vpn":
		return h.showVPN(ctx, chatID, int64(messageID), userID)
	case data == "subscription":
		return h.showSubscription(ctx, chatID, int64(messageID), userID)
	case data == "promocode":
		return h.showPromocode(ctx, chatID, messageID)
	case data == "referral":
		return h.showReferral(ctx, chatID, int64(messageID), userID)
	case data == "settings":
		return h.showSettings(ctx, chatID, int(messageID))
	case data == "help":
		return h.showHelp(ctx, chatID, messageID)
	case strings.HasPrefix(data, "vpn_"):
		return h.handleVPNAction(ctx, chatID, int64(messageID), userID, data)
	case strings.HasPrefix(data, "subscription_"):
		return h.handleSubscriptionAction(ctx, chatID, int64(messageID), userID, data)
	case strings.HasPrefix(data, "referral_"):
		return h.handleReferralAction(ctx, chatID, int64(messageID), userID, data)
	case strings.HasPrefix(data, "settings_"):
		return h.handleSettingsAction(ctx, chatID, int64(messageID), userID, data)
	default:
		return h.showMainMenu(ctx, chatID, messageID)
	}
}

// showMainMenu показывает главное меню
func (h *CallbackHandler) showMainMenu(ctx context.Context, chatID int64, messageID int) error {
	message := `
🏠 <b>Главное меню</b>

Выберите нужный раздел:
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)

	return h.EditMessageText(ctx, chatID, messageID, message, keyboard)
}

// showProfile показывает профиль пользователя
func (h *CallbackHandler) showProfile(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	// Создаем обработчик профиля
	profileHandler := NewProfileHandler(h.UseCaseManager)
	profileHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return profileHandler.Handle(ctx, update)
}

// showVPN показывает VPN подключение
func (h *CallbackHandler) showVPN(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	vpnHandler := NewVPNHandler(h.UseCaseManager)
	vpnHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return vpnHandler.Handle(ctx, update)
}

// showSubscription показывает подписку
func (h *CallbackHandler) showSubscription(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	subscriptionHandler := NewSubscriptionHandler(h.UseCaseManager)
	subscriptionHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return subscriptionHandler.Handle(ctx, update)
}

// showPromocode показывает промокоды
func (h *CallbackHandler) showPromocode(ctx context.Context, chatID int64, messageID int) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	promocodeHandler := NewPromocodeHandler(h.UseCaseManager)
	promocodeHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return promocodeHandler.Handle(ctx, update)
}

// showReferral показывает рефералы
func (h *CallbackHandler) showReferral(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: userID},
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	referralHandler := NewReferralHandler(h.UseCaseManager)
	referralHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return referralHandler.Handle(ctx, update)
}

// showSettings показывает настройки
func (h *CallbackHandler) showSettings(ctx context.Context, chatID int64, messageID int) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	settingsHandler := NewSettingsHandler(h.UseCaseManager)
	settingsHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return settingsHandler.Handle(ctx, update)
}

// showHelp показывает помощь
func (h *CallbackHandler) showHelp(ctx context.Context, chatID int64, messageID int) error {
	// Создаем обновление для обработчика
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: chatID},
		},
	}

	helpHandler := NewHelpHandler(h.UseCaseManager)
	helpHandler.BaseHandlerInterface = h.BaseHandlerInterface
	return helpHandler.Handle(ctx, update)
}

// handleVPNAction обрабатывает действия с VPN
func (h *CallbackHandler) handleVPNAction(ctx context.Context, chatID int64, messageID int64, userID int64, action string) error {
	switch action {
	case "vpn_create":
		return h.createVPNConnection(ctx, chatID, messageID, userID)
	case "vpn_refresh":
		return h.refreshVPNConnection(ctx, chatID, messageID, userID)
	case "vpn_delete":
		return h.deleteVPNConnection(ctx, chatID, messageID, userID)
	default:
		return h.showVPN(ctx, chatID, messageID, userID)
	}
}

// handleSubscriptionAction обрабатывает действия с подпиской
func (h *CallbackHandler) handleSubscriptionAction(ctx context.Context, chatID int64, messageID int64, userID int64, action string) error {
	switch action {
	case "subscription_extend":
		return h.showSubscriptionPlans(ctx, chatID, messageID)
	case "subscription_1_month":
		return h.createPayment(ctx, chatID, messageID, userID, 1, 100)
	case "subscription_3_months":
		return h.createPayment(ctx, chatID, messageID, userID, 3, 250)
	case "subscription_6_months":
		return h.createPayment(ctx, chatID, messageID, userID, 6, 450)
	case "subscription_1_year":
		return h.createPayment(ctx, chatID, messageID, userID, 12, 800)
	default:
		return h.showSubscription(ctx, chatID, messageID, userID)
	}
}

// handleReferralAction обрабатывает действия с рефералами
func (h *CallbackHandler) handleReferralAction(ctx context.Context, chatID int64, messageID int64, userID int64, action string) error {
	switch action {
	case "referral_share":
		return h.shareReferralLink(ctx, chatID, messageID, userID)
	case "referral_stats":
		return h.showReferralStats(ctx, chatID, messageID, userID)
	default:
		return h.showReferral(ctx, chatID, messageID, userID)
	}
}

// handleSettingsAction обрабатывает действия с настройками
func (h *CallbackHandler) handleSettingsAction(ctx context.Context, chatID int64, messageID int64, userID int64, action string) error {
	switch action {
	case "settings_language":
		return h.showLanguageSettings(ctx, chatID, messageID)
	case "settings_notifications":
		return h.showNotificationSettings(ctx, chatID, messageID)
	case "settings_stats":
		return h.showUserStats(ctx, chatID, messageID, userID)
	case "settings_support":
		return h.showSupport(ctx, chatID, messageID)
	case "settings_about":
		return h.showAbout(ctx, chatID, messageID)
	case "settings_export":
		return h.exportUserData(ctx, chatID, messageID, userID)
	case "settings_delete":
		return h.showDeleteAccount(ctx, chatID, messageID, userID)
	default:
		return h.showSettings(ctx, chatID, int(messageID))
	}
}

// Command возвращает команду обработчика
func (h *CallbackHandler) Command() string {
	return "callback"
}

// Description возвращает описание обработчика
func (h *CallbackHandler) Description() string {
	return "Обработка callback query"
}

// Заглушки для недостающих методов
func (h *CallbackHandler) createVPNConnection(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showVPN(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) refreshVPNConnection(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showVPN(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) deleteVPNConnection(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showVPN(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) showSubscriptionPlans(ctx context.Context, chatID int64, messageID int64) error {
	return h.showSubscription(ctx, chatID, messageID, 0)
}

func (h *CallbackHandler) createPayment(ctx context.Context, chatID int64, messageID int64, userID int64, months int, amount int) error {
	return h.showSubscription(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) shareReferralLink(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showReferral(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) showReferralStats(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showReferral(ctx, chatID, messageID, userID)
}

func (h *CallbackHandler) showLanguageSettings(ctx context.Context, chatID int64, messageID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) showNotificationSettings(ctx context.Context, chatID int64, messageID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) showUserStats(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) showSupport(ctx context.Context, chatID int64, messageID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) showAbout(ctx context.Context, chatID int64, messageID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) exportUserData(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}

func (h *CallbackHandler) showDeleteAccount(ctx context.Context, chatID int64, messageID int64, userID int64) error {
	return h.showSettings(ctx, chatID, int(messageID))
}
