package handlers

import (
	"context"

	"3xui-bot/internal/interfaces"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PaymentHandler обрабатывает команду /payment
type PaymentHandler struct {
	interfaces.BaseHandlerInterface
	UseCaseManager *usecase.UseCaseManager
}

// NewPaymentHandler создает новый обработчик команды /payment
func NewPaymentHandler(useCaseManager *usecase.UseCaseManager) *PaymentHandler {
	return &PaymentHandler{
		UseCaseManager: useCaseManager,
	}
}

// Handle обрабатывает команду /payment
func (h *PaymentHandler) Handle(ctx context.Context, update tgbotapi.Update) error {
	chatID := h.GetChatID(update)

	// Отправляем информацию о платежах
	message := `
💳 <b>Платежи</b>

💡 <b>Доступные способы оплаты:</b>

💰 <b>Криптовалюты</b>
• Bitcoin (BTC)
• Ethereum (ETH)
• USDT (TRC20)
• USDC (ERC20)

💳 <b>Банковские карты</b>
• Visa
• MasterCard
• МИР

📱 <b>Электронные кошельки</b>
• YooMoney
• QIWI
• WebMoney

⭐ <b>Telegram Stars</b>
• Платежи через Telegram

<b>Выберите способ оплаты:</b>
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Криптовалюты", "payment_crypto"),
			tgbotapi.NewInlineKeyboardButtonData("💳 Банковские карты", "payment_cards"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 Электронные кошельки", "payment_wallets"),
			tgbotapi.NewInlineKeyboardButtonData("⭐ Telegram Stars", "payment_stars"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Подписка", "subscription"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "main_menu"),
		),
	)

	return h.SendMessage(ctx, chatID, message, keyboard)
}

// Command возвращает команду обработчика
func (h *PaymentHandler) Command() string {
	return "payment"
}

// Description возвращает описание обработчика
func (h *PaymentHandler) Description() string {
	return "Управление платежами"
}
