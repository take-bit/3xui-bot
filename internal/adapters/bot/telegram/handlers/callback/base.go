package callback

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/core"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BaseHandler contains shared dependencies and methods for all callback handlers
type BaseHandler struct {
	userUC     *usecase.UserUseCase
	subUC      *usecase.SubscriptionUseCase
	paymentUC  *usecase.PaymentUseCase
	vpnUC      *usecase.VPNUseCase
	referralUC *usecase.ReferralUseCase
	notifUC    *usecase.NotificationUseCase
	msg        *service.MessageService
	// Состояние для переименования подписок
	renamingUsers map[int64]string // userID -> subscriptionID
	mu            sync.RWMutex
}

// NewBaseHandler creates a new BaseHandler with all dependencies
func NewBaseHandler(
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
	msg *service.MessageService,
) *BaseHandler {
	return &BaseHandler{
		userUC:        userUC,
		subUC:         subUC,
		paymentUC:     paymentUC,
		vpnUC:         vpnUC,
		referralUC:    referralUC,
		notifUC:       notifUC,
		msg:           msg,
		renamingUsers: make(map[int64]string),
	}
}

// ============================================================================
// SHARED HELPER METHODS
// ============================================================================

// getUser retrieves a user by ID
func (h *BaseHandler) getUser(ctx context.Context, userID int64) (*core.User, error) {
	user, err := h.userUC.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// getSubscription retrieves a subscription by ID for a user
func (h *BaseHandler) getSubscription(ctx context.Context, userID int64, subscriptionID string) (*core.Subscription, error) {
	subscription, err := h.subUC.GetSubscription(ctx, userID, subscriptionID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, fmt.Errorf("subscription not found")
	}
	return subscription, nil
}

// getUserSubscriptions retrieves all subscriptions for a user
func (h *BaseHandler) getUserSubscriptions(ctx context.Context, userID int64) ([]*core.Subscription, error) {
	subscriptions, err := h.subUC.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*core.Subscription
	for _, sub := range subscriptions {
		result = append(result, sub)
	}
	return result, nil
}

// getPlan retrieves a plan by ID
func (h *BaseHandler) getPlan(ctx context.Context, planID string) (*core.Plan, error) {
	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, fmt.Errorf("plan not found")
	}
	return plan, nil
}

// getPlans retrieves all available plans
func (h *BaseHandler) getPlans(ctx context.Context) ([]*core.Plan, error) {
	plans, err := h.subUC.GetPlans(ctx)
	if err != nil {
		return nil, err
	}

	var result []*core.Plan
	for _, plan := range plans {
		result = append(result, plan)
	}
	return result, nil
}

// getVPNConnectionsBySubscriptionID retrieves VPN connections for a subscription
func (h *BaseHandler) getVPNConnectionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]*core.VPNConnection, error) {
	// This method doesn't exist in VPNUseCase, so we'll return empty for now
	// TODO: Implement this method in VPNUseCase if needed
	return []*core.VPNConnection{}, nil
}

// createSubscription creates a new subscription
func (h *BaseHandler) createSubscription(ctx context.Context, dto usecase.CreateSubscriptionDTO) (interface{}, error) {
	return h.subUC.CreateSubscription(ctx, dto)
}

// createVPNForSubscription creates a VPN connection for a subscription
func (h *BaseHandler) createVPNForSubscription(ctx context.Context, userID int64, subscriptionID string) (interface{}, error) {
	return h.vpnUC.CreateVPNForSubscription(ctx, userID, subscriptionID)
}

// updateSubscriptionName updates a subscription's name
func (h *BaseHandler) updateSubscriptionName(ctx context.Context, userID int64, subscriptionID, name string) error {
	return h.subUC.UpdateSubscriptionName(ctx, userID, subscriptionID, name)
}

// activateTrial activates trial access for a user
func (h *BaseHandler) activateTrial(ctx context.Context, userID int64) (bool, error) {
	return h.userUC.ActivateTrial(ctx, userID)
}

// logError logs an error with context
func (h *BaseHandler) logError(err error, context string) {
	if err != nil {
		slog.Error("Error", "context", context, "error", err)
	}
}

// sendError sends an error message to the user
func (h *BaseHandler) sendError(chatID int64, message string) error {
	return h.msg.SendMessage(context.Background(), chatID, message)
}

// ============================================================================
// UPDATE EXTRACTION HELPERS
// ============================================================================

// getUserID extracts user ID from update
func getUserID(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

// getChatID extracts chat ID from update
func getChatID(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return 0
}

// getMessageID extracts message ID from update
func getMessageID(update tgbotapi.Update) int {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.MessageID
	}
	return 0
}
