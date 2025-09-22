package handlers

import (
	"context"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DefaultHandler обрабатывает неизвестные команды
type DefaultHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewDefaultHandler создает новый обработчик по умолчанию
func NewDefaultHandler(useCaseManager *usecase.UseCaseManager) *DefaultHandler {
	return &DefaultHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает неизвестные команды
func (h *DefaultHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	chatID := h.GetChatID(update)

	// Отправляем сообщение о неизвестной команде
	message := `
❓ <b>Неизвестная команда</b>

Я не понимаю эту команду. Используйте кнопки ниже или команды:

<b>Основные команды:</b>
/start - Начать работу с ботом
/help - Показать справку
/profile - Мой профиль
/subscription - Управление подпиской
/vpn - VPN подключение
/promocode - Применить промокод
/referral - Реферальная программа
/settings - Настройки
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
func (h *DefaultHandler) Command() string {
	return "default"
}

// Description возвращает описание обработчика
func (h *DefaultHandler) Description() string {
	return "Обработчик по умолчанию"
}
