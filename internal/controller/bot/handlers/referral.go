package handlers

import (
	"context"
	"fmt"

	"3xui-bot/internal/usecase"
	"3xui-bot/internal/interfaces"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ReferralHandler обрабатывает команду /referral
type ReferralHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewReferralHandler создает новый обработчик команды /referral
func NewReferralHandler(useCaseManager *usecase.UseCaseManager) *ReferralHandler {
	return &ReferralHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /referral
func (h *ReferralHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	userID := h.GetUserID(update)
	chatID := h.GetChatID(update)

	// Получаем статистику рефералов
	referralStats, err := h.UseCaseManager.GetReferralUseCase().GetReferralStats(ctx, userID)
	if err != nil {
		h.HandleError(ctx, chatID, err, "❌ Ошибка при получении статистики рефералов. Попробуйте позже.")
		return err
	}

	// Отправляем информацию о рефералах
	return h.sendReferralInfo(ctx, chatID, referralStats)
}

// sendReferralInfo отправляет информацию о рефералах
func (h *ReferralHandler) sendReferralInfo(ctx context.Context, chatID int64, referralStats *usecase.ReferralStatsInfo) error {
	message := fmt.Sprintf(`
👥 <b>Реферальная программа</b>

<b>Ваша реферальная ссылка:</b>
<code>%s</code>

<b>Статистика:</b>
👥 Всего рефералов: %d
💰 Заработано дней: %d
⏳ Ожидает выплаты: %d
✅ Выплачено: %d

<b>Как работает программа:</b>
1️⃣ Поделитесь своей реферальной ссылкой с друзьями
2️⃣ Друзья регистрируются по вашей ссылке
3️⃣ После их первого платежа вы получаете вознаграждение
4️⃣ Вознаграждение зачисляется на ваш аккаунт

<b>Условия:</b>
• Вознаграждение начисляется после первого платежа реферала
• Размер вознаграждения: 10%% от суммы платежа
• Вознаграждение зачисляется в виде дополнительных дней подписки
`,
		referralStats.ReferralLink,
		referralStats.ReferralStats.TotalReferrals,
		referralStats.TotalEarnings,
		referralStats.PendingEarnings,
		referralStats.PaidEarnings,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📤 Поделиться ссылкой", "referral_share"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Подробная статистика", "referral_stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Промокод", "promocode"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Мой профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *ReferralHandler) Command() string {
	return "referral"
}

// Description возвращает описание обработчика
func (h *ReferralHandler) Description() string {
	return "Реферальная программа"
}
