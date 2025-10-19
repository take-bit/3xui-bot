package handlers

import (
	"context"

	"3xui-bot/internal/adapters/bot/telegram/handlers/callbacks"
	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackHandler обрабатывает все callback query (refactored version)
type CallbackHandler struct {
	router *callbacks.Router
}

// NewCallbackHandler создает новый обработчик callback'ов
func NewCallbackHandler(
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	paymentUC *usecase.PaymentUseCase,
	vpnUC *usecase.VPNUseCase,
	referralUC *usecase.ReferralUseCase,
	notifUC *usecase.NotificationUseCase,
	bot *tgbotapi.BotAPI,
) *CallbackHandler {
	// Create MessageService
	msgService := service.NewMessageService(bot)
	
	// Create Router with all dependencies
	router := callbacks.NewRouter(userUC, subUC, paymentUC, vpnUC, referralUC, notifUC, msgService)
	
	return &CallbackHandler{
		router: router,
	}
}

// CanHandle проверяет, может ли обработчик обработать обновление
func (h *CallbackHandler) CanHandle(update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

// Handle обрабатывает callback query
func (h *CallbackHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	return h.router.Handle(ctx, update)
}

// HandleTextMessage обрабатывает текстовые сообщения (для переименования подписок)
func (h *CallbackHandler) HandleTextMessage(ctx context.Context, userID int64, chatID int64, messageText string) (bool, error) {
	return h.router.HandleTextMessage(ctx, userID, chatID, messageText)
}
