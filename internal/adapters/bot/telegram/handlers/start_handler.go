package handlers

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/service"
	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	bot      *tgbotapi.BotAPI
	notifier ports.Notifier
	userUC   *usecase.UserUseCase
	subUC    *usecase.SubscriptionUseCase
	msg      *service.MessageService
}

func NewStartHandler(bot *tgbotapi.BotAPI, notifier ports.Notifier, userUC *usecase.UserUseCase, subUC *usecase.SubscriptionUseCase) *StartHandler {

	return &StartHandler{
		bot:      bot,
		notifier: notifier,
		userUC:   userUC,
		subUC:    subUC,
		msg:      service.NewMessageService(bot),
	}
}

func (h *StartHandler) Handle(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	chatID := message.Chat.ID

	slog.Info("Handling /start command", "user_id", userID)

	user, err := h.userUC.GetUser(ctx, userID)

	isNewUser := false
	if err != nil {
		slog.Info("Creating new user", "user_id", userID)

		createUserDTO := usecase.CreateUserDTO{
			TelegramID:   userID,
			Username:     message.From.UserName,
			FirstName:    message.From.FirstName,
			LastName:     message.From.LastName,
			LanguageCode: message.From.LanguageCode,
		}

		user, err = h.userUC.CreateUser(ctx, createUserDTO)
		if err != nil {
			slog.Error("Failed to create user", "user_id", userID, "error", err)

			return h.sendError(chatID, "Произошла ошибка при регистрации. Попробуйте еще раз.")
		}

		isNewUser = true
		slog.Info("New user created", "user_id", userID)
	}

	if isNewUser {
		firstName := message.From.FirstName
		if firstName == "" {
			firstName = "друг"
		}
		text := ui.GetWelcomeText(firstName, user.HasTrial)
		keyboard := ui.GetWelcomeKeyboard(user.HasTrial)
		slog.Info("Showing welcome message for new user", "user_id", userID, "is_new_user", isNewUser)

		return h.msg.SendPhotoWithMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
	}

	slog.Info("Showing main menu for existing user", "user_id", userID, "is_new_user", isNewUser)

	subscriptions, err := h.subUC.GetUserSubscriptions(ctx, userID)
	if err != nil {
		subscriptions = []*core.Subscription{}
	}

	isPremium := false
	for _, sub := range subscriptions {
		if sub.IsActive && !sub.IsExpired() {
			isPremium = true
			break
		}
	}

	text := ui.GetMainMenuWithProfileText(user, subscriptions)
	keyboard := ui.GetMainMenuWithProfileKeyboard(isPremium)

	return h.msg.SendPhotoWithPreEscapedMarkdown(ctx, chatID, "static/images/bot_banner.png", text, keyboard)
}

func (h *StartHandler) sendError(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	_, err := h.bot.Send(msg)

	return err
}
