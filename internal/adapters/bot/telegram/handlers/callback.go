package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"3xui-bot/internal/adapters/bot/telegram/ui"
	"3xui-bot/internal/core"
	"3xui-bot/internal/ports"
	"3xui-bot/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CallbackHandler обрабатывает все callback query
type CallbackHandler struct {
	userUC *usecase.UserUseCase
	subUC  *usecase.SubscriptionUseCase
	vpnUC  *usecase.VPNUseCase
	bot    ports.BotPort
	log    *slog.Logger
	route  map[string]func(context.Context, meta) error
}

// meta содержит метаданные callback query
type meta struct {
	userID    int64
	chatID    int64
	messageID int
	cbID      string
	data      string
}

// NewCallbackHandler создает новый типобезопасный обработчик callback'ов
func NewCallbackHandler(
	userUC *usecase.UserUseCase,
	subUC *usecase.SubscriptionUseCase,
	vpnUC *usecase.VPNUseCase,
	bot ports.BotPort,
	log *slog.Logger,
) *CallbackHandler {
	h := &CallbackHandler{
		userUC: userUC,
		subUC:  subUC,
		vpnUC:  vpnUC,
		bot:    bot,
		log:    log,
	}

	// Регистрируем роуты для простых (непараметризованных) callback'ов
	h.route = map[string]func(context.Context, meta) error{
		"get_trial":           h.handleGetTrial,
		"open_menu":           h.handleOpenMenu,
		"open_profile":        h.handleOpenProfile,
		"open_pricing":        h.handleOpenPricing,
		"my_subscriptions":    h.handleMySubscriptions,
		"create_subscription": h.handleCreateSubscription,
		"open_keys":           h.handleOpenKeys,
		"open_referrals":      h.handleOpenReferrals,
		"open_support":        h.handleOpenSupport,
		"my_configs":          h.handleMyConfigs,
		"referral_stats":      h.handleReferralStats,
		"my_referrals":        h.handleMyReferrals,
		"my_referral_link":    h.handleMyReferralLink,
		"create_wireguard":    h.handleCreateWireguard,
		"create_shadowsocks":  h.handleCreateShadowsocks,
	}

	return h
}

// CanHandle проверяет, может ли обработчик обработать обновление
func (h *CallbackHandler) CanHandle(update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

// Handle обрабатывает callback query
func (h *CallbackHandler) Handle(ctx context.Context, upd tgbotapi.Update) error {
	cb := upd.CallbackQuery
	if cb == nil {
		return nil
	}

	m := meta{
		userID:    cb.From.ID,
		chatID:    cb.Message.Chat.ID,
		messageID: cb.Message.MessageID,
		cbID:      cb.ID,
		data:      cb.Data,
	}

	h.info("Handling callback", "data", m.data, "user_id", m.userID)

	// Проверяем простые роуты
	if fn, ok := h.route[m.data]; ok {
		return fn(ctx, m)
	}

	// Параметризованные колбэки
	if planID, ok := ui.ParsePlanCallback(m.data); ok {
		return h.handlePlanSelection(ctx, m, planID)
	}
	if planID, ok := ui.ParseSelectPlanCallback(m.data); ok {
		return h.handlePlanSelection(ctx, m, planID)
	}
	if planID, ok := ui.ParseCreatePlanCallback(m.data); ok {
		return h.handleCreateSubscriptionByPlan(ctx, m, planID)
	}
	if subID, ok := ui.ParseViewSubscriptionCallback(m.data); ok {
		return h.handleViewSubscription(ctx, m, subID)
	}
	if subID, ok := ui.ParseRenameSubscriptionCallback(m.data); ok {
		return h.handleRenameSubscription(ctx, m, subID)
	}
	if subID, ok := ui.ParseExtendSubscriptionCallback(m.data); ok {
		return h.handleExtendSubscription(ctx, m, subID)
	}
	if subID, ok := ui.ParseDeleteSubscriptionCallback(m.data); ok {
		return h.handleDeleteSubscription(ctx, m, subID)
	}
	if planID, subID, ok := ui.ParseExtendPlanCallback(m.data); ok {
		return h.handleExtendSubscriptionByPlan(ctx, m, planID, subID)
	}
	if planID, ok := ui.ParsePayCardCallback(m.data); ok {
		return h.handlePayCard(ctx, m, planID)
	}
	if planID, ok := ui.ParsePaySBPCallback(m.data); ok {
		return h.handlePaySBP(ctx, m, planID)
	}
	if planID, ok := ui.ParsePayStarsCallback(m.data); ok {
		return h.handlePayStars(ctx, m, planID)
	}
	if configID, ok := ui.ParseViewConfigCallback(m.data); ok {
		return h.handleViewConfig(ctx, m, configID)
	}
	if subID, ok := ui.ParseConnectionGuideCallback(m.data); ok {
		return h.handleConnectionGuide(ctx, m, subID)
	}

	// Неизвестная команда
	h.warn("Unknown callback", "data", m.data)
	return h.bot.Edit(ctx, m.chatID, m.messageID, "❓ Неизвестная команда", ui.GetUnknownCommandKeyboard())
}

// ============================================================================
// ОСНОВНЫЕ ОБРАБОТЧИКИ
// ============================================================================

func (h *CallbackHandler) handleGetTrial(ctx context.Context, m meta) error {
	h.info("Handle get trial", "user_id", m.userID)

	// Получаем пользователя
	user, err := h.userUC.GetUser(ctx, m.userID)
	if err != nil || user == nil {
		h.err("GetUser", err)
		return err
	}

	// Активируем пробный доступ
	success, err := h.userUC.ActivateTrial(ctx, m.userID)
	if err != nil {
		h.err("ActivateTrial", err)
		return err
	}

	var text string
	if success {
		text = "🎉 Пробный доступ активирован на 3 дня!"
	} else {
		text = "❌ Пробный доступ уже был использован"
	}

	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetWelcomeKeyboard(user.HasTrial))
}

func (h *CallbackHandler) handleOpenMenu(ctx context.Context, m meta) error {
	h.info("Handle open menu", "user_id", m.userID)

	user, err := h.userUC.GetUser(ctx, m.userID)
	if err != nil || user == nil {
		h.err("GetUser", err)
		return err
	}

	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetWelcomeText(user.FirstName, user.HasTrial), ui.GetWelcomeKeyboard(user.HasTrial))
}

func (h *CallbackHandler) handleOpenProfile(ctx context.Context, m meta) error {
	h.info("Handle open profile", "user_id", m.userID)

	user, err := h.userUC.GetUser(ctx, m.userID)
	if err != nil || user == nil {
		h.err("GetUser", err)
		return err
	}

	// Проверяем наличие активной подписки
	subscriptions, _ := h.subUC.GetUserSubscriptions(ctx, m.userID)
	isPremium := false
	statusText := "Free"
	subUntilText := "—"

	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if sub.IsActive {
				isPremium = true
				statusText = "Premium"
				subUntilText = sub.EndDate.Format("02.01.2006")
				break
			}
		}
	}

	text := ui.GetProfileText(user, isPremium, statusText, subUntilText)
	keyboard := ui.GetProfileKeyboard(isPremium)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenPricing(ctx context.Context, m meta) error {
	h.info("Handle open pricing", "user_id", m.userID)

	plans, err := h.subUC.GetPlans(ctx)
	if err != nil {
		h.err("GetPlans", err)
		return err
	}

	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetPricingText(plans), ui.GetPricingKeyboard(plans))
}

func (h *CallbackHandler) handleMySubscriptions(ctx context.Context, m meta) error {
	h.info("Handle my subscriptions", "user_id", m.userID)

	subscriptions, err := h.subUC.GetUserSubscriptions(ctx, m.userID)
	if err != nil {
		h.err("GetUserSubscriptions", err)
		return err
	}

	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetSubscriptionsText(subscriptions), ui.GetSubscriptionsKeyboard(subscriptions))
}

func (h *CallbackHandler) handleCreateSubscription(ctx context.Context, m meta) error {
	h.info("Handle create subscription", "user_id", m.userID)

	plans, err := h.subUC.GetPlans(ctx)
	if err != nil {
		h.err("GetPlans", err)
		return err
	}

	text := ui.GetCreateSubscriptionText()
	keyboard := ui.GetCreateSubscriptionKeyboard(plans)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleOpenKeys(ctx context.Context, m meta) error {
	h.info("Handle open keys", "user_id", m.userID)
	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetKeysText(), ui.GetKeysKeyboard())
}

func (h *CallbackHandler) handleOpenReferrals(ctx context.Context, m meta) error {
	h.info("Handle open referrals", "user_id", m.userID)
	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetReferralsText(), ui.GetReferralsKeyboard())
}

func (h *CallbackHandler) handleOpenSupport(ctx context.Context, m meta) error {
	h.info("Handle open support", "user_id", m.userID)

	user, err := h.userUC.GetUser(ctx, m.userID)
	if err != nil || user == nil {
		h.err("GetUser", err)
		return err
	}

	return h.bot.Edit(ctx, m.chatID, m.messageID, ui.GetSupportText(), ui.GetWelcomeKeyboard(user.HasTrial))
}

func (h *CallbackHandler) handleMyConfigs(ctx context.Context, m meta) error {
	h.info("Handle my configs", "user_id", m.userID)
	text := "📋 Ваши VPN конфигурации\n\nПока конфигураций нет."
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetKeysKeyboard())
}

func (h *CallbackHandler) handleReferralStats(ctx context.Context, m meta) error {
	h.info("Handle referral stats", "user_id", m.userID)
	text := "📊 Статистика рефералов\n\nПока статистики нет."
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetReferralsKeyboard())
}

func (h *CallbackHandler) handleMyReferrals(ctx context.Context, m meta) error {
	h.info("Handle my referrals", "user_id", m.userID)
	text := "👥 Ваши рефералы\n\nПока рефералов нет."
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetReferralsKeyboard())
}

func (h *CallbackHandler) handleMyReferralLink(ctx context.Context, m meta) error {
	h.info("Handle my referral link", "user_id", m.userID)
	text := "🔗 Ваша реферальная ссылка\n\nhttps://t.me/your_bot?start=ref_123456"
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetReferralsKeyboard())
}

func (h *CallbackHandler) handleCreateWireguard(ctx context.Context, m meta) error {
	h.info("Handle create wireguard", "user_id", m.userID)
	text := "🔑 Создание WireGuard конфигурации\n\nВведите название для конфигурации:"
	return h.bot.Send(ctx, m.chatID, text, nil)
}

func (h *CallbackHandler) handleCreateShadowsocks(ctx context.Context, m meta) error {
	h.info("Handle create shadowsocks", "user_id", m.userID)
	text := "🔑 Создание Shadowsocks конфигурации\n\nВведите название для конфигурации:"
	return h.bot.Send(ctx, m.chatID, text, nil)
}

// ============================================================================
// ПАРАМЕТРИЗОВАННЫЕ ОБРАБОТЧИКИ
// ============================================================================

func (h *CallbackHandler) handlePlanSelection(ctx context.Context, m meta, planID string) error {
	h.info("Handle plan selection", "plan_id", planID, "user_id", m.userID)

	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	text := fmt.Sprintf("📦 План: %s\n💵 Цена: %.0f₽\n⏰ Длительность: %d дней\n\nВыберите способ оплаты:", plan.Name, plan.Price, plan.Days)
	keyboard := ui.GetPaymentMethodKeyboard(planID)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleCreateSubscriptionByPlan(ctx context.Context, m meta, planID string) error {
	h.info("Handle create subscription by plan", "plan_id", planID, "user_id", m.userID)

	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	// Создаем подписку
	dto := usecase.CreateSubscriptionDTO{
		UserID: m.userID,
		Name:   "Основная",
		PlanID: planID,
		Days:   plan.Days,
	}

	subscription, err := h.subUC.CreateSubscription(ctx, dto)
	if err != nil {
		h.err("CreateSubscription", err)
		return err
	}

	// Создаем VPN для подписки
	_, err = h.vpnUC.CreateVPNForSubscription(ctx, m.userID, subscription.ID)
	if err != nil {
		h.err("CreateVPN", err)
	}

	user, _ := h.userUC.GetUser(ctx, m.userID)
	text := fmt.Sprintf("✅ Подписка '%s' создана успешно!\n⏰ Длительность: %d дней", plan.Name, plan.Days)
	hasT := false
	if user != nil {
		hasT = user.HasTrial
	}
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetWelcomeKeyboard(hasT))
}

func (h *CallbackHandler) handleViewSubscription(ctx context.Context, m meta, subscriptionID string) error {
	h.info("Handle view subscription", "subscription_id", subscriptionID, "user_id", m.userID)

	subscription, err := h.subUC.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		h.err("GetSubscription", err)
		return err
	}

	// Проверка прав доступа
	if subscription.UserID != m.userID {
		h.warn("Access denied to subscription", "user_id", m.userID, "owner_id", subscription.UserID)
		return h.bot.Edit(ctx, m.chatID, m.messageID, "❌ У вас нет доступа к этой подписке", nil)
	}

	// Получаем план
	plan, err := h.subUC.GetPlan(ctx, subscription.PlanID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	// Получаем VPN конфигурации
	vpnConfigs, err := h.vpnUC.GetVPNConnectionsBySubscription(ctx, subscriptionID)
	if err != nil {
		h.err("GetVPNConnections", err)
		vpnConfigs = []*core.VPNConnection{}
	}

	text := ui.GetSubscriptionDetailText(subscription, plan, vpnConfigs)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleRenameSubscription(ctx context.Context, m meta, subscriptionID string) error {
	h.info("Handle rename subscription", "subscription_id", subscriptionID, "user_id", m.userID)

	subscription, err := h.subUC.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		h.err("GetSubscription", err)
		return err
	}

	if subscription.UserID != m.userID {
		return h.bot.Edit(ctx, m.chatID, m.messageID, "❌ У вас нет доступа к этой подписке", nil)
	}

	text := "✏️ Введите новое название для подписки:\n\n(Максимум 50 символов)"
	keyboard := ui.GetCancelKeyboard()
	return h.bot.Send(ctx, m.chatID, text, keyboard)
}

func (h *CallbackHandler) handleExtendSubscription(ctx context.Context, m meta, subscriptionID string) error {
	h.info("Handle extend subscription", "subscription_id", subscriptionID, "user_id", m.userID)

	subscription, err := h.subUC.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		h.err("GetSubscription", err)
		return err
	}

	if subscription.UserID != m.userID {
		return h.bot.Edit(ctx, m.chatID, m.messageID, "❌ У вас нет доступа к этой подписке", nil)
	}

	plans, err := h.subUC.GetPlans(ctx)
	if err != nil {
		h.err("GetPlans", err)
		return err
	}

	text := ui.GetExtendSubscriptionText(subscription)
	keyboard := ui.GetExtendSubscriptionKeyboard(subscriptionID, plans)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleDeleteSubscription(ctx context.Context, m meta, subscriptionID string) error {
	h.info("Handle delete subscription", "subscription_id", subscriptionID, "user_id", m.userID)

	subscription, err := h.subUC.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		h.err("GetSubscription", err)
		return err
	}

	if subscription.UserID != m.userID {
		return h.bot.Edit(ctx, m.chatID, m.messageID, "❌ У вас нет доступа к этой подписке", nil)
	}

	vpnConfigs, _ := h.vpnUC.GetVPNConnectionsBySubscription(ctx, subscriptionID)
	text := ui.GetDeleteSubscriptionText(subscription)
	keyboard := ui.GetSubscriptionDetailKeyboard(subscription, vpnConfigs)
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, keyboard)
}

func (h *CallbackHandler) handleExtendSubscriptionByPlan(ctx context.Context, m meta, planID, subscriptionID string) error {
	h.info("Handle extend subscription by plan", "subscription_id", subscriptionID, "plan_id", planID, "user_id", m.userID)

	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	err = h.subUC.ExtendSubscription(ctx, m.userID, subscriptionID, plan.Days)
	if err != nil {
		h.err("ExtendSubscription", err)
		return err
	}

	user, _ := h.userUC.GetUser(ctx, m.userID)
	text := fmt.Sprintf("✅ Подписка продлена на %d дней!", plan.Days)
	hasT := false
	if user != nil {
		hasT = user.HasTrial
	}
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetWelcomeKeyboard(hasT))
}

func (h *CallbackHandler) handlePayCard(ctx context.Context, m meta, planID string) error {
	h.info("Handle pay card", "plan_id", planID, "user_id", m.userID)

	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	// MOCK: Картой всегда успешно
	h.info("Creating MOCK card payment (auto-success)", "user_id", m.userID, "plan_id", planID)

	dto := usecase.CreateSubscriptionDTO{
		UserID: m.userID,
		Name:   "Основная",
		PlanID: planID,
		Days:   plan.Days,
	}

	subscription, err := h.subUC.CreateSubscription(ctx, dto)
	if err != nil {
		h.err("CreateSubscription", err)
		return h.bot.Send(ctx, m.chatID, "❌ Ошибка создания подписки", nil)
	}

	// Создаем VPN для подписки
	_, err = h.vpnUC.CreateVPNForSubscription(ctx, m.userID, subscription.ID)
	if err != nil {
		h.err("CreateVPN", err)
	}

	text := fmt.Sprintf("✅ Оплата успешна!\n\n🎉 Подписка '%s' активирована на %d дней", plan.Name, plan.Days)
	return h.bot.Send(ctx, m.chatID, text, ui.GetWelcomeKeyboard(false))
}

func (h *CallbackHandler) handlePaySBP(ctx context.Context, m meta, planID string) error {
	h.info("Handle pay SBP", "plan_id", planID, "user_id", m.userID)
	// MOCK: СБП работает аналогично карте
	return h.handlePayCard(ctx, m, planID)
}

func (h *CallbackHandler) handlePayStars(ctx context.Context, m meta, planID string) error {
	h.info("Handle pay stars", "plan_id", planID, "user_id", m.userID)

	plan, err := h.subUC.GetPlan(ctx, planID)
	if err != nil {
		h.err("GetPlan", err)
		return err
	}

	text := fmt.Sprintf("⭐ Оплата Telegram Stars\n\nПлан: %s\nЦена: %.0f₽\n\nФункция в разработке", plan.Name, plan.Price)
	return h.bot.Send(ctx, m.chatID, text, ui.GetBackToPricingKeyboard())
}

func (h *CallbackHandler) handleViewConfig(ctx context.Context, m meta, configID string) error {
	h.info("Handle view config", "config_id", configID, "user_id", m.userID)
	text := "🔑 Детали VPN конфигурации\n\nФункция в разработке"
	return h.bot.Edit(ctx, m.chatID, m.messageID, text, ui.GetKeysKeyboard())
}

func (h *CallbackHandler) handleConnectionGuide(ctx context.Context, m meta, subscriptionID string) error {
	h.info("Handle connection guide", "subscription_id", subscriptionID, "user_id", m.userID)

	subscription, err := h.subUC.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		h.err("GetSubscription", err)
		return h.bot.Send(ctx, m.chatID, "❌ Подписка не найдена", nil)
	}

	if subscription.UserID != m.userID {
		return h.bot.Send(ctx, m.chatID, "❌ У вас нет доступа к этой подписке", nil)
	}

	if !subscription.IsActive {
		return h.bot.Send(ctx, m.chatID, "❌ Подписка неактивна", nil)
	}

	connectionURL := fmt.Sprintf("https://3xui.com/connect/%s", subscription.ID)

	text := fmt.Sprintf(`📖 *Инструкция по подключению*

*🔗 URL подключения:*
`+"`%s`"+`

*📱 Для подключения выполните следующие шаги:*

1️⃣ *Скачайте VPN клиент:*
   • WireGuard для Android/iOS/Windows/Mac
   • Или используйте встроенный клиент

2️⃣ *Импортируйте конфигурацию:*
   • Откройте приложение WireGuard
   • Нажмите "Добавить туннель"
   • Выберите "Импортировать из файла или архива"

3️⃣ *Подключитесь:*
   • Найдите вашу конфигурацию в списке
   • Нажмите переключатель для подключения
   • Готово! Ваш трафик защищен

*💡 Полезные советы:*
   • Держите приложение обновленным
   • При проблемах попробуйте переподключиться
   • Следите за сроком действия подписки

*🆘 Нужна помощь?*
   Обратитесь в поддержку через кнопку ниже`, connectionURL)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🆘 Поддержка", "open_support"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к подписке", fmt.Sprintf("view_subscription_%s", subscriptionID)),
		),
	)

	return h.bot.Send(ctx, m.chatID, text, keyboard)
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
// ============================================================================

func (h *CallbackHandler) info(msg string, args ...any) {
	if h.log != nil {
		h.log.Info(msg, args...)
	} else {
		slog.Info(msg, args...)
	}
}

func (h *CallbackHandler) warn(msg string, args ...any) {
	if h.log != nil {
		h.log.Warn(msg, args...)
	} else {
		slog.Warn(msg, args...)
	}
}

func (h *CallbackHandler) err(msg string, e error, args ...any) {
	if e == nil {
		return
	}
	allArgs := append([]any{"error", e}, args...)
	if h.log != nil {
		h.log.Error(msg, allArgs...)
	} else {
		slog.Error(msg, allArgs...)
	}
}
