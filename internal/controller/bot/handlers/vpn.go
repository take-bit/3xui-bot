package handlers

import (
	"context"
	"fmt"

	"3xui-bot/internal/domain"
	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// VPNHandler обрабатывает команду /vpn
type VPNHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewVPNHandler создает новый обработчик команды /vpn
func NewVPNHandler(useCaseManager *usecase.UseCaseManager) *VPNHandler {
	return &VPNHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /vpn
func (h *VPNHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)

	// Получаем информацию о VPN подключении
	connection, err := h.UseCaseManager.GetVPNUseCase().GetVPNConnectionInfo(ctx, userID)
	if err != nil {
		// Если подключение не найдено, предлагаем создать
		return h.handleNoConnection(ctx, chatID, userID)
	}

	// Отправляем информацию о подключении
	return h.sendConnectionInfo(ctx, chatID, connection)
}

// handleNoConnection обрабатывает случай, когда VPN подключение не найдено
func (h *VPNHandler) handleNoConnection(ctx context.Context, chatID, userID int64) error {
	message := `
🔗 <b>VPN подключение</b>

❌ У вас нет активного VPN подключения

💡 Для создания подключения:
1️⃣ Убедитесь, что у вас есть активная подписка
2️⃣ Нажмите кнопку "Создать подключение" ниже
3️⃣ Скачайте конфигурацию
4️⃣ Настройте VPN на своем устройстве
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подключение", "vpn_create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// sendConnectionInfo отправляет информацию о VPN подключении
func (h *VPNHandler) sendConnectionInfo(ctx context.Context, chatID int64, connection *domain.VPNConnection) error {
	message := fmt.Sprintf(`
🔗 <b>VPN подключение</b>

✅ <b>Статус:</b> Активно
🆔 <b>ID подключения:</b> <code>%s</code>
🌍 <b>Сервер:</b> %s
📅 <b>Создано:</b> %s
⏰ <b>Истекает:</b> %s

<b>Конфигурация:</b>
🔗 <b>Ссылка для скачивания:</b>
<code>%s</code>

<b>Инструкция по настройке:</b>
1️⃣ Скачайте конфигурацию по ссылке выше
2️⃣ Установите VPN клиент (WireGuard, V2Ray, etc.)
3️⃣ Импортируйте конфигурацию
4️⃣ Подключитесь к VPN
`,
		connection.UUID,
		connection.ServerID,
		connection.CreatedAt.Format("02.01.2006 15:04"),
		connection.ExpiresAt.Format("02.01.2006 15:04"),
		connection.ConfigURL,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📥 Скачать конфигурацию", connection.ConfigURL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить подключение", "vpn_refresh"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить подключение", "vpn_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *VPNHandler) Command() string {
	return "vpn"
}

// Description возвращает описание обработчика
func (h *VPNHandler) Description() string {
	return "Управление VPN подключением"
}
