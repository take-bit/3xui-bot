package ui

import (
	"3xui-bot/internal/core"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetWelcomeKeyboard(hasTrialUsed bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	if !hasTrialUsed {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉 Получить пробный доступ", "get_trial"),
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔑 Мои подписки", "my_subscriptions"),
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Профиль", "open_profile"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Пригласить друзей", "open_referrals"),
			tgbotapi.NewInlineKeyboardButtonData("📖 Инструкция", "show_instruction"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "open_support"),
		),
	)
}
func GetMainMenuWithProfileKeyboard(isPremium bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		tgbotapi.NewInlineKeyboardButtonData("👥 Пригласить друзей", "open_referrals"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📖 Инструкция", "show_instruction"),
		tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "open_support"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetProfileKeyboard(isPremium bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👥 Реферальная программа", "open_referrals"),
		tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "open_support"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetPricingKeyboard(plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 %s - %.0f₽ (%s)", plan.Name, plan.Price, FormatDuration(plan.Days))
			callbackData := fmt.Sprintf("select_plan_%s", plan.ID)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetPaymentMethodKeyboard(planID string) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Картой", fmt.Sprintf("pay_card_%s", planID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏦 СБП", fmt.Sprintf("pay_sbp_%s", planID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💎 Stars", fmt.Sprintf("pay_stars_%s", planID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_pricing"),
		),
	)
}
func GetSubscriptionsKeyboard(subscriptions []*core.Subscription) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	if len(subscriptions) == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	} else {
		for _, sub := range subscriptions {
			var statusIcon string
			if sub.IsActive && !sub.IsExpired() {
				statusIcon = "🟢"
			} else {
				statusIcon = "⚪"
			}
			viewCallbackData := fmt.Sprintf("view_subscription_%s", sub.ID)
			viewButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", statusIcon, sub.GetDisplayName()),
				viewCallbackData)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(viewButton))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetSubscriptionDetailKeyboardOld(subscriptionID string) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать", fmt.Sprintf("rename_subscription_%s", subscriptionID)),
			tgbotapi.NewInlineKeyboardButtonData("📈 Продлить", fmt.Sprintf("extend_subscription_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_subscription_%s", subscriptionID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "my_subscriptions"),
		),
	)
}
func GetExtendSubscriptionKeyboard(subscriptionID string, plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 +%s - %.0f₽", FormatDuration(plan.Days), plan.Price)
			callbackData := fmt.Sprintf("extend_plan_%s_sub_%s", plan.ID, subscriptionID)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("view_subscription_%s", subscriptionID)),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetCreateSubscriptionKeyboard(plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 %s - %.0f₽ (%s)", plan.Name, plan.Price, FormatDuration(plan.Days))
			callbackData := fmt.Sprintf("create_plan_%s", plan.ID)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetKeysKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔑 WireGuard", "create_wireguard"),
			tgbotapi.NewInlineKeyboardButtonData("🔑 Shadowsocks", "create_shadowsocks"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои конфиги", "my_configs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_profile"),
		),
	)
}
func GetReferralsKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "referral_stats"),
			tgbotapi.NewInlineKeyboardButtonData("👥 Мои рефералы", "my_referrals"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Рейтинг", "referral_ranking"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Моя ссылка", "my_referral_link"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
		),
	)
}
func GetWelcomeText(firstName string, hasTrialUsed bool) string {
	greeting := fmt.Sprintf("👋 Привет, %s!\n\n", firstName)
	text := greeting + `🎉 Добро пожаловать в VPN Bot!
Этот бот поможет вам:
• 🔐 Создавать безопасные VPN подключения
• 💳 Управлять подписками
• 👥 Приглашать друзей и получать бонусы
• 📊 Отслеживать статистику использования
`
	if !hasTrialUsed {
		text += "\n🎁 Нажмите кнопку ниже, чтобы получить пробный доступ на 3 дня бесплатно!"
	} else {
		text += "\n💰 Выберите подходящий тариф для покупки подписки."
	}

	return text
}
func GetMainMenuText() string {

	return `📱 Главное меню
Выберите нужное действие:`
}
func GetMainMenuWithProfileText(user *core.User, subscriptions []*core.Subscription) string {
	text := "👤 *Ваш профиль*\n\n"
	text += fmt.Sprintf("🆔 *ID:* `%d`\n", user.TelegramID)
	text += fmt.Sprintf("🙋‍♂️ *Имя:* `%s`\n", user.GetDisplayName())
	langDisplay := user.LanguageCode
	if langDisplay == "" {
		langDisplay = "ru"
	}
	text += fmt.Sprintf("🌐 *Язык интерфейса:* `%s`\n", langDisplay)
	activeSubscriptions := make([]*core.Subscription, 0)
	for _, sub := range subscriptions {
		if sub.IsActive && !sub.IsExpired() {
			activeSubscriptions = append(activeSubscriptions, sub)
		}
	}
	if len(activeSubscriptions) == 0 {
		text += "💫 *Статус:* 🆓 *Бесплатный*\n"
	} else {
		text += "💫 *Статус:* ⭐️ *Premium*\n"
		for _, sub := range activeSubscriptions {
			subName := EscapeMarkdownV2(sub.GetDisplayName())
			subUntil := sub.EndDate.Format("02\\.01\\.2006")
			text += fmt.Sprintf("⏰ *Подписка активна до:* `%s` \\(%s\\)\n", subUntil, subName)
		}
	}
	text += "\n━━━━━━━━━━━━━━━\n"
	text += "📍 *Главное меню*\n"
	text += "Выберите нужное действие ниже ⤵️"

	return text
}
func GetInstructionText() string {

	return `📖 Инструкция по использованию
🔹 Как начать:
1. Получите пробный доступ или купите подписку
2. После оплаты вы автоматически получите VPN конфигурацию
3. Скачайте приложение для вашей платформы
4. Добавьте конфигурацию в приложение
🔹 Рекомендуемые приложения:
📱 iOS:
• Shadowrocket (платно)
• FoXray (бесплатно)
🤖 Android:
• v2rayNG (бесплатно)
• NekoBox (бесплатно)
💻 Windows:
• v2rayN (бесплатно)
• Hiddify (бесплатно)
🍎 macOS:
• V2Box (бесплатно)
• FoXray (бесплатно)
🔹 Как подключиться:
1. Откройте приложение
2. Нажмите "Добавить конфигурацию"
3. Отсканируйте QR-код или вставьте ссылку
4. Нажмите "Подключиться"
💡 Если возникли проблемы - обратитесь в поддержку!`
}

func GetInstructionWithConnectionText(subscriptionID string) string {
	connectionURL := fmt.Sprintf("https://3xui.com/connect/%s", subscriptionID)

	return fmt.Sprintf(`📖 Инструкция по подключению

🔗 Ваша ссылка для подключения:
`+"`%s`"+`

🔹 Как начать:
1. Скачайте приложение для вашей платформы
2. Откройте приложение
3. Нажмите "Добавить конфигурацию"
4. Отсканируйте QR-код или вставьте ссылку выше
5. Нажмите "Подключиться"

🔹 Рекомендуемые приложения:

📱 iOS:
• Shadowrocket (платно)
• FoXray (бесплатно)

🤖 Android:
• v2rayNG (бесплатно)
• NekoBox (бесплатно)

💻 Windows:
• v2rayN (бесплатно)
• Hiddify (бесплатно)

🍎 macOS:
• V2Box (бесплатно)
• FoXray (бесплатно)

💡 Если возникли проблемы - обратитесь в поддержку!`, connectionURL)
}
func GetProfileText(user *core.User, isPremium bool, statusText, subUntilText string) string {
	text := "👤 Ваш профиль\n\n"
	text += fmt.Sprintf("🆔 ID: %d\n", user.TelegramID)
	text += fmt.Sprintf("👋 Имя: %s\n", user.GetDisplayName())
	text += fmt.Sprintf("🌐 Язык: %s\n", user.LanguageCode)
	text += fmt.Sprintf("📊 Статус: %s\n", statusText)
	if isPremium && subUntilText != "" {
		text += fmt.Sprintf("⏰ Подписка до: %s\n", subUntilText)
	}

	return text
}
func GetPricingText(plans []*core.Plan) string {

	return "💳 Выберите тарифный план для создания нового ключа:"
}
func GetPaymentMethodText(plan *core.Plan) string {

	return fmt.Sprintf(`💳 Оплата подписки
📦 План: %s
💵 Сумма: %.0f₽
⏰ Длительность: %s
Выберите способ оплаты:`, plan.Name, plan.Price, FormatDuration(plan.Days))
}
func GetSubscriptionsText(subscriptions []*core.Subscription) string {
	text := "*🔑 Список ваших подписок:*\n\n"
	if len(subscriptions) == 0 {
		text += "У вас пока нет подписок\\.\n\n"
		text += "💡 Создайте подписку, чтобы получить доступ к VPN сервисам\\!"
	} else {
		for i, sub := range subscriptions {
			displayName := EscapeMarkdownV2(sub.GetDisplayName())
			dateStr := EscapeMarkdownV2(sub.EndDate.Format("02.01.06, 15:04"))
			var statusIcon string
			if sub.IsActive && !sub.IsExpired() {
				statusIcon = "🟢"
			} else {
				statusIcon = "⚪"
			}
			text += fmt.Sprintf("> %s • %s \\(до %s\\) »\n", statusIcon, displayName, dateStr)
			if i < len(subscriptions)-1 {
				text += "\n"
			}
		}
		text += "\n_Нажмите на подписку, чтобы просмотреть детали и управлять ею\\._"
	}

	return text
}
func EscapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(text)
}
func GetProfileTextWithQuotes(user *core.User, subscriptionCount int) string {
	displayName := EscapeMarkdownV2(user.GetDisplayName())
	text := fmt.Sprintf("Профиль: %s\n\n", displayName)
	text += fmt.Sprintf("> \\-\\- ID: %d\n", user.TelegramID)
	text += fmt.Sprintf("> \\-\\- К\\-во подписок: %d »\n\n", subscriptionCount)
	text += "👉 Наш канал 👈\n"
	text += "👉 Поддержка 👈\n\n"
	text += "> _Приглашай друзей, получай больше\\!_ »"

	return text
}
func GetProfileKeyboardWithQuotes() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("👉 Наш канал 👈", "https://t.me/your_channel"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("👉 Поддержка 👈", "https://t.me/your_support"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetSubscriptionDetailTextOld(sub *core.Subscription) string {
	text := "📋 Детали подписки\n\n"
	text += fmt.Sprintf("📝 Название: %s\n", sub.GetDisplayName())
	text += fmt.Sprintf("📊 Статус: %s\n", sub.GetStatusText())
	text += fmt.Sprintf("📅 Начало: %s\n", sub.StartDate.Format("02.01.2006 15:04"))
	text += fmt.Sprintf("📅 Окончание: %s\n", sub.EndDate.Format("02.01.2006 15:04"))
	if sub.IsActive && !sub.IsExpired() {
		text += fmt.Sprintf("⏰ Осталось: %d дней\n", sub.DaysRemaining())
	}
	text += fmt.Sprintf("🆔 ID: %s\n", sub.ID)
	text += fmt.Sprintf("📋 План: %s\n", sub.PlanID)

	return text
}
func GetCreateSubscriptionText() string {

	return `➕ Создание новой подписки
Выберите подходящий тарифный план для вашей подписки.
💡 После создания подписки вы получите доступ к VPN сервисам!`
}
func GetRenameSubscriptionText(sub *core.Subscription) string {

	return fmt.Sprintf(`✏️ Переименование подписки
Текущее название: %s
Введите новое название для подписки:`, sub.GetDisplayName())
}
func GetExtendSubscriptionText(sub *core.Subscription) string {
	text := fmt.Sprintf(`📈 Продление подписки
Подписка: %s
Текущее окончание: %s
Выберите период продления:`,
		sub.GetDisplayName(),
		sub.EndDate.Format("02.01.2006"))

	return text
}
func GetDeleteSubscriptionText(sub *core.Subscription) string {

	return fmt.Sprintf(`🗑️ Удаление подписки
Подписка: %s
Статус: %s
⚠️ Внимание! Это действие нельзя отменить.
Все связанные с этой подпиской данные будут удалены.
Вы уверены, что хотите удалить эту подписку?`,
		sub.GetDisplayName(),
		sub.GetStatusText())
}
func GetKeysText() string {

	return `🔑 Управление ключами
Здесь вы можете:
• Создавать новые VPN конфигурации
• Просматривать существующие конфиги
• Управлять доступом к серверам
Выберите тип конфигурации:`
}
func GetReferralsText() string {

	return `👥 Реферальная программа
Приглашайте друзей и получайте бонусы!
🎁 За каждого приглашенного друга вы получите:
• 7 дней бесплатной подписки
• Доступ к дополнительным функциям
📊 Отслеживайте статистику приглашений и заработанные бонусы.`
}
func GetSupportText() string {

	return `💬 Поддержка
Если у вас возникли вопросы или проблемы, обратитесь к нашей поддержке:
📧 Email: support@3xui.com
💬 Telegram: @3xui_support
🌐 Сайт: https:
⏰ Время ответа: до 24 часов`
}
func GetReferralRankingText() string {

	return `🏆 Рейтинг рефералов
Здесь можно увидеть топ людей, которые пригласили наибольшее количество рефералов в сервис.
Твоё место в рейтинге:
Ты еще не приглашал пользователей в проект.
🏆 Топ-5 пригласивших:
1. 57956***** - 156 чел.
2. 80000***** - 105 чел.
3. 52587***** - 12 чел.
4. 63999***** - 7 чел.
5. 10149***** - 6 чел.`
}
func GetReferralRankingKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
		),
	)
}
func FormatPrice(price float64) string {

	return fmt.Sprintf("%.0f₽", price)
}
func FormatDuration(days int) string {
	if days >= 365 {
		years := days / 365
		remainingDays := days % 365
		if remainingDays == 0 {
			if years == 1 {

				return "1 год"
			}

			return fmt.Sprintf("%d лет", years)
		}
		if years == 1 {

			return fmt.Sprintf("1 год %d дней", remainingDays)
		}

		return fmt.Sprintf("%d лет %d дней", years, remainingDays)
	} else if days >= 30 {
		months := days / 30
		remainingDays := days % 30
		if remainingDays == 0 {
			if months == 1 {

				return "1 месяц"
			} else if months >= 2 && months <= 4 {

				return fmt.Sprintf("%d месяца", months)
			}

			return fmt.Sprintf("%d месяцев", months)
		}
		if months == 1 {

			return fmt.Sprintf("1 месяц %d дней", remainingDays)
		} else if months >= 2 && months <= 4 {

			return fmt.Sprintf("%d месяца %d дней", months, remainingDays)
		}

		return fmt.Sprintf("%d месяцев %d дней", months, remainingDays)
	} else if days >= 7 {
		weeks := days / 7
		remainingDays := days % 7
		if remainingDays == 0 {
			if weeks == 1 {

				return "1 неделя"
			} else if weeks >= 2 && weeks <= 4 {

				return fmt.Sprintf("%d недели", weeks)
			}

			return fmt.Sprintf("%d недель", weeks)
		}
		if weeks == 1 {

			return fmt.Sprintf("1 неделя %d дней", remainingDays)
		} else if weeks >= 2 && weeks <= 4 {

			return fmt.Sprintf("%d недели %d дней", weeks, remainingDays)
		}

		return fmt.Sprintf("%d недель %d дней", weeks, remainingDays)
	} else if days == 1 {

		return "1 день"
	} else if days >= 2 && days <= 4 {

		return fmt.Sprintf("%d дня", days)
	}

	return fmt.Sprintf("%d дней", days)
}
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {

		return s
	}

	return s[:maxLen-3] + "..."
}
func EscapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}
func GetSubscriptionDetailText(subscription *core.Subscription, plan *core.Plan, vpnConfigs []*core.VPNConnection) string {
	var text strings.Builder
	text.WriteString("*📋 Детали подписки*\n\n")
	planName := EscapeMarkdownV2(plan.Name)
	text.WriteString(fmt.Sprintf("> 📦 *План:* %s\n", planName))
	text.WriteString(fmt.Sprintf("> 💰 *Цена:* %.0f₽\n", plan.Price))
	text.WriteString(fmt.Sprintf("> ⏰ *Длительность:* %d дней »\n\n", plan.Days))
	if subscription.IsActive {
		text.WriteString("✅ *Статус:* Активна\n")
		endDate := EscapeMarkdownV2(subscription.EndDate.Format("02.01.06, 15:04"))
		text.WriteString(fmt.Sprintf("📅 *Активна до:* %s\n", endDate))
	} else {
		text.WriteString("❌ *Статус:* Неактивна\n")
	}
	startDate := EscapeMarkdownV2(subscription.StartDate.Format("02.01.06"))
	text.WriteString(fmt.Sprintf("📅 *Создана:* %s\n\n", startDate))
	if subscription.IsActive {
		connectionURL := fmt.Sprintf("https://3xui.com/connect/%s", subscription.ID)
		text.WriteString("*🔗 URL подключения:*\n")
		text.WriteString(fmt.Sprintf("`%s`", connectionURL))
	}

	return text.String()
}
func GetSubscriptionDetailKeyboard(subscription *core.Subscription, vpnConfigs []*core.VPNConnection) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	if subscription.IsActive {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Инструкция по подключению", fmt.Sprintf("connection_guide_%s", subscription.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать", fmt.Sprintf("rename_subscription_%s", subscription.ID)),
	))
	if subscription.IsActive {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Продлить подписку", fmt.Sprintf("extend_subscription_%s", subscription.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к подпискам", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetVPNConfigDetailText(config *core.VPNConnection) string {
	var text strings.Builder
	text.WriteString("🔑 Детали VPN конфигурации\n\n")
	text.WriteString(fmt.Sprintf("📝 Название: %s\n", config.Name))
	text.WriteString(fmt.Sprintf("🔐 Username: %s\n", config.MarzbanUsername))
	text.WriteString(fmt.Sprintf("📊 Статус: %s\n", config.Status))
	if config.ExpireAt != nil {
		text.WriteString(fmt.Sprintf("⏰ Истекает: %s\n", config.ExpireAt.Format("02.01.2006 15:04")))
		remaining := time.Until(*config.ExpireAt)
		if remaining > 0 {
			days := int(remaining.Hours() / 24)
			hours := int(remaining.Hours()) % 24
			text.WriteString(fmt.Sprintf("⏳ Осталось: %d дней %d часов\n", days, hours))
		}
	}
	if config.DataLimitBytes != nil && *config.DataLimitBytes > 0 {
		dataLimitGB := float64(*config.DataLimitBytes) / (1024 * 1024 * 1024)
		text.WriteString(fmt.Sprintf("💾 Лимит трафика: %.1f GB\n", dataLimitGB))
		if config.DataUsedBytes != nil {
			dataUsedGB := float64(*config.DataUsedBytes) / (1024 * 1024 * 1024)
			dataRemainingGB := dataLimitGB - dataUsedGB
			usagePercent := (dataUsedGB / dataLimitGB) * 100
			text.WriteString(fmt.Sprintf("📊 Использовано: %.2f GB (%.1f%%)\n", dataUsedGB, usagePercent))
			text.WriteString(fmt.Sprintf("📉 Осталось: %.2f GB\n", dataRemainingGB))
		}
	}
	text.WriteString(fmt.Sprintf("\n📅 Создано: %s\n", config.CreatedAt.Format("02.01.2006 15:04")))
	text.WriteString(fmt.Sprintf("🔄 Обновлено: %s\n", config.UpdatedAt.Format("02.01.2006 15:04")))
	text.WriteString("\nНажмите \"Получить ключ\" для получения конфигурации.")

	return text.String()
}
func GetVPNConfigDetailKeyboard(config *core.VPNConnection) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔑 Получить ключ", fmt.Sprintf("get_vpn_key_%s", config.ID)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать", fmt.Sprintf("rename_config_%s", config.ID)),
	))
	if !config.IsActive || config.IsExpired() {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_config_%s", config.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetUnknownCommandText() string {

	return `🤖 Я пока не умею отвечать на такие сообщения
❓ У вас вопрос или возникли сложности?
Свяжитесь с нашей поддержкой — мы поможем как можно скорее.
🔒 Управление подпиской
Всё, что касается вашего VPN — тарифы, продления, подключение — доступно в личном кабинете 👇`
}
func GetUnknownCommandKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("💬 Поддержка", "https://t.me/your_support_chat"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func GetCancelKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "my_subscriptions"),
		),
	)
}
func GetBackToPricingKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к тарифам", "open_pricing"),
		),
	)
}

func GetBackToMenuKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
		),
	)
}

func GetBackToSubscriptionsKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 К моим подпискам", "my_subscriptions"),
		),
	)
}
