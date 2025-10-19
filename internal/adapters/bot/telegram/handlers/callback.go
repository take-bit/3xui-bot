package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackHandler обрабатывает все callback query
type CallbackHandler struct {
	controller interface{}
	// Состояние для переименования подписок
	renamingUsers map[int64]string // userID -> subscriptionID
	mu            sync.RWMutex
}

// NewCallbackHandler создает новый обработчик callback'ов
func NewCallbackHandler(controller interface{}) *CallbackHandler {
	return &CallbackHandler{
		controller:    controller,
		renamingUsers: make(map[int64]string),
	}
}

// CanHandle проверяет, может ли обработчик обработать обновление
func (h *CallbackHandler) CanHandle(update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

// Handle обрабатывает callback query
func (h *CallbackHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	userID := h.getUserID(update)
	chatID := h.getChatID(update)
	messageID := h.getMessageID(update)
	callbackData := update.CallbackQuery.Data

	slog.Info("Handling callback", "data", callbackData, "user_id", userID)

	// Отвечаем на callback query
	err := h.answerCallbackQuery(ctx, update.CallbackQuery.ID, "", false)
	if err != nil {
		slog.Error("Error answering callback query", "error", err)
	}

	// Обрабатываем callback
	switch callbackData {
	// Основные команды
	case "get_trial":
		return h.handleGetTrial(ctx, userID, chatID, messageID)
	case "open_menu":
		return h.handleOpenMenu(ctx, userID, chatID, messageID)
	case "open_profile":
		return h.handleOpenProfile(ctx, userID, chatID, messageID)
	case "open_pricing":
		return h.handleOpenPricing(ctx, userID, chatID, messageID)
	case "my_subscriptions":
		return h.handleMySubscriptions(ctx, userID, chatID, messageID)
	case "create_subscription":
		return h.handleCreateSubscription(ctx, userID, chatID, messageID)
	case "open_keys":
		return h.handleOpenKeys(ctx, userID, chatID, messageID)
	case "open_referrals":
		return h.handleOpenReferrals(ctx, userID, chatID, messageID)
	case "open_support":
		return h.handleOpenSupport(ctx, userID, chatID, messageID)
	case "my_configs":
		return h.handleMyConfigs(ctx, userID, chatID, messageID)
	case "referral_stats":
		return h.handleReferralStats(ctx, userID, chatID, messageID)
	case "my_referrals":
		return h.handleMyReferrals(ctx, userID, chatID, messageID)
	case "my_referral_link":
		return h.handleMyReferralLink(ctx, userID, chatID, messageID)
	case "create_wireguard":
		return h.handleCreateWireguard(ctx, userID, chatID, messageID)
	case "create_shadowsocks":
		return h.handleCreateShadowsocks(ctx, userID, chatID, messageID)
	default:
		return h.handleParameterizedCallback(ctx, userID, chatID, messageID, callbackData)
	}
}

// ============================================================================
// ОСНОВНЫЕ ОБРАБОТЧИКИ
// ============================================================================

func (h *CallbackHandler) handleGetTrial(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling get trial", "user_id", userID)

	// Получаем пользователя
	userObj, err := h.getUser(ctx, userID)
	if err != nil || userObj == nil {
		h.logError(err, "GetUser")
		return err
	}
	user := userObj.(*core.User)

	// Активируем пробный доступ
	success, err := h.activateTrial(ctx, userID)
	if err != nil {
		h.logError(err, "ActivateTrial")
		return err
	}

	var text string
	if success {
		text = "🎉 Пробный доступ активирован на 3 дня!"

		// Создаем пробную подписку
		err = h.createTrialSubscription(ctx, userID)
		if err != nil {
			h.logError(err, "CreateTrialSubscription")
			// Не возвращаем ошибку, т.к. пробный доступ уже активирован
		}

		// Обновляем user.HasTrial для правильного отображения клавиатуры
		user.HasTrial = true
	} else {
		text = "❌ Пробный доступ уже был использован"
	}

	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
	// Используем deleteAndSendMessage, т.к. приветственное сообщение может быть с фото
	return h.deleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenMenu(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open menu", "user_id", userID)

	// Получаем пользователя
	userObj, err := h.getUser(ctx, userID)
	if err != nil || userObj == nil {
		h.logError(err, "GetUser")
		return err
	}
	user := userObj.(*core.User)

	// Проверяем активные подписки
	subscriptions, _ := h.getUserSubscriptions(ctx, userID)
	isPremium := false
	statusText := "🆓 Бесплатный"
	subUntilText := ""

	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if !sub.IsExpired() {
				isPremium = true
				statusText = "⭐ Premium"
				subUntilText = sub.EndDate.Format("02.01.2006")
				break
			}
		}
	}

	text := ui.GetMainMenuWithProfileText(user, isPremium, statusText, subUntilText)
	keyboard := ui.GetMainMenuWithProfileKeyboard(isPremium)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenProfile(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open profile", "user_id", userID)

	// Получаем пользователя
	userObj, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}

	user := userObj.(*core.User)

	// Проверяем наличие активной подписки
	subsObj, _ := h.getUserSubscriptions(ctx, userID)
	isPremium := false
	statusText := "Free"
	subUntilText := "—"

	if len(subsObj) > 0 {
		for _, sub := range subsObj {
			if !sub.IsExpired() {
				isPremium = true
				statusText = "Premium"
				subUntilText = sub.EndDate.Format("02.01.2006")
				break
			}
		}
	}

	text := ui.GetProfileText(user, isPremium, statusText, subUntilText)
	keyboard := ui.GetProfileKeyboard(isPremium)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenPricing(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open pricing", "user_id", userID)

	// Получаем планы
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)
	// Используем deleteAndSendMessage, т.к. главное меню может быть с фото
	return h.deleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMySubscriptions(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my subscriptions", "user_id", userID)

	// Получаем подписки пользователя
	subscriptions, err := h.getUserSubscriptions(ctx, userID)
	if err != nil {
		h.logError(err, "GetUserSubscriptions")
		return err
	}

	text := ui.GetSubscriptionsText(subscriptions)
	keyboard := ui.GetSubscriptionsKeyboard(subscriptions)
	// Используем deleteAndSendMessageWithMarkdownV2 для корректного форматирования
	return h.deleteAndSendMessageWithMarkdownV2(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleCreateSubscription(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create subscription", "user_id", userID)

	// Получаем планы
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetCreateSubscriptionText()
	keyboard := ui.GetCreateSubscriptionKeyboard(plans)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenKeys(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open keys", "user_id", userID)
	text := ui.GetKeysText()
	keyboard := ui.GetKeysKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open referrals", "user_id", userID)
	text := ui.GetReferralsText()
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenSupport(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open support", "user_id", userID)

	// Получаем пользователя
	userObj, err := h.getUser(ctx, userID)
	if err != nil || userObj == nil {
		h.logError(err, "GetUser")
		return err
	}
	user := userObj.(*core.User)

	text := ui.GetSupportText()
	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyConfigs(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my configs", "user_id", userID)
	text := "📋 Ваши VPN конфигурации\n\nПока конфигураций нет."
	keyboard := ui.GetKeysKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleReferralStats(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling referral stats", "user_id", userID)
	text := "📊 Статистика рефералов\n\nПока статистики нет."
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my referrals", "user_id", userID)
	text := "👥 Ваши рефералы\n\nПока рефералов нет."
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyReferralLink(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my referral link", "user_id", userID)
	text := "🔗 Ваша реферальная ссылка\n\nhttps://t.me/your_bot?start=ref_123456"
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleCreateWireguard(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create wireguard", "user_id", userID)
	text := "🔑 Создание WireGuard конфигурации\n\nВведите название для конфигурации:"
	return h.sendMessage(ctx, chatID, text)
}

func (h *CallbackHandler) handleCreateShadowsocks(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create shadowsocks", "user_id", userID)
	text := "🔑 Создание Shadowsocks конфигурации\n\nВведите название для конфигурации:"
	return h.sendMessage(ctx, chatID, text)
}

// ============================================================================
// ОБРАБОТЧИКИ ПАРАМЕТРИЗОВАННЫХ CALLBACK'ОВ
// ============================================================================

func (h *CallbackHandler) handleParameterizedCallback(ctx context.Context, userID, chatID int64, messageID int, callbackData string) error {
	slog.Info("Handling parameterized callback", "data", callbackData, "user_id", userID)

	// Планы подписок
	if planID, ok := ui.ParsePlanCallback(callbackData); ok {
		return h.handlePlanSelection(ctx, userID, chatID, messageID, planID)
	}

	// Выбор плана (select_plan_)
	if planID, ok := ui.ParseSelectPlanCallback(callbackData); ok {
		return h.handleSelectPlan(ctx, userID, chatID, messageID, planID)
	}

	// Оплата картой
	if planID, ok := ui.ParsePayCardCallback(callbackData); ok {
		return h.handlePayCard(ctx, userID, chatID, messageID, planID)
	}

	// Оплата СБП
	if planID, ok := ui.ParsePaySBPCallback(callbackData); ok {
		return h.handlePaySBP(ctx, userID, chatID, messageID, planID)
	}

	// Оплата Stars
	if planID, ok := ui.ParsePayStarsCallback(callbackData); ok {
		return h.handlePayStars(ctx, userID, chatID, messageID, planID)
	}

	// Создание подписки по плану
	if planID, ok := ui.ParseCreatePlanCallback(callbackData); ok {
		return h.handleCreateSubscriptionByPlan(ctx, userID, chatID, messageID, planID)
	}

	// Просмотр подписки
	if subscriptionID, ok := ui.ParseViewSubscriptionCallback(callbackData); ok {
		return h.handleViewSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Переименование подписки
	if subscriptionID, ok := ui.ParseRenameSubscriptionCallback(callbackData); ok {
		return h.handleRenameSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Продление подписки
	if subscriptionID, ok := ui.ParseExtendSubscriptionCallback(callbackData); ok {
		return h.handleExtendSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Удаление подписки
	if subscriptionID, ok := ui.ParseDeleteSubscriptionCallback(callbackData); ok {
		return h.handleDeleteSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Продление подписки по плану
	if planID, subscriptionID, ok := ui.ParseExtendPlanCallback(callbackData); ok {
		return h.handleExtendSubscriptionByPlan(ctx, userID, chatID, messageID, planID, subscriptionID)
	}

	text := "❓ Неизвестная команда"
	return h.editMessageText(ctx, chatID, messageID, text, nil)
}

func (h *CallbackHandler) handlePlanSelection(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling plan selection", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := fmt.Sprintf("📦 План: %s\n💵 Цена: %.0f₽\n⏰ Длительность: %d дней\n\nСоздать подписку?", plan.Name, plan.Price, plan.Days)
	keyboard := ui.GetPricingKeyboard([]*core.Plan{plan})
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleSelectPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling select plan", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := fmt.Sprintf("📦 План: %s\n💵 Цена: %.0f₽\n⏰ Длительность: %d дней\n\nВыберите способ оплаты:", plan.Name, plan.Price, plan.Days)
	keyboard := ui.GetPaymentMethodKeyboard(planID)
	// Используем deleteAndSendMessage, т.к. может быть фото
	return h.deleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handlePayCard(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay card", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// MOCK: Картой всегда успешно
	slog.Info("Creating MOCK card payment (auto-success)", "user_id", userID, "plan_id", planID)

	now := time.Now()

	// Генерируем осмысленное название подписки на основе плана
	subscriptionName := h.generateSubscriptionName(plan)

	dto := usecase.CreateSubscriptionDTO{
		UserID:    userID,
		Name:      subscriptionName,
		PlanID:    planID,
		Days:      plan.Days,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.Days),
		IsActive:  true,
	}

	subscription, err := h.createSubscription(ctx, dto)
	if err != nil {
		h.logError(err, "CreateSubscription")
		return h.sendError(chatID, "Ошибка создания подписки")
	}

	// Создаем VPN для подписки
	sub, ok := subscription.(*core.Subscription)
	if ok && sub != nil {
		if vpnUC, ok := h.controller.(interface {
			VpnUC() *usecase.VPNUseCase
		}); ok {
			_, err = vpnUC.VpnUC().CreateVPNForSubscription(ctx, userID, sub.ID)
			if err != nil {
				h.logError(err, "CreateVPN")
			}
		}
	}

	text := fmt.Sprintf("✅ Оплата успешна!\n\n🎉 Подписка '%s' активирована на %d дней", plan.Name, plan.Days)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В главное меню", "open_menu"),
		),
	)
	return h.deleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handlePaySBP(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay SBP", "plan_id", planID, "user_id", userID)
	// MOCK: СБП работает аналогично карте
	return h.handlePayCard(ctx, userID, chatID, messageID, planID)
}

func (h *CallbackHandler) handlePayStars(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling pay stars", "plan_id", planID, "user_id", userID)

	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	text := fmt.Sprintf("⭐ Оплата Telegram Stars\n\nПлан: %s\nЦена: %.0f₽\n\nФункция в разработке", plan.Name, plan.Price)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к тарифам", "open_pricing"),
		),
	)
	return h.deleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) sendError(chatID int64, text string) error {
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		msg := tgbotapi.NewMessage(chatID, "❌ "+text)
		_, err := bot.Bot().Send(msg)
		return err
	}
	return nil
}

func (h *CallbackHandler) handleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling create subscription by plan", "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// Создаем подписку
	createSubDTO := usecase.CreateSubscriptionDTO{
		UserID: userID,
		Name:   "Основная",
		PlanID: planID,
		Days:   plan.Days,
	}

	_, err = h.createSubscription(ctx, createSubDTO)
	if err != nil {
		h.logError(err, "CreateSubscription")
		return err
	}

	// Получаем пользователя для клавиатуры
	userObj, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}
	user := userObj.(*core.User)

	text := fmt.Sprintf("✅ Подписка '%s' создана успешно!\n⏰ Длительность: %d дней", plan.Name, plan.Days)
	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleViewSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling view subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Получаем план подписки
	plan, err := h.getPlan(ctx, subscription.PlanID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// Получаем VPN конфигурации подписки
	vpnConfigs, err := h.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		h.logError(err, "GetVPNConnections")
		vpnConfigs = []*core.VPNConnection{} // Пустой массив, если ошибка
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	return h.editMessageWithMarkdownV2(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleRenameSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling rename subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Сохраняем состояние переименования
	h.mu.Lock()
	h.renamingUsers[userID] = subscriptionID
	h.mu.Unlock()

	text := ui.GetRenameSubscriptionText(subscription)
	return h.sendMessage(ctx, chatID, text)
}

func (h *CallbackHandler) handleExtendSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling extend subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Получаем планы для продления
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetExtendSubscriptionText(subscription)
	keyboard := ui.GetExtendSubscriptionKeyboard(subscriptionID, plans)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleDeleteSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	slog.Info("Handling delete subscription", "subscription_id", subscriptionID, "user_id", userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	// Получаем VPN конфигурации подписки
	vpnConfigs, err := h.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		h.logError(err, "GetVPNConnections")
		vpnConfigs = []*core.VPNConnection{} // Пустой массив, если ошибка
	}

	text := ui.GetDeleteSubscriptionText(subscription)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleExtendSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID, subscriptionID string) error {
	slog.Info("Handling extend subscription by plan", "subscription_id", subscriptionID, "plan_id", planID, "user_id", userID)

	// Получаем план
	plan, err := h.getPlan(ctx, planID)
	if err != nil {
		h.logError(err, "GetPlan")
		return err
	}

	// Продлеваем подписку
	err = h.extendSubscription(ctx, userID, subscriptionID, plan.Days)
	if err != nil {
		h.logError(err, "ExtendSubscription")
		return err
	}

	// Получаем пользователя для клавиатуры
	userObj, err := h.getUser(ctx, userID)
	if err != nil {
		h.logError(err, "GetUser")
		return err
	}
	user := userObj.(*core.User)

	text := fmt.Sprintf("✅ Подписка продлена на %d дней!", plan.Days)
	keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
// ============================================================================

func (h *CallbackHandler) getUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

func (h *CallbackHandler) getChatID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return 0
}

func (h *CallbackHandler) getMessageID(update tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.MessageID
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.MessageID
	}
	return 0
}

func (h *CallbackHandler) sendMessage(ctx context.Context, chatID int64, text string) error {
	if bot, ok := h.controller.(interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
	}); ok {
		return bot.SendMessage(ctx, chatID, text)
	}
	return nil
}

func (h *CallbackHandler) sendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard interface{}) error {
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "MarkdownV2"
		if keyboard != nil {
			if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
				msg.ReplyMarkup = kb
			}
		}
		_, err := bot.Bot().Send(msg)
		return err
	}
	return nil
}

func (h *CallbackHandler) editMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	if bot, ok := h.controller.(interface {
		EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error
	}); ok {
		return bot.EditMessageText(ctx, chatID, messageID, text, replyMarkup)
	}
	return nil
}

// editMessageWithMarkdownV2 редактирует сообщение с MarkdownV2 форматированием
func (h *CallbackHandler) editMessageWithMarkdownV2(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
		editMsg.ParseMode = "MarkdownV2"

		if replyMarkup != nil {
			if kb, ok := replyMarkup.(tgbotapi.InlineKeyboardMarkup); ok {
				editMsg.ReplyMarkup = &kb
			}
		}

		_, err := bot.Bot().Send(editMsg)
		return err
	}
	return nil
}

func (h *CallbackHandler) answerCallbackQuery(ctx context.Context, callbackQueryID, text string, showAlert bool) error {
	if bot, ok := h.controller.(interface {
		AnswerCallbackQuery(ctx context.Context, callbackQueryID, text string, showAlert bool) error
	}); ok {
		return bot.AnswerCallbackQuery(ctx, callbackQueryID, text, showAlert)
	}
	return nil
}

func (h *CallbackHandler) deleteMessage(ctx context.Context, chatID int64, messageID int) error {
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		_, err := bot.Bot().Request(deleteMsg)
		return err
	}
	return nil
}

func (h *CallbackHandler) deleteAndSendMessage(ctx context.Context, chatID int64, messageID int, text string, keyboard interface{}) error {
	// Удаляем старое сообщение (игнорируем ошибку)
	_ = h.deleteMessage(ctx, chatID, messageID)

	// Отправляем новое сообщение
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		msg := tgbotapi.NewMessage(chatID, text)
		if keyboard != nil {
			if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
				msg.ReplyMarkup = kb
			}
		}
		_, err := bot.Bot().Send(msg)
		return err
	}
	return nil
}

// deleteAndSendMessageWithMarkdownV2 удаляет старое сообщение и отправляет новое с MarkdownV2
func (h *CallbackHandler) deleteAndSendMessageWithMarkdownV2(ctx context.Context, chatID int64, messageID int, text string, keyboard interface{}) error {
	// Удаляем старое сообщение (игнорируем ошибку)
	_ = h.deleteMessage(ctx, chatID, messageID)

	// Отправляем новое сообщение с MarkdownV2
	if bot, ok := h.controller.(interface {
		Bot() *tgbotapi.BotAPI
	}); ok {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "MarkdownV2"
		if keyboard != nil {
			if kb, ok := keyboard.(tgbotapi.InlineKeyboardMarkup); ok {
				msg.ReplyMarkup = kb
			}
		}
		_, err := bot.Bot().Send(msg)
		return err
	}
	return nil
}

// ============================================================================
// USE CASE МЕТОДЫ
// ============================================================================

func (h *CallbackHandler) activateTrial(ctx context.Context, userID int64) (bool, error) {
	if userUC, ok := h.controller.(interface {
		UserUC() *usecase.UserUseCase
	}); ok {
		return userUC.UserUC().ActivateTrial(ctx, userID)
	}
	return false, nil
}

func (h *CallbackHandler) getUser(ctx context.Context, userID int64) (interface{}, error) {
	if userUC, ok := h.controller.(interface {
		UserUC() *usecase.UserUseCase
	}); ok {
		return userUC.UserUC().GetUser(ctx, userID)
	}
	return nil, nil
}

func (h *CallbackHandler) getUserSubscriptions(ctx context.Context, userID int64) ([]*core.Subscription, error) {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().GetUserSubscriptions(ctx, userID)
	}
	return nil, nil
}

func (h *CallbackHandler) getVPNConnectionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]*core.VPNConnection, error) {
	if vpnUC, ok := h.controller.(interface {
		VpnUC() *usecase.VPNUseCase
	}); ok {
		return vpnUC.VpnUC().GetVPNConnectionsBySubscription(ctx, subscriptionID)
	}
	return nil, nil
}

func (h *CallbackHandler) getPlans(ctx context.Context) ([]*core.Plan, error) {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().GetPlans(ctx)
	}
	return nil, nil
}

func (h *CallbackHandler) getPlan(ctx context.Context, planID string) (*core.Plan, error) {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().GetPlan(ctx, planID)
	}
	return nil, nil
}

func (h *CallbackHandler) createSubscription(ctx context.Context, dto usecase.CreateSubscriptionDTO) (interface{}, error) {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().CreateSubscription(ctx, dto)
	}
	return nil, nil
}

func (h *CallbackHandler) getSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.Subscription, error) {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().GetSubscription(ctx, userID, subscriptionID)
	}
	return nil, nil
}

func (h *CallbackHandler) extendSubscription(ctx context.Context, userID int64, subscriptionID string, days int) error {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().ExtendSubscription(ctx, userID, subscriptionID, days)
	}
	return nil
}

func (h *CallbackHandler) logError(err error, context string) {
	if logger, ok := h.controller.(interface {
		LogError(err error, context string)
	}); ok {
		logger.LogError(err, context)
	}
}

// generateSubscriptionName генерирует осмысленное название подписки на основе плана
func (h *CallbackHandler) generateSubscriptionName(plan *core.Plan) string {
	// Базовое название на основе длительности
	var baseName string
	switch {
	case plan.Days <= 7:
		baseName = "Недельная"
	case plan.Days <= 30:
		baseName = "Месячная"
	case plan.Days <= 90:
		baseName = "Квартальная"
	case plan.Days <= 365:
		baseName = "Годовая"
	default:
		baseName = "Долгосрочная"
	}

	// Добавляем дату создания для уникальности
	now := time.Now()
	dateStr := now.Format("02.01")

	return fmt.Sprintf("%s (%s)", baseName, dateStr)
}

// createTrialSubscription создает пробную подписку на 3 дня
func (h *CallbackHandler) createTrialSubscription(ctx context.Context, userID int64) error {
	now := time.Now()
	dto := usecase.CreateSubscriptionDTO{
		UserID:    userID,
		Name:      "Пробная подписка",
		PlanID:    "trial", // Специальный ID для пробной подписки
		Days:      3,       // 3 дня пробного доступа
		StartDate: now,
		EndDate:   now.AddDate(0, 0, 3),
		IsActive:  true,
	}

	subscription, err := h.createSubscription(ctx, dto)
	if err != nil {
		return err
	}

	// Создаем VPN для пробной подписки
	sub, ok := subscription.(*core.Subscription)
	if ok && sub != nil {
		if vpnUC, ok := h.controller.(interface {
			VpnUC() *usecase.VPNUseCase
		}); ok {
			_, err = vpnUC.VpnUC().CreateVPNForSubscription(ctx, userID, sub.ID)
			if err != nil {
				// Логируем ошибку, но не возвращаем её
				slog.Error("Failed to create VPN for trial subscription", "error", err)
			}
		}
	}

	return nil
}

// HandleTextMessage обрабатывает текстовые сообщения (для переименования подписок)
func (h *CallbackHandler) HandleTextMessage(ctx context.Context, userID int64, chatID int64, messageText string) (bool, error) {
	h.mu.RLock()
	subscriptionID, isRenaming := h.renamingUsers[userID]
	h.mu.RUnlock()

	if !isRenaming {
		return false, nil // Не обрабатываем, если пользователь не в процессе переименования
	}

	// Очищаем состояние переименования
	h.mu.Lock()
	delete(h.renamingUsers, userID)
	h.mu.Unlock()

	// Валидация нового названия
	if len(messageText) < 1 || len(messageText) > 50 {
		text := "❌ Название подписки должно содержать от 1 до 50 символов"
		return true, h.sendMessage(ctx, chatID, text)
	}

	// Обновляем название подписки
	err := h.updateSubscriptionName(ctx, userID, subscriptionID, messageText)
	if err != nil {
		h.logError(err, "UpdateSubscriptionName")
		text := "❌ Ошибка при переименовании подписки"
		return true, h.sendMessage(ctx, chatID, text)
	}

	// Получаем обновленную подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		text := "✅ Подписка переименована, но произошла ошибка при получении обновленной информации"
		return true, h.sendMessage(ctx, chatID, text)
	}

	// Показываем обновленную информацию о подписке
	plan, err := h.getPlan(ctx, subscription.PlanID)
	if err != nil {
		h.logError(err, "GetPlan")
		plan = &core.Plan{Name: "Неизвестный план"}
	}

	vpnConfigs, err := h.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		vpnConfigs = []*core.VPNConnection{}
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	
	return true, h.sendMessageWithKeyboard(ctx, chatID, text, keyboard)
}

// updateSubscriptionName обновляет название подписки
func (h *CallbackHandler) updateSubscriptionName(ctx context.Context, userID int64, subscriptionID, name string) error {
	if subUC, ok := h.controller.(interface {
		SubUC() *usecase.SubscriptionUseCase
	}); ok {
		return subUC.SubUC().UpdateSubscriptionName(ctx, userID, subscriptionID, name)
	}
	return fmt.Errorf("subscription use case not available")
}
