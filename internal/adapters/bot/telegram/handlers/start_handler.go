package handlers

import (
	"context"
	"log/slog"
	"os"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	bot      *tgbotapi.BotAPI
	notifier ports.Notifier
	userUC   *usecase.UserUseCase
	subUC    *usecase.SubscriptionUseCase
}

// NewStartHandler создает новый обработчик команды /start
func NewStartHandler(bot *tgbotapi.BotAPI, notifier ports.Notifier, userUC *usecase.UserUseCase, subUC *usecase.SubscriptionUseCase) *StartHandler {
	return &StartHandler{
		bot:      bot,
		notifier: notifier,
		userUC:   userUC,
		subUC:    subUC,
	}
}

// Handle обрабатывает команду /start
func (h *StartHandler) Handle(ctx context.Context, message *tgbotapi.Message) error {
	userID := message.From.ID
	chatID := message.Chat.ID

	slog.Info("Handling /start command", "user_id", userID)

	// Проверяем существует ли пользователь
	user, err := h.userUC.GetUser(ctx, userID)

	isNewUser := false
	if err != nil {
		// Пользователь не найден - создаем нового
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

	// Формируем приветственное сообщение
	var text string
	var keyboard tgbotapi.InlineKeyboardMarkup

	if isNewUser {
		// Для нового пользователя - приветствие с кнопкой триала
		firstName := message.From.FirstName
		if firstName == "" {
			firstName = "друг"
		}
		text = ui.GetWelcomeText(firstName, user.HasTrial)
		keyboard = ui.GetWelcomeKeyboard(user.HasTrial)
		slog.Info("Showing welcome message for new user", "user_id", userID, "is_new_user", isNewUser)
	} else {
		// Для существующего - объединенное меню с профилем
		slog.Info("Showing main menu for existing user", "user_id", userID, "is_new_user", isNewUser)

		// Проверяем активные подписки
		subscriptions, err := h.subUC.GetUserSubscriptions(ctx, userID)
		isPremium := err == nil && len(subscriptions) > 0

		statusText := "🆓 Бесплатный"
		subUntilText := ""

		if isPremium && len(subscriptions) > 0 {
			statusText = "⭐ Premium"
			subUntilText = subscriptions[0].EndDate.Format("02.01.2006")
		}

		text = ui.GetMainMenuWithProfileText(user, isPremium, statusText, subUntilText)
		keyboard = ui.GetMainMenuWithProfileKeyboard(isPremium)
	}

	// Отправляем сообщение с фото
	return h.sendMessageWithPhoto(ctx, chatID, text, keyboard)
}

// sendMessageWithPhoto отправляет сообщение с фото (универсальный метод)
func (h *StartHandler) sendMessageWithPhoto(ctx context.Context, chatID int64, caption string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	photoPath := "static/images/bot_banner.png"

	// Проверяем существование файла
	if _, fileErr := os.Stat(photoPath); fileErr == nil {
		err := h.notifier.SendPhotoFromFile(ctx, chatID, photoPath, caption, keyboard)
		if err != nil {
			// Если не удалось отправить фото, отправляем обычное сообщение
			slog.Warn("Failed to send photo, sending text message", "error", err)
			msg := tgbotapi.NewMessage(chatID, caption)
			msg.ReplyMarkup = keyboard
			msg.ParseMode = "HTML"
			_, err = h.bot.Send(msg)
			return err
		}
		return nil
	}

	// Если файла нет, отправляем обычное сообщение
	slog.Info("Photo file not found, sending text message", "path", photoPath)
	msg := tgbotapi.NewMessage(chatID, caption)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "HTML"
	_, err := h.bot.Send(msg)
	return err
}

// sendError отправляет сообщение об ошибке
func (h *StartHandler) sendError(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, "❌ "+text)
	_, err := h.bot.Send(msg)
	return err
}
