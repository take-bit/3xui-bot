package ui

const (
	CommandStart = "start"
	CommandHelp  = "help"
)

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

const (
	CallbackPrefixPlan       = "plan_"
	CallbackPrefixSelectPlan = "select_plan_"
	CallbackPrefixCreatePlan = "create_plan_"
	CallbackPrefixExtendPlan = "extend_plan_"

	CallbackPrefixPayCard  = "pay_card_"
	CallbackPrefixPaySBP   = "pay_sbp_"
	CallbackPrefixPayStars = "pay_stars_"

	CallbackPrefixViewSubscription   = "view_subscription_"
	CallbackPrefixRenameSubscription = "rename_subscription_"
	CallbackPrefixExtendSubscription = "extend_subscription_"
	CallbackPrefixDeleteSubscription = "delete_subscription_"

	CallbackPrefixCreateWireguard   = "create_wireguard"
	CallbackPrefixCreateShadowsocks = "create_shadowsocks"
	CallbackPrefixViewConfig        = "view_config_"
	CallbackPrefixDeleteConfig      = "delete_config_"
	CallbackPrefixConnectionGuide   = "connection_guide_"
)

func ParsePlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPlan) && callbackData[:len(CallbackPrefixPlan)] == CallbackPrefixPlan {

		return callbackData[len(CallbackPrefixPlan):], true
	}

	return "", false
}

func ParseSelectPlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixSelectPlan) && callbackData[:len(CallbackPrefixSelectPlan)] == CallbackPrefixSelectPlan {

		return callbackData[len(CallbackPrefixSelectPlan):], true
	}

	return "", false
}

func ParseCreatePlanCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixCreatePlan) && callbackData[:len(CallbackPrefixCreatePlan)] == CallbackPrefixCreatePlan {

		return callbackData[len(CallbackPrefixCreatePlan):], true
	}

	return "", false
}

func ParsePayCardCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPayCard) && callbackData[:len(CallbackPrefixPayCard)] == CallbackPrefixPayCard {

		return callbackData[len(CallbackPrefixPayCard):], true
	}

	return "", false
}

func ParsePaySBPCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPaySBP) && callbackData[:len(CallbackPrefixPaySBP)] == CallbackPrefixPaySBP {

		return callbackData[len(CallbackPrefixPaySBP):], true
	}

	return "", false
}

func ParsePayStarsCallback(callbackData string) (planID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixPayStars) && callbackData[:len(CallbackPrefixPayStars)] == CallbackPrefixPayStars {

		return callbackData[len(CallbackPrefixPayStars):], true
	}

	return "", false
}

func ParseExtendPlanCallback(callbackData string) (planID, subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixExtendPlan) && callbackData[:len(CallbackPrefixExtendPlan)] == CallbackPrefixExtendPlan {
		rest := callbackData[len(CallbackPrefixExtendPlan):]

		if idx := len(rest) - len("_sub_"); idx > 0 && rest[idx:] == "_sub_" {
			planID = rest[:idx]
			subscriptionID = rest[idx+len("_sub_"):]

			return planID, subscriptionID, true
		}
	}

	return "", "", false
}

func ParseViewSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixViewSubscription) && callbackData[:len(CallbackPrefixViewSubscription)] == CallbackPrefixViewSubscription {

		return callbackData[len(CallbackPrefixViewSubscription):], true
	}

	return "", false
}

func ParseRenameSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixRenameSubscription) && callbackData[:len(CallbackPrefixRenameSubscription)] == CallbackPrefixRenameSubscription {

		return callbackData[len(CallbackPrefixRenameSubscription):], true
	}

	return "", false
}

func ParseExtendSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixExtendSubscription) && callbackData[:len(CallbackPrefixExtendSubscription)] == CallbackPrefixExtendSubscription {

		return callbackData[len(CallbackPrefixExtendSubscription):], true
	}

	return "", false
}

func ParseDeleteSubscriptionCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixDeleteSubscription) && callbackData[:len(CallbackPrefixDeleteSubscription)] == CallbackPrefixDeleteSubscription {

		return callbackData[len(CallbackPrefixDeleteSubscription):], true
	}

	return "", false
}

func ParseViewConfigCallback(callbackData string) (configID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixViewConfig) && callbackData[:len(CallbackPrefixViewConfig)] == CallbackPrefixViewConfig {

		return callbackData[len(CallbackPrefixViewConfig):], true
	}

	return "", false
}

func ParseDeleteConfigCallback(callbackData string) (configID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixDeleteConfig) && callbackData[:len(CallbackPrefixDeleteConfig)] == CallbackPrefixDeleteConfig {

		return callbackData[len(CallbackPrefixDeleteConfig):], true
	}

	return "", false
}

func ParseConnectionGuideCallback(callbackData string) (subscriptionID string, ok bool) {
	if len(callbackData) > len(CallbackPrefixConnectionGuide) && callbackData[:len(CallbackPrefixConnectionGuide)] == CallbackPrefixConnectionGuide {

		return callbackData[len(CallbackPrefixConnectionGuide):], true
	}

	return "", false
}
