package callback

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
)

func (h *BaseHandler) HandleOpenKeys(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open keys", "user_id", userID)

	text := ui.GetKeysText()
	keyboard := ui.GetKeysKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleMyConfigs(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my configs", "user_id", userID)

	text := "🔑 Мои конфигурации\n\nЗдесь будут ваши VPN конфигурации"
	keyboard := ui.GetBackToPricingKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleCreateWireguard(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create wireguard", "user_id", userID)
	text := "🔑 Создание WireGuard конфигурации\n\nВведите название для конфигурации:"

	return h.msg.SendMessage(ctx, chatID, text)
}

func (h *BaseHandler) HandleCreateShadowsocks(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling create shadowsocks", "user_id", userID)
	text := "🔑 Создание Shadowsocks конфигурации\n\nВведите название для конфигурации:"

	return h.msg.SendMessage(ctx, chatID, text)
}
