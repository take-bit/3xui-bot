package ui

// Основные команды
const (
	CommandStart = "start"
	CommandHelp  = "help"
)

// Основные callback'ы
const (
	CallbackGetTrial           = "get_trial"
	CallbackOpenMenu           = "open_menu"
	CallbackOpenProfile        = "open_profile"
	CallbackOpenPricing        = "open_pricing"
	CallbackMySubscriptions    = "my_subscriptions"
	CallbackCreateSubscription = "create_subscription"
	CallbackOpenKeys           = "open_keys"
	CallbackOpenReferrals      = "open_referrals"
	CallbackOpenSupport        = "open_support"
	CallbackShowInstruction    = "show_instruction"
	CallbackMyConfigs          = "my_configs"
	CallbackReferralStats      = "referral_stats"
	CallbackMyReferrals        = "my_referrals"
	CallbackMyReferralLink     = "my_referral_link"
	CallbackReferralRanking    = "referral_ranking"
)

// Префиксы для динамических callback'ов
const (
	// Планы подписок
	CallbackPrefixPlan       = "plan_"        // plan_plan_1m
	CallbackPrefixSelectPlan = "select_plan_" // select_plan_plan_1m (выбор плана)
	CallbackPrefixCreatePlan = "create_plan_" // create_plan_plan_1m
	CallbackPrefixExtendPlan = "extend_plan_" // extend_plan_plan_1m_sub_sub123

	// Способы оплаты
	CallbackPrefixPayCard  = "pay_card_"  // pay_card_plan_1m
	CallbackPrefixPaySBP   = "pay_sbp_"   // pay_sbp_plan_1m
	CallbackPrefixPayStars = "pay_stars_" // pay_stars_plan_1m

	// Подписки
	CallbackPrefixViewSubscription   = "view_subscription_"   // view_subscription_sub123
	CallbackPrefixRenameSubscription = "rename_subscription_" // rename_subscription_sub123
	CallbackPrefixExtendSubscription = "extend_subscription_" // extend_subscription_sub123
	CallbackPrefixDeleteSubscription = "delete_subscription_" // delete_subscription_sub123

	// VPN конфигурации
	CallbackPrefixCreateWireguard   = "create_wireguard"
	CallbackPrefixCreateShadowsocks = "create_shadowsocks"
	CallbackPrefixViewConfig        = "view_config_"      // view_config_config123
	CallbackPrefixDeleteConfig      = "delete_config_"    // delete_config_config123
	CallbackPrefixConnectionGuide   = "connection_guide_" // connection_guide_sub123
)

// Вспомогательные функции для работы с callback'ами

// ParsePlanCallback извлекает ID плана из callback'а
func ParsePlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPlan) && callbackData[:len(CallbackPrefixPlan)] == CallbackPrefixPlan {
		return callbackData[len(CallbackPrefixPlan):], true
	}
	return "", false
}

// ParseSelectPlanCallback извлекает ID плана из callback'а выбора плана
func ParseSelectPlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixSelectPlan) && callbackData[:len(CallbackPrefixSelectPlan)] == CallbackPrefixSelectPlan {
		return callbackData[len(CallbackPrefixSelectPlan):], true
	}
	return "", false
}

// ParseCreatePlanCallback извлекает ID плана из callback'а создания
func ParseCreatePlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixCreatePlan) && callbackData[:len(CallbackPrefixCreatePlan)] == CallbackPrefixCreatePlan {
		return callbackData[len(CallbackPrefixCreatePlan):], true
	}
	return "", false
}

// ParsePayCardCallback извлекает ID плана из callback'а оплаты картой
func ParsePayCardCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPayCard) && callbackData[:len(CallbackPrefixPayCard)] == CallbackPrefixPayCard {
		return callbackData[len(CallbackPrefixPayCard):], true
	}
	return "", false
}

// ParsePaySBPCallback извлекает ID плана из callback'а оплаты СБП
func ParsePaySBPCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPaySBP) && callbackData[:len(CallbackPrefixPaySBP)] == CallbackPrefixPaySBP {
		return callbackData[len(CallbackPrefixPaySBP):], true
	}
	return "", false
}

// ParsePayStarsCallback извлекает ID плана из callback'а оплаты Stars
func ParsePayStarsCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPayStars) && callbackData[:len(CallbackPrefixPayStars)] == CallbackPrefixPayStars {
		return callbackData[len(CallbackPrefixPayStars):], true
	}
	return "", false
}

// ParseExtendPlanCallback извлекает ID плана и подписки из callback'а продления
func ParseExtendPlanCallback(callbackData string) (planID, subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixExtendPlan) && callbackData[:len(CallbackPrefixExtendPlan)] == CallbackPrefixExtendPlan {
		// Формат: extend_plan_plan_1m_sub_sub123
		rest := callbackData[len(CallbackPrefixExtendPlan):]

		// Ищем разделитель "_sub_"
		if idx := len(rest) - len("_sub_"); idx > 0 && rest[idx:] == "_sub_" {
			planID = rest[:idx]
			subscriptionID = rest[idx+len("_sub_"):]
			return planID, subscriptionID, true
		}
	}
	return "", "", false
}

// ParseViewSubscriptionCallback извлекает ID подписки из callback'а просмотра
func ParseViewSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixViewSubscription) && callbackData[:len(CallbackPrefixViewSubscription)] == CallbackPrefixViewSubscription {
		return callbackData[len(CallbackPrefixViewSubscription):], true
	}
	return "", false
}

// ParseRenameSubscriptionCallback извлекает ID подписки из callback'а переименования
func ParseRenameSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixRenameSubscription) && callbackData[:len(CallbackPrefixRenameSubscription)] == CallbackPrefixRenameSubscription {
		return callbackData[len(CallbackPrefixRenameSubscription):], true
	}
	return "", false
}

// ParseExtendSubscriptionCallback извлекает ID подписки из callback'а продления
func ParseExtendSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixExtendSubscription) && callbackData[:len(CallbackPrefixExtendSubscription)] == CallbackPrefixExtendSubscription {
		return callbackData[len(CallbackPrefixExtendSubscription):], true
	}
	return "", false
}

// ParseDeleteSubscriptionCallback извлекает ID подписки из callback'а удаления
func ParseDeleteSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixDeleteSubscription) && callbackData[:len(CallbackPrefixDeleteSubscription)] == CallbackPrefixDeleteSubscription {
		return callbackData[len(CallbackPrefixDeleteSubscription):], true
	}
	return "", false
}

// ParseViewConfigCallback извлекает ID конфигурации из callback'а просмотра
func ParseViewConfigCallback(callbackData string) (configID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixViewConfig) && callbackData[:len(CallbackPrefixViewConfig)] == CallbackPrefixViewConfig {
		return callbackData[len(CallbackPrefixViewConfig):], true
	}
	return "", false
}

// ParseDeleteConfigCallback извлекает ID конфигурации из callback'а удаления
func ParseDeleteConfigCallback(callbackData string) (configID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixDeleteConfig) && callbackData[:len(CallbackPrefixDeleteConfig)] == CallbackPrefixDeleteConfig {
		return callbackData[len(CallbackPrefixDeleteConfig):], true
	}
	return "", false
}

// ParseConnectionGuideCallback извлекает ID подписки из callback'а инструкции по подключению
func ParseConnectionGuideCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixConnectionGuide) && callbackData[:len(CallbackPrefixConnectionGuide)] == CallbackPrefixConnectionGuide {
		return callbackData[len(CallbackPrefixConnectionGuide):], true
	}
	return "", false
}
