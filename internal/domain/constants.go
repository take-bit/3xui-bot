package domain

// Поддерживаемые валюты
const (
	CurrencyRUB = "RUB"
	CurrencyUSD = "USD"
	CurrencyXTR = "XTR" // Telegram Stars
)

// Поддерживаемые языки
const (
	LanguageRussian = "ru"
	LanguageEnglish = "en"
)

// Поддерживаемые периоды подписки (в днях)
const (
	Duration30  = 30
	Duration60  = 60
	Duration180 = 180
	Duration365 = 365
)

// Лимиты
const (
	MaxUsernameLength     = 32
	MaxFirstNameLength    = 64
	MaxLastNameLength     = 64
	MaxPromocodeLength    = 50
	MaxNotificationLength = 4096
	MaxTitleLength        = 256
)

// Telegram Bot Commands
const (
	CommandStart   = "/start"
	CommandHelp    = "/help"
	CommandProfile = "/profile"
	CommandSupport = "/support"
	CommandAdmin   = "/admin"
)

// Callback Data Prefixes
const (
	CallbackBuyPlan      = "buy_plan_"
	CallbackPayment      = "payment_"
	CallbackPromocode    = "promocode_"
	CallbackReferral     = "referral_"
	CallbackAdmin        = "admin_"
	CallbackNotification = "notification_"
)

// Webhook Paths
const (
	WebhookCryptomus = "/webhook/cryptomus"
	WebhookHeleket   = "/webhook/heleket"
	WebhookYooKassa  = "/webhook/yookassa"
	WebhookYooMoney  = "/webhook/yoomoney"
)

// Default Values
const (
	DefaultTrialDays        = 3
	DefaultReferralReward   = 7
	DefaultMaxClients       = 1000
	DefaultServerPort       = 2096
	DefaultSubscriptionPath = "/user/"
)

// Message Templates
const (
	MsgWelcome = `Добро пожаловать в VPN Shop! 🚀

Здесь вы можете приобрести VPN подписку для безопасного и быстрого интернета.

Выберите действие:`

	MsgProfile = `👤 Ваш профиль

🆔 ID: %d
👤 Имя: %s
📅 Регистрация: %s
📊 Статус: %s`

	MsgSubscriptionActive = `✅ Активная подписка

📅 Осталось дней: %d
📱 Устройств: %d
🔗 Ссылка: %s`

	MsgSubscriptionExpired = `❌ Подписка истекла

Продлите подписку для продолжения использования VPN.`
)
