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

type BaseHandler struct {
	userUC        *usecase.UserUseCase
	subUC         *usecase.SubscriptionUseCase
	paymentUC     *usecase.PaymentUseCase
	vpnUC         *usecase.VPNUseCase
	referralUC    *usecase.ReferralUseCase
	notifUC       *usecase.NotificationUseCase
	msg           *service.MessageService
	renamingUsers map[int64]string
	mu            sync.RWMutex
}

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

func (h *BaseHandler) getVPNConnectionsBySubscriptionID(ctx context.Context, subscriptionID string) ([]*core.VPNConnection, error) {

	return []*core.VPNConnection{}, nil
}

func (h *BaseHandler) createSubscription(ctx context.Context, dto usecase.CreateSubscriptionDTO) (interface{}, error) {

	return h.subUC.CreateSubscription(ctx, dto)
}

func (h *BaseHandler) createVPNForSubscription(ctx context.Context, userID int64, subscriptionID string) (interface{}, error) {

	return h.vpnUC.CreateVPNForSubscription(ctx, userID, subscriptionID)
}

func (h *BaseHandler) updateSubscriptionName(ctx context.Context, userID int64, subscriptionID, name string) error {

	return h.subUC.UpdateSubscriptionName(ctx, userID, subscriptionID, name)
}

func (h *BaseHandler) activateTrial(ctx context.Context, userID int64) (bool, error) {

	return h.userUC.ActivateTrial(ctx, userID)
}

func (h *BaseHandler) logError(err error, context string) {
	if err != nil {
		slog.Error("Error", "context", context, "error", err)
	}
}

func (h *BaseHandler) sendError(chatID int64, message string) error {

	return h.msg.SendMessage(context.Background(), chatID, message)
}

func getUserID(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {

		return update.CallbackQuery.From.ID
	}

	return 0
}

func getChatID(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {

		return update.CallbackQuery.Message.Chat.ID
	}

	return 0
}

func getMessageID(update tgbotapi.Update) int {
	if update.CallbackQuery != nil {

		return update.CallbackQuery.Message.MessageID
	}

	return 0
}
