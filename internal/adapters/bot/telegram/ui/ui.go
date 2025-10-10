package ui

import (
	"fmt"
	"strings"

	"3xui-bot/internal/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================================
// КЛАВИАТУРЫ
// ============================================================================

// GetWelcomeKeyboard возвращает клавиатуру приветствия
func GetWelcomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎉 Получить пробный доступ", "get_trial"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Тарифы", "open_pricing"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Профиль", "open_profile"),
		),
	)
}

// GetProfileKeyboard возвращает клавиатуру профиля
func GetProfileKeyboard(isPremium bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопки подписки
	if isPremium {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
			tgbotapi.NewInlineKeyboardButtonData("🔑 Мои ключи/конфиги", "open_keys"),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Мои подписки", "my_subscriptions"),
			tgbotapi.NewInlineKeyboardButtonData("🔑 Мои ключи/конфиги", "open_keys"),
		))
	}

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
			callbackData := fmt.Sprintf("plan_%s", plan.ID)

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

// GetSubscriptionsKeyboard возвращает клавиатуру со списком подписок
func GetSubscriptionsKeyboard(subscriptions []*core.Subscription) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	if len(subscriptions) == 0 {
		// Если подписок нет, показываем кнопку создания
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
		))
	} else {
		// Создаем кнопки для каждой подписки
		for _, sub := range subscriptions {
			buttonText := fmt.Sprintf("📋 %s", sub.GetDisplayName())
			callbackData := fmt.Sprintf("view_subscription_%s", sub.ID)

			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
			))
		}

		// Кнопка создания новой подписки
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать подписку", "create_subscription"),
		))
	}

	// Кнопка назад
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_profile"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// GetSubscriptionDetailKeyboard возвращает клавиатуру деталей подписки
func GetSubscriptionDetailKeyboard(subscriptionID string) tgbotapi.InlineKeyboardMarkup {
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
			tgbotapi.NewInlineKeyboardButtonData("🔗 Моя ссылка", "my_referral_link"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "open_profile"),
		),
	)
}

// ============================================================================
// ТЕКСТЫ
// ============================================================================

// GetWelcomeText возвращает текст приветствия
func GetWelcomeText() string {
	return `🎉 Добро пожаловать в 3xui-bot!

Этот бот поможет вам:
• 🔐 Создавать VPN конфигурации
• 💳 Управлять подписками
• 👥 Приглашать друзей
• 📊 Отслеживать статистику

Нажмите кнопку ниже, чтобы получить пробный доступ на 3 дня!`
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
	text := "💰 Тарифные планы\n\n"

	for _, plan := range plans {
		if plan.IsActive {
			text += fmt.Sprintf("📦 %s\n", plan.Name)
			text += fmt.Sprintf("   💵 Цена: %.0f₽\n", plan.Price)
			text += fmt.Sprintf("   ⏰ Длительность: %s\n", FormatDuration(plan.Days))
			text += fmt.Sprintf("   💰 Цена за день: %.2f₽\n", plan.GetPricePerDay())
			if plan.Description != "" {
				text += fmt.Sprintf("   📝 %s\n", plan.Description)
			}

			// Показываем скидку если есть
			if discount := plan.GetDiscount(); discount > 0 {
				text += fmt.Sprintf("   🎯 Скидка: %.0f%%\n", discount)
			}
			text += "\n"
		}
	}

	text += "💡 Выберите подходящий план для покупки:"

	return text
}

// GetSubscriptionsText возвращает текст со списком подписок
func GetSubscriptionsText(subscriptions []*core.Subscription) string {
	text := "💳 Ваши подписки\n\n"

	if len(subscriptions) == 0 {
		text += "У вас пока нет активных подписок.\n\n"
		text += "💡 Создайте подписку, чтобы получить доступ к VPN сервисам!"
	} else {
		for i, sub := range subscriptions {
			text += fmt.Sprintf("%d. 📋 %s\n", i+1, sub.GetDisplayName())
			text += fmt.Sprintf("   📊 Статус: %s\n", sub.GetStatusText())
			text += fmt.Sprintf("   📅 Начало: %s\n", sub.StartDate.Format("02.01.2006"))
			text += fmt.Sprintf("   📅 Окончание: %s\n", sub.EndDate.Format("02.01.2006"))

			if sub.IsActive && !sub.IsExpired() {
				text += fmt.Sprintf("   ⏰ Осталось: %d дней\n", sub.DaysRemaining())
			}
			text += "\n"
		}

		text += "💡 Нажмите на подписку для управления или создайте новую!"
	}

	return text
}

// GetSubscriptionDetailText возвращает текст деталей подписки
func GetSubscriptionDetailText(sub *core.Subscription) string {
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
