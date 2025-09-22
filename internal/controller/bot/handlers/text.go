package handlers

import (
	"context"
	"fmt"
	"strings"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TextHandler обрабатывает текстовые сообщения
type TextHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewTextHandler создает новый обработчик текстовых сообщений
func NewTextHandler(useCaseManager *usecase.UseCaseManager) *TextHandler {
	return &TextHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает текстовые сообщения
func (h *TextHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)
	text := strings.TrimSpace(h.GetText(update))

	// Проверяем, что пользователь зарегистрирован
	err := h.EnsureUserRegistered(ctx, update)
	if err != nil {
		h.HandleError(ctx, chatID, err, "❌ Ошибка при регистрации. Попробуйте позже.")
		return err
	}

	// Обрабатываем различные типы текстовых сообщений
	if h.isPromocode(text) {
		return h.handlePromocode(ctx, userID, chatID, text)
	}

	if h.isSupportMessage(text) {
		return h.handleSupportMessage(ctx, userID, chatID, text)
	}

	// Если сообщение не распознано, показываем главное меню
	return h.showMainMenu(ctx, chatID)
}

// isPromocode проверяет, является ли текст промокодом
func (h *TextHandler) isPromocode(text string) bool {
	// Простая проверка: промокод должен быть в верхнем регистре и содержать только буквы и цифры
	if len(text) < 3 || len(text) > 20 {
		return false
	}

	// Проверяем, что текст содержит только буквы и цифры
	for _, char := range text {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// isSupportMessage проверяет, является ли сообщение обращением в поддержку
func (h *TextHandler) isSupportMessage(text string) bool {
	supportKeywords := []string{"поддержка", "помощь", "проблема", "ошибка", "не работает", "help", "support"}
	textLower := strings.ToLower(text)

	for _, keyword := range supportKeywords {
		if strings.Contains(textLower, keyword) {
			return true
		}
	}

	return false
}

// handlePromocode обрабатывает промокод
func (h *TextHandler) handlePromocode(ctx context.Context, userID, chatID int64, promocode string) error {
	// Применяем промокод
	result, err := h.UseCaseManager.GetPromocodeUseCase().ApplyPromocode(ctx, userID, promocode)
	if err != nil {
		h.HandleError(ctx, chatID, err, "❌ Промокод недействителен или уже использован.")
		return err
	}

	// Отправляем сообщение об успешном применении промокода
	message := fmt.Sprintf(`
🎉 <b>Промокод применен!</b>

%s

<b>Что дальше:</b>
• Проверьте свой профиль для подтверждения
• Используйте VPN подключение
• Поделитесь с друзьями реферальной ссылкой
`, result.Message)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// handleSupportMessage обрабатывает обращение в поддержку
func (h *TextHandler) handleSupportMessage(ctx context.Context, userID, chatID int64, message string) error {
	// Отправляем сообщение в поддержку
	supportMessage := fmt.Sprintf(`
🆘 <b>Обращение в поддержку</b>

<b>Пользователь:</b> @%s (ID: %d)
<b>Сообщение:</b> %s

<b>Время:</b> %s
`,
		h.GetUsername(tgbotapi.Update{}), // TODO: Получить username из контекста
		userID,
		message,
		"сейчас", // TODO: Получить текущее время
	)

	// TODO: Отправить сообщение в чат поддержки
	// h.SendMessage(ctx, supportChatID, supportMessage, nil)
	_ = supportMessage

	// Отправляем подтверждение пользователю
	responseMessage := `
✅ <b>Сообщение отправлено в поддержку</b>

Ваше обращение получено. Мы ответим в ближайшее время.

<b>Что вы можете сделать:</b>
• Проверить FAQ в разделе "Помощь"
• Обратиться в поддержку через кнопку "Настройки"
• Продолжить пользоваться ботом
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, responseMessage, keyboard)
}

// showMainMenu показывает главное меню
func (h *TextHandler) showMainMenu(ctx context.Context, chatID int64) error {
	message := `
🏠 <b>Главное меню</b>

Выберите нужный раздел:
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 VPN подключение", "vpn"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Рефералы", "referral"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", "help"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *TextHandler) Command() string {
	return "text"
}

// Description возвращает описание обработчика
func (h *TextHandler) Description() string {
	return "Обработка текстовых сообщений"
}
