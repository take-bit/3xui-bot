package handlers

import (
	"context"

	"3xui-bot/internal/adapters/bot/telegram/handlers/callback"
	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	router *callback.Router
}

func NewCallbackHandler(
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
	bot *tgbotapi.BotAPI,
) *CallbackHandler {
	msgService := service.NewMessageService(bot)
	router := callback.NewRouter(userUC, subUC, paymentUC, vpnUC, referralUC, notifUC, msgService)

	return &CallbackHandler{
		router: router,
	}
}

func (h *CallbackHandler) CanHandle(update tgbotapi.Update) bool {

	return update.CallbackQuery != nil
}

func (h *CallbackHandler) Handle(ctx context.Context, update tgbotapi.Update) error {

	return h.router.Handle(ctx, update)
}

func (h *CallbackHandler) HandleTextMessage(ctx context.Context, userID int64, chatID int64, messageText string) (bool, error) {

	return h.router.HandleTextMessage(ctx, userID, chatID, messageText)
}
