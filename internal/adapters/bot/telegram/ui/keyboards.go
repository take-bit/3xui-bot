package ui

import (
	"fmt"
	"strings"
	"time"

	"3xui-bot/internal/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================================
// КЛАВИАТУРЫ
// ============================================================================

// GetWelcomeKeyboard возвращает клавиатуру приветствия для нового пользователя
func GetWelcomeKeyboard(hasTrialUsed bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Если пользователь еще не использовал триал
	if !hasTrialUsed {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉 Получить пробный доступ", "get_trial"),
		))
		// Показываем кнопку покупки только если нет пробного доступа
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	} else {
		// Если пробный доступ уже использован, показываем "Мои подписки"
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔑 Мои подписки", "my_subscriptions"),
		))
		// И кнопку покупки дополнительной подписки
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetMainMenuKeyboard возвращает главное меню (простое, для новых пользователей)
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

// GetMainMenuWithProfileKeyboard возвращает объединенную клавиатуру профиля и главного меню
func GetMainMenuWithProfileKeyboard(isPremium bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка подписок (ключи будут внутри подписок)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
	))

	// Основные кнопки меню
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

// GetProfileKeyboard возвращает клавиатуру профиля (устарела, используем объединенное меню)
func GetProfileKeyboard(isPremium bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка подписок (ключи будут внутри подписок)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
	))

	// Кнопки рефералов
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👥 Реферальная программа", "open_referrals"),
		tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "open_support"),
	))

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetPricingKeyboard возвращает клавиатуру с тарифами
func GetPricingKeyboard(plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для каждого плана
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 %s - %.0f₽ (%s)", plan.Name, plan.Price, FormatDuration(plan.Days))
			callbackData := fmt.Sprintf("select_plan_%s", plan.ID)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetPaymentMethodKeyboard возвращает клавиатуру выбора способа оплаты
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

// GetSubscriptionsKeyboard возвращает клавиатуру со списком всех подписок
func GetSubscriptionsKeyboard(subscriptions []*core.Subscription) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	if len(subscriptions) == 0 {
		// Если подписок нет, показываем кнопку покупки
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	} else {
		// Создаем кнопки для каждой подписки
		for _, sub := range subscriptions {
			// Определяем статус подписки для кнопки
			var statusIcon string
			if sub.IsActive && !sub.IsExpired() {
				statusIcon = "🟢" // Активна
			} else {
				statusIcon = "⚪" // Истекла
			}

			// Кнопка просмотра подписки с статус-индикатором
			viewCallbackData := fmt.Sprintf("view_subscription_%s", sub.ID)
			viewButton := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", statusIcon, sub.GetDisplayName()),
				viewCallbackData)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(viewButton))
		}

		// Кнопка покупки новой подписки
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Купить подписку", "open_pricing"),
		))
	}

	// Кнопка назад (изменили на "Личный кабинет")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetSubscriptionDetailKeyboardOld возвращает клавиатуру деталей подписки (старая версия)
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

// GetExtendSubscriptionKeyboard возвращает клавиатуру продления подписки
func GetExtendSubscriptionKeyboard(subscriptionID string, plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для каждого плана продления
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 +%s - %.0f₽", FormatDuration(plan.Days), plan.Price)
			callbackData := fmt.Sprintf("extend_plan_%s_sub_%s", plan.ID, subscriptionID)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("view_subscription_%s", subscriptionID)),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetCreateSubscriptionKeyboard возвращает клавиатуру создания подписки
func GetCreateSubscriptionKeyboard(plans []*core.Plan) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Создаем кнопки для каждого плана
	for _, plan := range plans {
		if plan.IsActive {
			buttonText := fmt.Sprintf("📦 %s - %.0f₽ (%s)", plan.Name, plan.Price, FormatDuration(plan.Days))
			callbackData := fmt.Sprintf("create_plan_%s", plan.ID)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}
	}

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetKeysKeyboard возвращает клавиатуру управления ключами
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

// GetReferralsKeyboard возвращает клавиатуру реферальной программы
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

// ============================================================================
// ТЕКСТЫ
// ============================================================================

// GetWelcomeText возвращает текст приветствия
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

// GetMainMenuText возвращает текст главного меню
func GetMainMenuText() string {
	return `📱 Главное меню

Выберите нужное действие:`
}

// GetMainMenuWithProfileText возвращает объединенный текст профиля и главного меню
func GetMainMenuWithProfileText(user *core.User, isPremium bool, statusText, subUntilText string) string {
	text := "👤 Ваш профиль\n\n"
	text += fmt.Sprintf("🆔 ID: %d\n", user.TelegramID)
	text += fmt.Sprintf("👋 Имя: %s\n", user.GetDisplayName())
	text += fmt.Sprintf("🌐 Язык: %s\n", user.LanguageCode)
	text += fmt.Sprintf("📊 Статус: %s\n", statusText)

	if isPremium && subUntilText != "" {
		text += fmt.Sprintf("⏰ Подписка до: %s\n", subUntilText)
	}

	text += "\n📱 Главное меню\n"
	text += "Выберите нужное действие:"

	return text
}

// GetInstructionText возвращает текст инструкции
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

// GetProfileText возвращает текст профиля
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

// GetPricingText возвращает текст с тарифами
func GetPricingText(plans []*core.Plan) string {
	return "💳 Выберите тарифный план для создания нового ключа:"
}

// GetPaymentMethodText возвращает текст выбора способа оплаты
func GetPaymentMethodText(plan *core.Plan) string {
	return fmt.Sprintf(`💳 Оплата подписки

📦 План: %s
💵 Сумма: %.0f₽
⏰ Длительность: %s

Выберите способ оплаты:`, plan.Name, plan.Price, FormatDuration(plan.Days))
}

// GetSubscriptionsText возвращает текст со списком всех подписок в MarkdownV2 формате
func GetSubscriptionsText(subscriptions []*core.Subscription) string {
	text := "*🔑 Список ваших подписок:*\n\n"

	if len(subscriptions) == 0 {
		text += "У вас пока нет подписок\\.\n\n"
		text += "💡 Создайте подписку, чтобы получить доступ к VPN сервисам\\!"
	} else {
		for i, sub := range subscriptions {
			// Экранируем специальные символы для MarkdownV2
			displayName := EscapeMarkdownV2(sub.GetDisplayName())
			dateStr := EscapeMarkdownV2(sub.EndDate.Format("02.01.06, 15:04"))

			// Определяем статус подписки
			var statusIcon string
			if sub.IsActive && !sub.IsExpired() {
				statusIcon = "🟢" // Активна
			} else {
				statusIcon = "⚪" // Истекла
			}

			// Отображаем каждую подписку в отдельной цитате с статус-индикатором
			text += fmt.Sprintf("> %s • %s \\(до %s\\) »\n", statusIcon, displayName, dateStr)

			// Добавляем пустую строку между подписками (кроме последней)
			if i < len(subscriptions)-1 {
				text += "\n"
			}
		}
		text += "\n_Нажмите на подписку, чтобы просмотреть детали и управлять ею\\._"
	}

	return text
}

// EscapeMarkdownV2 экранирует специальные символы для MarkdownV2
func EscapeMarkdownV2(text string) string {
	// Символы, которые нужно экранировать в MarkdownV2
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

// GetProfileTextWithQuotes возвращает текст профиля с цитатами в MarkdownV2 формате
func GetProfileTextWithQuotes(user *core.User, subscriptionCount int) string {
	// Экранируем имя пользователя
	displayName := EscapeMarkdownV2(user.GetDisplayName())

	text := fmt.Sprintf("Профиль: %s\n\n", displayName)

	// Первая цитата с информацией профиля
	text += fmt.Sprintf("> \\-\\- ID: %d\n", user.TelegramID)
	text += fmt.Sprintf("> \\-\\- К\\-во подписок: %d »\n\n", subscriptionCount)

	// Кнопки с эмодзи
	text += "👉 Наш канал 👈\n"
	text += "👉 Поддержка 👈\n\n"

	// Вторая цитата с приглашением друзей
	text += "> _Приглашай друзей, получай больше\\!_ »"

	return text
}

// GetProfileKeyboardWithQuotes возвращает клавиатуру для профиля с цитатами
func GetProfileKeyboardWithQuotes() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка "Наш канал"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("👉 Наш канал 👈", "https://t.me/your_channel"),
	))

	// Кнопка "Поддержка"
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("👉 Поддержка 👈", "https://t.me/your_support"),
	))

	// Кнопка назад в главное меню
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetSubscriptionDetailTextOld возвращает текст деталей подписки (старая версия, оставлена для совместимости)
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

// GetCreateSubscriptionText возвращает текст создания подписки
func GetCreateSubscriptionText() string {
	return `➕ Создание новой подписки

Выберите подходящий тарифный план для вашей подписки.

💡 После создания подписки вы получите доступ к VPN сервисам!`
}

// GetRenameSubscriptionText возвращает текст переименования подписки
func GetRenameSubscriptionText(sub *core.Subscription) string {
	return fmt.Sprintf(`✏️ Переименование подписки

Текущее название: %s

Введите новое название для подписки:`, sub.GetDisplayName())
}

// GetExtendSubscriptionText возвращает текст продления подписки
func GetExtendSubscriptionText(sub *core.Subscription) string {
	text := fmt.Sprintf(`📈 Продление подписки

Подписка: %s
Текущее окончание: %s

Выберите период продления:`,
		sub.GetDisplayName(),
		sub.EndDate.Format("02.01.2006"))

	return text
}

// GetDeleteSubscriptionText возвращает текст удаления подписки
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

// GetKeysText возвращает текст управления ключами
func GetKeysText() string {
	return `🔑 Управление ключами

Здесь вы можете:
• Создавать новые VPN конфигурации
• Просматривать существующие конфиги
• Управлять доступом к серверам

Выберите тип конфигурации:`
}

// GetReferralsText возвращает текст реферальной программы
func GetReferralsText() string {
	return `👥 Реферальная программа

Приглашайте друзей и получайте бонусы!

🎁 За каждого приглашенного друга вы получите:
• 7 дней бесплатной подписки
• Доступ к дополнительным функциям

📊 Отслеживайте статистику приглашений и заработанные бонусы.`
}

// GetSupportText возвращает текст поддержки
func GetSupportText() string {
	return `💬 Поддержка

Если у вас возникли вопросы или проблемы, обратитесь к нашей поддержке:

📧 Email: support@3xui.com
💬 Telegram: @3xui_support
🌐 Сайт: https://3xui.com

⏰ Время ответа: до 24 часов`
}

// GetReferralRankingText возвращает текст реферального рейтинга
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

// GetReferralRankingKeyboard возвращает клавиатуру реферального рейтинга
func GetReferralRankingKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_menu"),
		),
	)
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================================

// FormatPrice форматирует цену
func FormatPrice(price float64) string {
	return fmt.Sprintf("%.0f₽", price)
}

// FormatDuration форматирует длительность
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

// TruncateString обрезает строку до указанной длины
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// EscapeMarkdown экранирует специальные символы Markdown
func EscapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}

// GetSubscriptionDetailText возвращает детальную информацию о подписке с VPN конфигами в MarkdownV2
func GetSubscriptionDetailText(subscription *core.Subscription, plan *core.Plan, vpnConfigs []*core.VPNConnection) string {
	var text strings.Builder

	// Заголовок
	text.WriteString("*📋 Детали подписки*\n\n")

	// Информация о плане в цитате
	planName := EscapeMarkdownV2(plan.Name)
	text.WriteString(fmt.Sprintf("> 📦 *План:* %s\n", planName))
	text.WriteString(fmt.Sprintf("> 💰 *Цена:* %.0f₽\n", plan.Price))
	text.WriteString(fmt.Sprintf("> ⏰ *Длительность:* %d дней »\n\n", plan.Days))

	// Статус подписки
	if subscription.IsActive {
		text.WriteString("✅ *Статус:* Активна\n")
		endDate := EscapeMarkdownV2(subscription.EndDate.Format("02.01.06, 15:04"))
		text.WriteString(fmt.Sprintf("📅 *Активна до:* %s\n", endDate))
	} else {
		text.WriteString("❌ *Статус:* Неактивна\n")
	}

	startDate := EscapeMarkdownV2(subscription.StartDate.Format("02.01.06"))
	text.WriteString(fmt.Sprintf("📅 *Создана:* %s\n\n", startDate))

	// URL подключения (если подписка активна)
	if subscription.IsActive {
		// Генерируем URL подключения (пока простой)
		connectionURL := fmt.Sprintf("https://3xui.com/connect/%s", subscription.ID)
		text.WriteString("*🔗 URL подключения:*\n")
		text.WriteString(fmt.Sprintf("`%s`", connectionURL))
	}

	return text.String()
}

// GetSubscriptionDetailKeyboard возвращает клавиатуру для детальной информации о подписке
func GetSubscriptionDetailKeyboard(subscription *core.Subscription, vpnConfigs []*core.VPNConnection) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка инструкции по подключению (если подписка активна)
	if subscription.IsActive {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Инструкция по подключению", fmt.Sprintf("connection_guide_%s", subscription.ID)),
		))
	}

	// Кнопка редактирования подписки
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать", fmt.Sprintf("rename_subscription_%s", subscription.ID)),
	))

	// Кнопка продления подписки (если она активна)
	if subscription.IsActive {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Продлить подписку", fmt.Sprintf("extend_subscription_%s", subscription.ID)),
		))
	}

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к подпискам", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetVPNConfigDetailText возвращает детальную информацию о VPN конфигурации
func GetVPNConfigDetailText(config *core.VPNConnection) string {
	var text strings.Builder

	text.WriteString("🔑 Детали VPN конфигурации\n\n")
	text.WriteString(fmt.Sprintf("📝 Название: %s\n", config.Name))
	text.WriteString(fmt.Sprintf("🔐 Username: %s\n", config.MarzbanUsername))
	text.WriteString(fmt.Sprintf("📊 Статус: %s\n", config.Status))

	if config.ExpireAt != nil {
		text.WriteString(fmt.Sprintf("⏰ Истекает: %s\n", config.ExpireAt.Format("02.01.2006 15:04")))

		// Считаем оставшееся время
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

// GetVPNConfigDetailKeyboard возвращает клавиатуру для детальной информации о VPN конфигурации
func GetVPNConfigDetailKeyboard(config *core.VPNConnection) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка получения ключа/конфигурации
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔑 Получить ключ", fmt.Sprintf("get_vpn_key_%s", config.ID)),
	))

	// Кнопка переименования
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Переименовать", fmt.Sprintf("rename_config_%s", config.ID)),
	))

	// Кнопка удаления (если конфиг неактивен или истек)
	if !config.IsActive || config.IsExpired() {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_config_%s", config.ID)),
		))
	}

	// Кнопка назад (нужно знать subscription_id - для упрощения возвращаемся к подпискам)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "my_subscriptions"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetUnknownCommandText возвращает текст для неизвестной команды
func GetUnknownCommandText() string {
	return `🤖 Я пока не умею отвечать на такие сообщения

❓ У вас вопрос или возникли сложности?
Свяжитесь с нашей поддержкой — мы поможем как можно скорее.

🔒 Управление подпиской
Всё, что касается вашего VPN — тарифы, продления, подключение — доступно в личном кабинете 👇`
}

// GetUnknownCommandKeyboard возвращает клавиатуру для неизвестной команды
func GetUnknownCommandKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопка поддержки
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("💬 Поддержка", "https://t.me/your_support_chat"),
	))

	// Кнопка личного кабинета (возврат в главное меню)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "open_menu"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetCancelKeyboard возвращает клавиатуру с кнопкой отмены
func GetCancelKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "my_subscriptions"),
		),
	)
}

// GetBackToPricingKeyboard возвращает клавиатуру для возврата к тарифам
func GetBackToPricingKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к тарифам", "open_pricing"),
		),
	)
}
