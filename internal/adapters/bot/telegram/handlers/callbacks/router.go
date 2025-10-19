package callbacks

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackFunc represents a callback handler function
type CallbackFunc func(ctx context.Context, userID, chatID int64, messageID int) error

// ParameterizedCallbackFunc represents a callback handler function with parameters
type ParameterizedCallbackFunc func(ctx context.Context, userID, chatID int64, messageID int, params ...string) error

// Router handles routing of callback queries to appropriate handlers
type Router struct {
	baseHandler *BaseHandler
	routes      map[string]CallbackFunc
}

// NewRouter creates a new Router with all dependencies
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

// setupRoutes initializes all route mappings
func (r *Router) setupRoutes() {
	// Trial routes
	r.routes["get_trial"] = r.baseHandler.HandleGetTrial

	// Menu routes
	r.routes["open_menu"] = r.baseHandler.HandleOpenMenu
	r.routes["open_profile"] = r.baseHandler.HandleOpenProfile
	r.routes["open_pricing"] = r.baseHandler.HandleOpenPricing
	r.routes["open_support"] = r.baseHandler.HandleOpenSupport

	// Subscription routes
	r.routes["my_subscriptions"] = r.baseHandler.HandleMySubscriptions
	r.routes["create_subscription"] = r.baseHandler.HandleCreateSubscription

	// VPN routes
	r.routes["open_keys"] = r.baseHandler.HandleOpenKeys
	r.routes["my_configs"] = r.baseHandler.HandleMyConfigs
	r.routes["create_wireguard"] = r.baseHandler.HandleCreateWireguard
	r.routes["create_shadowsocks"] = r.baseHandler.HandleCreateShadowsocks

	// Referral routes
	r.routes["open_referrals"] = r.baseHandler.HandleOpenReferrals
	r.routes["referral_stats"] = r.baseHandler.HandleReferralStats
	r.routes["my_referrals"] = r.baseHandler.HandleMyReferrals
	r.routes["my_referral_link"] = r.baseHandler.HandleMyReferralLink
}

// Handle processes a callback query and routes it to the appropriate handler
func (r *Router) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	userID := getUserID(update)
	chatID := getChatID(update)
	messageID := getMessageID(update)
	callbackData := update.CallbackQuery.Data

	slog.Info("Handling callback", "data", callbackData, "user_id", userID)

	// Try exact match first
	if handler, exists := r.routes[callbackData]; exists {
		return handler(ctx, userID, chatID, messageID)
	}

	// Handle parameterized callbacks
	return r.handleParameterizedCallback(ctx, userID, chatID, messageID, callbackData)
}

// handleParameterizedCallback handles callbacks with parameters
func (r *Router) handleParameterizedCallback(ctx context.Context, userID, chatID int64, messageID int, callbackData string) error {
	// Select plan callbacks
	if planID, ok := ui.ParseSelectPlanCallback(callbackData); ok {
		return r.baseHandler.HandleSelectPlan(ctx, userID, chatID, messageID, planID)
	}

	// Payment callbacks
	if planID, ok := ui.ParsePayCardCallback(callbackData); ok {
		return r.baseHandler.HandlePayCard(ctx, userID, chatID, messageID, planID)
	}
	if planID, ok := ui.ParsePaySBPCallback(callbackData); ok {
		return r.baseHandler.HandlePaySBP(ctx, userID, chatID, messageID, planID)
	}
	if planID, ok := ui.ParsePayStarsCallback(callbackData); ok {
		return r.baseHandler.HandlePayStars(ctx, userID, chatID, messageID, planID)
	}

	// Create subscription by plan
	if planID, ok := ui.ParseCreatePlanCallback(callbackData); ok {
		return r.baseHandler.HandleCreateSubscriptionByPlan(ctx, userID, chatID, messageID, planID)
	}

	// View subscription
	if subscriptionID, ok := ui.ParseViewSubscriptionCallback(callbackData); ok {
		return r.baseHandler.HandleViewSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Rename subscription
	if subscriptionID, ok := ui.ParseRenameSubscriptionCallback(callbackData); ok {
		return r.baseHandler.HandleRenameSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Extend subscription
	if subscriptionID, ok := ui.ParseExtendSubscriptionCallback(callbackData); ok {
		return r.baseHandler.HandleExtendSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Delete subscription
	if subscriptionID, ok := ui.ParseDeleteSubscriptionCallback(callbackData); ok {
		return r.baseHandler.HandleDeleteSubscription(ctx, userID, chatID, messageID, subscriptionID)
	}

	// Extend subscription by plan
	if planID, subscriptionID, ok := ui.ParseExtendPlanCallback(callbackData); ok {
		return r.baseHandler.HandleExtendSubscriptionByPlan(ctx, userID, chatID, messageID, planID, subscriptionID)
	}

	// View config
	if configID, ok := ui.ParseViewConfigCallback(callbackData); ok {
		return r.handleViewConfig(ctx, userID, chatID, messageID, configID)
	}

	// Connection guide
	if configID, ok := ui.ParseConnectionGuideCallback(callbackData); ok {
		return r.handleConnectionGuide(ctx, userID, chatID, messageID, configID)
	}

	// Unknown callback
	text := ui.GetUnknownCommandText()
	keyboard := ui.GetUnknownCommandKeyboard()
	return r.baseHandler.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}

// HandleTextMessage handles text messages (for renaming subscriptions)
func (r *Router) HandleTextMessage(ctx context.Context, userID int64, chatID int64, messageText string) (bool, error) {
	// Check if user is in renaming state
	r.baseHandler.mu.RLock()
	subscriptionID, isRenaming := r.baseHandler.renamingUsers[userID]
	r.baseHandler.mu.RUnlock()

	if !isRenaming {
		return false, nil // Not handling, user is not in renaming process
	}

	// Clear renaming state
	r.baseHandler.mu.Lock()
	delete(r.baseHandler.renamingUsers, userID)
	r.baseHandler.mu.Unlock()

	// Validate new name
	if len(messageText) < 1 || len(messageText) > 50 {
		text := "❌ Название подписки должно содержать от 1 до 50 символов"
		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	// Update subscription name
	err := r.baseHandler.updateSubscriptionName(ctx, userID, subscriptionID, messageText)
	if err != nil {
		r.baseHandler.logError(err, "UpdateSubscriptionName")
		text := "❌ Ошибка при переименовании подписки"
		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	// Get updated subscription
	subscription, err := r.baseHandler.getSubscription(ctx, userID, subscriptionID)
	if err != nil {
		r.baseHandler.logError(err, "GetSubscription")
		text := "✅ Подписка переименована, но произошла ошибка при получении обновленной информации"
		return true, r.baseHandler.msg.SendMessage(ctx, chatID, text)
	}

	// Show updated subscription information
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

// ============================================================================
// HELPER METHODS FOR MISSING HANDLERS
// ============================================================================

// handleViewConfig handles view_config callback (placeholder)
func (r *Router) handleViewConfig(ctx context.Context, userID, chatID int64, messageID int, configID string) error {
	slog.Info("Handling view config", "config_id", configID, "user_id", userID)
	text := "🔑 Просмотр конфигурации\n\nКонфигурация: " + configID
	return r.baseHandler.msg.EditMessageText(ctx, chatID, messageID, text, nil)
}

// handleConnectionGuide handles connection_guide callback (placeholder)
func (r *Router) handleConnectionGuide(ctx context.Context, userID, chatID int64, messageID int, configID string) error {
	slog.Info("Handling connection guide", "config_id", configID, "user_id", userID)
	text := "📖 Руководство по подключению\n\nКонфигурация: " + configID
	return r.baseHandler.msg.EditMessageText(ctx, chatID, messageID, text, nil)
}

// handleCreateSubscriptionByPlan handles create_subscription_by_plan callback
func (r *Router) handleCreateSubscriptionByPlan(ctx context.Context, userID, chatID int64, messageID int, planID string) error {
	slog.Info("Handling create subscription by plan", "plan_id", planID, "user_id", userID)

	// Get plan
	plan, err := r.baseHandler.getPlan(ctx, planID)
	if err != nil {
		r.baseHandler.logError(err, "GetPlan")
		return err
	}

	text := ui.GetPaymentMethodText(plan)
	keyboard := ui.GetPaymentMethodKeyboard(planID)
	// Use deleteAndSendMessage as it might be a photo
	return r.baseHandler.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}
