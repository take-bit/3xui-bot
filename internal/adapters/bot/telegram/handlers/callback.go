package handlers

import (
	"context"
	"fmt"
	"log"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackHandler обрабатывает все callback query
type CallbackHandler struct {
	controller interface{}
}

// NewCallbackHandler создает новый обработчик callback'ов
func NewCallbackHandler(controller interface{}) *CallbackHandler {
	return &CallbackHandler{controller: controller}
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

	log.Printf("Handling callback: %s for user %d", callbackData, userID)

	// Отвечаем на callback query
	err := h.answerCallbackQuery(ctx, update.CallbackQuery.ID, "", false)
	if err != nil {
		log.Printf("Error answering callback query: %v", err)
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
	log.Printf("Handling get trial for user %d", userID)

	// Активируем пробный доступ
	success, err := h.activateTrial(ctx, userID)
	if err != nil {
		h.logError(err, "ActivateTrial")
		return err
	}

	var text string
	if success {
		text = "🎉 Пробный доступ активирован на 3 дня!"
	} else {
		text = "❌ Пробный доступ уже был использован"
	}

	keyboard := ui.GetWelcomeKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenMenu(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling open menu for user %d", userID)
	text := ui.GetWelcomeText()
	keyboard := ui.GetWelcomeKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenProfile(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling open profile for user %d", userID)

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
	log.Printf("Handling open pricing for user %d", userID)

	// Получаем планы
	plans, err := h.getPlans(ctx)
	if err != nil {
		h.logError(err, "GetPlans")
		return err
	}

	text := ui.GetPricingText(plans)
	keyboard := ui.GetPricingKeyboard(plans)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMySubscriptions(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling my subscriptions for user %d", userID)

	// Получаем подписки пользователя
	subscriptions, err := h.getUserSubscriptions(ctx, userID)
	if err != nil {
		h.logError(err, "GetUserSubscriptions")
		return err
	}

	text := ui.GetSubscriptionsText(subscriptions)
	keyboard := ui.GetSubscriptionsKeyboard(subscriptions)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleCreateSubscription(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling create subscription for user %d", userID)

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
	log.Printf("Handling open keys for user %d", userID)
	text := ui.GetKeysText()
	keyboard := ui.GetKeysKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling open referrals for user %d", userID)
	text := ui.GetReferralsText()
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenSupport(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling open support for user %d", userID)
	text := ui.GetSupportText()
	keyboard := ui.GetWelcomeKeyboard() // Используем базовую клавиатуру
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyConfigs(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling my configs for user %d", userID)
	text := "📋 Ваши VPN конфигурации\n\nПока конфигураций нет."
	keyboard := ui.GetKeysKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleReferralStats(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling referral stats for user %d", userID)
	text := "📊 Статистика рефералов\n\nПока статистики нет."
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling my referrals for user %d", userID)
	text := "👥 Ваши рефералы\n\nПока рефералов нет."
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleMyReferralLink(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling my referral link for user %d", userID)
	text := "🔗 Ваша реферальная ссылка\n\nhttps://t.me/your_bot?start=ref_123456"
	keyboard := ui.GetReferralsKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleCreateWireguard(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling create wireguard for user %d", userID)
	text := "🔑 Создание WireGuard конфигурации\n\nВведите название для конфигурации:"
	return h.sendMessage(ctx, chatID, text)
}

func (h *CallbackHandler) handleCreateShadowsocks(ctx context.Context, userID, chatID int64, messageID int) error {
	log.Printf("Handling create shadowsocks for user %d", userID)
	text := "🔑 Создание Shadowsocks конфигурации\n\nВведите название для конфигурации:"
	return h.sendMessage(ctx, chatID, text)
}

// ============================================================================
// ОБРАБОТЧИКИ ПАРАМЕТРИЗОВАННЫХ CALLBACK'ОВ
// ============================================================================

func (h *CallbackHandler) handleParameterizedCallback(ctx context.Context, userID, chatID int64, messageID int, callbackData string) error {
	log.Printf("Handling parameterized callback: %s for user %d", callbackData, userID)

	// Планы подписок
	if planID, ok := ui.ParsePlanCallback(callbackData); ok {
		return h.handlePlanSelection(ctx, userID, chatID, messageID, planID)
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
	log.Printf("Handling plan selection %s for user %d", planID, userID)

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

func (h *CallbackHandler) handleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	log.Printf("Handling create subscription by plan %s for user %d", planID, userID)

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

	text := fmt.Sprintf("✅ Подписка '%s' создана успешно!\n⏰ Длительность: %d дней", plan.Name, plan.Days)
	keyboard := ui.GetWelcomeKeyboard()
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleViewSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	log.Printf("Handling view subscription %s for user %d", subscriptionID, userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	text := ui.GetSubscriptionDetailText(subscription)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscriptionID)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleRenameSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	log.Printf("Handling rename subscription %s for user %d", subscriptionID, userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	text := ui.GetRenameSubscriptionText(subscription)
	return h.sendMessage(ctx, chatID, text)
}

func (h *CallbackHandler) handleExtendSubscription(ctx context.Context, userID, chatID int64, messageID int, subscriptionID string) error {
	log.Printf("Handling extend subscription %s for user %d", subscriptionID, userID)

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
	log.Printf("Handling delete subscription %s for user %d", subscriptionID, userID)

	// Получаем подписку
	subscription, err := h.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		h.logError(err, "GetSubscription")
		return err
	}

	text := ui.GetDeleteSubscriptionText(subscription)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscriptionID)
	return h.editMessageText(ctx, chatID, messageID, text, keyboard)
}

func (h *CallbackHandler) handleExtendSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID, subscriptionID string) error {
	log.Printf("Handling extend subscription %s by plan %s for user %d", subscriptionID, planID, userID)

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

	text := fmt.Sprintf("✅ Подписка продлена на %d дней!", plan.Days)
	keyboard := ui.GetWelcomeKeyboard()
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

func (h *CallbackHandler) editMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error {
	if bot, ok := h.controller.(interface {
		EditMessageText(ctx context.Context, chatID int64, messageID int, text string, replyMarkup interface{}) error
	}); ok {
		return bot.EditMessageText(ctx, chatID, messageID, text, replyMarkup)
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
