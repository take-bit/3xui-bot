package handlers

import (
	"context"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SettingsHandler обрабатывает команду /settings
type SettingsHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewSettingsHandler создает новый обработчик команды /settings
func NewSettingsHandler(useCaseManager *usecase.UseCaseManager) *SettingsHandler {
	return &SettingsHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /settings
func (h *SettingsHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	chatID := h.GetChatID(update)

	// Отправляем информацию о настройках
	message := `
⚙️ <b>Настройки</b>

<b>Доступные настройки:</b>

🌐 <b>Язык</b> - Изменить язык интерфейса
🔔 <b>Уведомления</b> - Настройка уведомлений
📊 <b>Статистика</b> - Просмотр статистики использования
🆘 <b>Поддержка</b> - Связаться с поддержкой
ℹ️ <b>О боте</b> - Информация о боте

<b>Дополнительные функции:</b>
📱 <b>Экспорт данных</b> - Скачать свои данные
🗑️ <b>Удалить аккаунт</b> - Удалить аккаунт и все данные
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "settings_language"),
			tgbotapi.NewInlineKeyboardButtonData("🔔 Уведомления", "settings_notifications"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "settings_stats"),
			tgbotapi.NewInlineKeyboardButtonData("🆘 Поддержка", "settings_support"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ℹ️ О боте", "settings_about"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 Экспорт данных", "settings_export"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить аккаунт", "settings_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *SettingsHandler) Command() string {
	return "settings"
}

// Description возвращает описание обработчика
func (h *SettingsHandler) Description() string {
	return "Настройки бота"
}
