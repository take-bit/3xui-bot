package callback

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackFunc func(ctx context.Context, userID, chatID int64, messageID int) error

type ParameterizedCallbackFunc func(ctx context.Context, userID, chatID int64, messageID int, params ...string) error

type Router struct {
	baseHandler *BaseHandler
	routes      map[string]CallbackFunc
}

func NewRouter(
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
	msg *service.MessageService,
) *Router {
	baseHandler := NewBaseHandler(userUC, subUC, paymentUC, vpnUC, referralUC, notifUC, msg)

	router := &Router{
		baseHandler: baseHandler,
		routes:      make(map[string]CallbackFunc),
	}

	router.setupRoutes()

	return router
}

func (r *Router) setupRoutes() {
	r.routes["get_trial"] = r.baseHandler.HandleGetTrial

	r.routes["open_menu"] = r.baseHandler.HandleOpenMenu
	r.routes["open_profile"] = r.baseHandler.HandleOpenProfile
	r.routes["open_pricing"] = r.baseHandler.HandleOpenPricing
	r.routes["open_support"] = r.baseHandler.HandleOpenSupport
	r.routes["show_instruction"] = r.baseHandler.HandleShowInstruction

	r.routes["my_subscriptions"] = r.baseHandler.HandleMySubscriptions
	r.routes["create_subscription"] = r.baseHandler.HandleCreateSubscription

	r.routes["open_keys"] = r.baseHandler.HandleOpenKeys
	r.routes["my_configs"] = r.baseHandler.HandleMyConfigs
	r.routes["create_wireguard"] = r.baseHandler.HandleCreateWireguard
	r.routes["create_shadowsocks"] = r.baseHandler.HandleCreateShadowsocks

	r.routes["open_referrals"] = r.baseHandler.HandleOpenReferrals
	r.routes["referral_stats"] = r.baseHandler.HandleReferralStats
	r.routes["my_referrals"] = r.baseHandler.HandleMyReferrals
	r.routes["my_referral_link"] = r.baseHandler.HandleMyReferralLink
}

func (r *Router) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {

		return nil
	}

	userID := getUserID(update)
	chatID := getChatID(update)
	messageID := getMessageID(update)
	callbackData := update.CallbackQuery.Data

	slog.Info("Handling callback", "data", callbackData, "user_id", userID)

	if handler, exists := r.routes[callbackData]; exists {

		return handler(ctx, userID, chatID, messageID)
	}

	return r.handleParameterizedCallback(ctx, userID, chatID, messageID, callbackData)
}

func (r *Router) handleParameterizedCallback(ctx context.Context, userID, chatID int64, messageID int, callbackData string) error {
	if planID, ok := ui.ParseSelectPlanCallback(callbackData); ok {

		return r.baseHandler.HandleSelectPlan(ctx, userID, chatID, messageID, planID)
	}

	if planID, ok := ui.ParsePayCardCallback(callbackData); ok {

		return r.baseHandler.HandlePayCard(ctx, userID, chatID, messageID, planID)
	}
	if planID, ok := ui.ParsePaySBPCallback(callbackData); ok {

		return r.baseHandler.HandlePaySBP(ctx, userID, chatID, messageID, planID)
	}
	if planID, ok := ui.ParsePayStarsCallback(callbackData); ok {

		return r.baseHandler.HandlePayStars(ctx, userID, chatID, messageID, planID)
	}

	if planID, ok := ui.ParseCreatePlanCallback(callbackData); ok {

		return r.baseHandler.HandleCreateSubscriptionByPlan(ctx, userID, chatID, messageID, planID)
	}

	if subscriptionID, ok := ui.ParseViewSubscriptionCallback(callbackData); ok {

		return r.baseHandler.HandleViewSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	if subscriptionID, ok := ui.ParseRenameSubscriptionCallback(callbackData); ok {

		return r.baseHandler.HandleRenameSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	if subscriptionID, ok := ui.ParseExtendSubscriptionCallback(callbackData); ok {

		return r.baseHandler.HandleExtendSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	if subscriptionID, ok := ui.ParseDeleteSubscriptionCallback(callbackData); ok {

		return r.baseHandler.HandleDeleteSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	if planID, subscriptionID, ok := ui.ParseExtendPlanCallback(callbackData); ok {

		return r.baseHandler.HandleExtendSubscriptionByPlan(ctx, userID, chatID, messageID, planID, subscriptionID)
	}

	if configID, ok := ui.ParseViewConfigCallback(callbackData); ok {

		return r.handleViewConfig(ctx, userID, chatID, messageID, configID)
	}

	if configID, ok := ui.ParseConnectionGuideCallback(callbackData); ok {

		return r.handleConnectionGuide(ctx, userID, chatID, messageID, configID)
	}

	text := ui.GetUnknownCommandText()
	keyboard := ui.GetUnknownCommandKeyboard()

	return r.baseHandler.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (r *Router) HandleTextMessage(ctx context.Context, userID int64, chatID int64, messageText string) (bool, error) {
	r.baseHandler.mu.RLock()
	subscriptionID, isRenaming := r.baseHandler.renamingUsers[userID]
	r.baseHandler.mu.RUnlock()

	if !isRenaming {

		return false, nil
	}

	r.baseHandler.mu.Lock()
	delete(r.baseHandler.renamingUsers, userID)
	r.baseHandler.mu.Unlock()

	if len(messageText) < 1 || len(messageText) > 50 {
		text := "❌ Название подписки должно содержать от 1 до 50 символов"

		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	err := r.baseHandler.updateSubscriptionName(ctx, userID, subscriptionID, messageText)
	if err != nil {
		r.baseHandler.logError(err, "UpdateSubscriptionName")
		text := "❌ Ошибка при переименовании подписки"

		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	subscription, err := r.baseHandler.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		r.baseHandler.logError(err, "GetSubscription")
		text := "✅ Подписка переименована, но произошла ошибка при получении обновленной информации"

		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	plan, err := r.baseHandler.getPlan(ctx, subscription.PlanID)
	if err != nil {
		r.baseHandler.logError(err, "GetPlan")
		plan = &core.Plan{Name: "Неизвестный план"}
	}

	vpnConfigs, err := r.baseHandler.getVPNConnectionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		vpnConfigs = []*core.VPNConnection{}
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)

	return true, r.baseHandler.msg.SendMessageWithMarkdownV2(ctx, chatID, text, keyboard)
}

func (r *Router) handleViewConfig(ctx context.Context, userID, chatID int64, messageID int, configID string) error {
	slog.Info("Handling view config", "config_id", configID, "user_id", userID)
	text := "🔑 Просмотр конфигурации\n\nКонфигурация: " + configID

	return r.baseHandler.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, nil)
}

func (r *Router) handleConnectionGuide(ctx context.Context, userID, chatID int64, messageID int, configID string) error {
	slog.Info("Handling connection guide", "config_id", configID, "user_id", userID)
	text := "📖 Руководство по подключению\n\nКонфигурация: " + configID

	return r.baseHandler.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, nil)
}

func (r *Router) handleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling create subscription by plan", "plan_id", planID, "user_id", userID)

	plan, err := r.baseHandler.getPlan(ctx, planID)
	if err != nil {
		r.baseHandler.logError(err, "GetPlan")

		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)

	return r.baseHandler.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}
