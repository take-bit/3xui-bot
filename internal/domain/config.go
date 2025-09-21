package domain

import "time"

// Config представляет конфигурацию приложения
type Config struct {
	// Telegram Bot
	BotToken string `json:"bot_token"`
	BotURL   string `json:"bot_url"`

	// 3X-UI Panel
	XUIUsername string `json:"xui_username"`
	XUIPassword string `json:"xui_password"`
	XUIToken    string `json:"xui_token"`
	XUIBaseURL  string `json:"xui_base_url"`

	// Database
	DatabaseURL string `json:"database_url"`

	// Payment Methods
	PaymentMethods PaymentConfig `json:"payment_methods"`

	// Trial and Referral
	TrialConfig    TrialConfig    `json:"trial_config"`
	ReferralConfig ReferralConfig `json:"referral_config"`

	// Server
	Server ServerConfig `json:"server"`

	// Logging
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`
}

// PaymentConfig представляет конфигурацию платежных методов
type PaymentConfig struct {
	Cryptomus     CryptomusConfig     `json:"cryptomus"`
	Heleket       HeleketConfig       `json:"heleket"`
	YooKassa      YooKassaConfig      `json:"yookassa"`
	YooMoney      YooMoneyConfig      `json:"yoomoney"`
	TelegramStars TelegramStarsConfig `json:"telegram_stars"`
}

// CryptomusConfig представляет конфигурацию Cryptomus
type CryptomusConfig struct {
	Enabled    bool   `json:"enabled"`
	APIKey     string `json:"api_key"`
	MerchantID string `json:"merchant_id"`
}

// HeleketConfig представляет конфигурацию Heleket
type HeleketConfig struct {
	Enabled    bool   `json:"enabled"`
	APIKey     string `json:"api_key"`
	MerchantID string `json:"merchant_id"`
}

// YooKassaConfig представляет конфигурацию YooKassa
type YooKassaConfig struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	ShopID  string `json:"shop_id"`
}

// YooMoneyConfig представляет конфигурацию YooMoney
type YooMoneyConfig struct {
	Enabled            bool   `json:"enabled"`
	WalletID           string `json:"wallet_id"`
	NotificationSecret string `json:"notification_secret"`
}

// TelegramStarsConfig представляет конфигурацию Telegram Stars
type TelegramStarsConfig struct {
	Enabled bool `json:"enabled"`
}

// TrialConfig представляет конфигурацию пробного периода
type TrialConfig struct {
	Enabled      bool `json:"enabled"`
	Days         int  `json:"days"`
	ExtendedDays int  `json:"extended_days"` // для рефералов
}

// ReferralConfig представляет конфигурацию реферальной программы
type ReferralConfig struct {
	Enabled      bool `json:"enabled"`
	Level1Reward int  `json:"level1_reward"` // награда за первого уровня
	Level2Reward int  `json:"level2_reward"` // награда за второй уровень
	MinPayment   int  `json:"min_payment"`   // минимальная сумма платежа для награды
}

// ServerConfig представляет конфигурацию сервера
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}
