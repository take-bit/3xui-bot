package callback

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
)

func (h *BaseHandler) HandleOpenSupport(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open support", "user_id", userID)
	text := ui.GetSupportText()
	keyboard := ui.GetBackToPricingKeyboard()
	return h.msg.EditMessageText(ctx, chatID, messageID, text, keyboard)
}
