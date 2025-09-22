package handlers

import (
	"context"

	"3xui-bot/internal/interfaces"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpHandler обрабатывает команду /help
type HelpHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewHelpHandler создает новый обработчик команды /help
func NewHelpHandler(useCaseManager *usecase.UseCaseManager) *HelpHandler {
	return &HelpHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /help
func (h *HelpHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	chatID := h.GetChatID(update)

	// Отправляем справочное сообщение
	message := `
❓ <b>Справка по боту</b>

<b>Основные команды:</b>
/start - Начать работу с ботом
/help - Показать эту справку
/profile - Мой профиль
/subscription - Управление подпиской
/vpn - VPN подключение
/promocode - Применить промокод
/referral - Реферальная программа
/settings - Настройки

<b>Как пользоваться:</b>
1️⃣ Нажмите /start для регистрации
2️⃣ Получите пробный период на 3 дня
3️⃣ Скачайте конфигурацию VPN
4️⃣ Настройте VPN на своем устройстве
5️⃣ Продлите подписку для продолжения использования

<b>Поддержка:</b>
Если у вас возникли вопросы, обратитесь в поддержку через кнопку "Настройки" → "Поддержка"

<b>Реферальная программа:</b>
Приглашайте друзей и получайте вознаграждения!
Ваша реферальная ссылка доступна в разделе "Рефералы"
`

	// Создаем клавиатуру с основными разделами
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
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *HelpHandler) Command() string {
	return "help"
}

// Description возвращает описание обработчика
func (h *HelpHandler) Description() string {
	return "Показать справку по боту"
}
