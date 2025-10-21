package callback

import (
	"context"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
)

func (h *BaseHandler) HandleOpenReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling open referrals", "user_id", userID)

	text := ui.GetReferralsText()
	keyboard := ui.GetReferralsKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleReferralStats(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling referral stats", "user_id", userID)

	text := ui.GetReferralRankingText()
	keyboard := ui.GetReferralRankingKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleMyReferrals(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my referrals", "user_id", userID)

	text := "👥 Мои реферралы\n\nЗдесь будет информация о ваших реферралах"
	keyboard := ui.GetBackToPricingKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}

func (h *BaseHandler) HandleMyReferralLink(ctx context.Context, userID, chatID int64, messageID int) error {
	slog.Info("Handling my referral link", "user_id", userID)

	text := "🔗 Моя реферальная ссылка\n\nЗдесь будет ваша реферальная ссылка"
	keyboard := ui.GetBackToPricingKeyboard()

	return h.msg.DeleteAndSendMessage(ctx, chatID, messageID, text, keyboard)
}
