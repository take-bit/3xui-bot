package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	// Telegram Bot
	BotToken string `json:"bot_token"`
	BotURL   string `json:"bot_url"`

	XUIServers []XUIServerConfig `json:"xui_servers"`

	DatabaseURL string         `json:"database_url"`
	Database    DatabaseConfig `json:"database"`

	PaymentMethods PaymentConfig `json:"payment_methods"`

	// Trial and Referral
	TrialConfig    TrialConfig    `json:"trial_config"`
	ReferralConfig ReferralConfig `json:"referral_config"`

	// Server Management
	ServerManagement ServerManagementConfig `json:"server_management"`

	// HTTP Server
	Server ServerConfig `json:"server"`

	// Logging
	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`

	// Plans Configuration
	PlansConfig PlansConfig `json:"plans_config"`

	// Notifications
	NotificationsConfig NotificationsConfig `json:"notifications_config"`

	// Security
	SecurityConfig SecurityConfig `json:"security_config"`

	// Performance
	PerformanceConfig PerformanceConfig `json:"performance_config"`
}

// DatabaseConfig представляет конфигурацию базы данных
type DatabaseConfig struct {
	Host                  string        `json:"host"`
	Port                  int           `json:"port"`
	User                  string        `json:"user"`
	Password              string        `json:"password"`
	Name                  string        `json:"name"`
	SSLMode               string        `json:"ssl_mode"`
	MaxConnections        int           `json:"max_connections"`
	MinConnections        int           `json:"min_connections"`
	MaxConnectionLifetime int           `json:"max_connection_lifetime"`
	MaxConnectionIdleTime int           `json:"max_connection_idle_time"`
	ConnectionTimeout     time.Duration `json:"connection_timeout"`
	QueryTimeout          time.Duration `json:"query_timeout"`
	EnableQueryLog        bool          `json:"enable_query_log"`
	SlowQueryThreshold    time.Duration `json:"slow_query_threshold"`
}

// XUIServerConfig представляет конфигурацию одного сервера 3X-UI
type XUIServerConfig struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	Enabled     bool   `json:"enabled"`
	Priority    int    `json:"priority"`
	MaxClients  int    `json:"max_clients"`
	Region      string `json:"region"`
	Description string `json:"description"`
}

// ServerManagementConfig представляет конфигурацию управления серверами
type ServerManagementConfig struct {
	SelectionStrategy    string        `json:"selection_strategy"`
	HealthCheckInterval  time.Duration `json:"health_check_interval"`
	HealthCheckTimeout   time.Duration `json:"health_check_timeout"`
	MaxRetries           int           `json:"max_retries"`
	RetryDelay           time.Duration `json:"retry_delay"`
	LoadBalanceThreshold float64       `json:"load_balance_threshold"`
	AutoFailover         bool          `json:"auto_failover"`
	GeographicRouting    bool          `json:"geographic_routing"`
}

// PaymentConfig представляет конфигурацию платежных методов
type PaymentConfig struct {
	Cryptomus          CryptomusConfig           `json:"cryptomus"`
	Heleket            HeleketConfig             `json:"heleket"`
	YooKassa           YooKassaConfig            `json:"yookassa"`
	YooMoney           YooMoneyConfig            `json:"yoomoney"`
	TelegramStars      TelegramStarsConfig       `json:"telegram_stars"`
	ProcessingSettings PaymentProcessingSettings `json:"processing_settings"`
}

// CryptomusConfig представляет конфигурацию Cryptomus
type CryptomusConfig struct {
	Enabled     bool   `json:"enabled"`
	APIKey      string `json:"api_key"`
	MerchantID  string `json:"merchant_id"`
	WebhookURL  string `json:"webhook_url"`
	CallbackURL string `json:"callback_url"`
	Currency    string `json:"currency"`
	Network     string `json:"network"`
}

// HeleketConfig представляет конфигурацию Heleket
type HeleketConfig struct {
	Enabled     bool   `json:"enabled"`
	APIKey      string `json:"api_key"`
	MerchantID  string `json:"merchant_id"`
	WebhookURL  string `json:"webhook_url"`
	CallbackURL string `json:"callback_url"`
	Currency    string `json:"currency"`
}

// YooKassaConfig представляет конфигурацию YooKassa
type YooKassaConfig struct {
	Enabled     bool   `json:"enabled"`
	Token       string `json:"token"`
	ShopID      string `json:"shop_id"`
	WebhookURL  string `json:"webhook_url"`
	CallbackURL string `json:"callback_url"`
	Currency    string `json:"currency"`
}

// YooMoneyConfig представляет конфигурацию YooMoney
type YooMoneyConfig struct {
	Enabled            bool   `json:"enabled"`
	WalletID           string `json:"wallet_id"`
	NotificationSecret string `json:"notification_secret"`
	WebhookURL         string `json:"webhook_url"`
	CallbackURL        string `json:"callback_url"`
	Currency           string `json:"currency"`
}

// TelegramStarsConfig представляет конфигурацию Telegram Stars
type TelegramStarsConfig struct {
	Enabled     bool   `json:"enabled"`
	WebhookURL  string `json:"webhook_url"`
	CallbackURL string `json:"callback_url"`
}

// PaymentProcessingSettings представляет настройки обработки платежей
type PaymentProcessingSettings struct {
	ProcessingTimeout   time.Duration              `json:"processing_timeout"`
	MaxRetries          int                        `json:"max_retries"`
	RetryDelay          time.Duration              `json:"retry_delay"`
	AutoConfirm         bool                       `json:"auto_confirm"`
	ConfirmationTimeout time.Duration              `json:"confirmation_timeout"`
	CurrencyConversion  CurrencyConversionSettings `json:"currency_conversion"`
}

// CurrencyConversionSettings представляет настройки конвертации валют
type CurrencyConversionSettings struct {
	Enabled         bool          `json:"enabled"`
	BaseCurrency    string        `json:"base_currency"`
	ExchangeRateAPI string        `json:"exchange_rate_api"`
	UpdateInterval  time.Duration `json:"update_interval"`
}

// TrialConfig представляет конфигурацию пробного периода
type TrialConfig struct {
	Enabled       bool               `json:"enabled"`
	Days          int                `json:"days"`
	ExtendedDays  int                `json:"extended_days"`
	Restrictions  TrialRestrictions  `json:"restrictions"`
	Notifications TrialNotifications `json:"notifications"`
}

// TrialRestrictions представляет ограничения пробного периода
type TrialRestrictions struct {
	MaxDevices     int   `json:"max_devices"`
	MaxConnections int   `json:"max_connections"`
	BandwidthLimit int64 `json:"bandwidth_limit"`
	SpeedLimit     int   `json:"speed_limit"`
}

// TrialNotifications представляет уведомления пробного периода
type TrialNotifications struct {
	Enabled            bool  `json:"enabled"`
	DaysBeforeExpiry   []int `json:"days_before_expiry"`
	ExpiryNotification bool  `json:"expiry_notification"`
}

// ReferralConfig представляет конфигурацию реферальной программы
type ReferralConfig struct {
	Enabled       bool                  `json:"enabled"`
	Level1Reward  int                   `json:"level1_reward"`
	Level2Reward  int                   `json:"level2_reward"`
	MinPayment    int                   `json:"min_payment"`
	Settings      ReferralSettings      `json:"settings"`
	Notifications ReferralNotifications `json:"notifications"`
}

// ReferralSettings представляет настройки реферальной программы
type ReferralSettings struct {
	MaxLevels        int     `json:"max_levels"`
	RewardType       string  `json:"reward_type"`
	RewardPercentage float64 `json:"reward_percentage"`
	MinReward        int     `json:"min_reward"`
	MaxReward        int     `json:"max_reward"`
}

// ReferralNotifications представляет уведомления реферальной программы
type ReferralNotifications struct {
	Enabled         bool `json:"enabled"`
	NewReferral     bool `json:"new_referral"`
	RewardEarned    bool `json:"reward_earned"`
	ReferralPayment bool `json:"referral_payment"`
}

// ServerConfig представляет конфигурацию HTTP сервера
type ServerConfig struct {
	Port         string                  `json:"port"`
	ReadTimeout  time.Duration           `json:"read_timeout"`
	WriteTimeout time.Duration           `json:"write_timeout"`
	IdleTimeout  time.Duration           `json:"idle_timeout"`
	Security     ServerSecurityConfig    `json:"security"`
	Performance  ServerPerformanceConfig `json:"performance"`
}

// ServerSecurityConfig представляет настройки безопасности сервера
type ServerSecurityConfig struct {
	EnableHTTPS     bool          `json:"enable_https"`
	TLSCertFile     string        `json:"tls_cert_file"`
	TLSKeyFile      string        `json:"tls_key_file"`
	AllowedIPs      []string      `json:"allowed_ips"`
	RateLimit       int           `json:"rate_limit"`
	RateLimitWindow time.Duration `json:"rate_limit_window"`
}

// ServerPerformanceConfig представляет настройки производительности сервера
type ServerPerformanceConfig struct {
	MaxConnections    int           `json:"max_connections"`
	MaxRequestSize    int64         `json:"max_request_size"`
	EnableCompression bool          `json:"enable_compression"`
	EnableKeepAlive   bool          `json:"enable_keep_alive"`
	KeepAliveTimeout  time.Duration `json:"keep_alive_timeout"`
}

// PlansConfig представляет конфигурацию планов
type PlansConfig struct {
	FilePath       string                `json:"file_path"`
	AutoReload     bool                  `json:"auto_reload"`
	ReloadInterval time.Duration         `json:"reload_interval"`
	DefaultPlans   bool                  `json:"default_plans"`
	Validation     PlansValidationConfig `json:"validation"`
}

// PlansValidationConfig представляет настройки валидации планов
type PlansValidationConfig struct {
	Enabled           bool `json:"enabled"`
	ValidatePrices    bool `json:"validate_prices"`
	ValidateDurations bool `json:"validate_durations"`
	MinPrice          int  `json:"min_price"`
	MaxPrice          int  `json:"max_price"`
	MinDuration       int  `json:"min_duration"`
	MaxDuration       int  `json:"max_duration"`
}

// NotificationsConfig представляет конфигурацию уведомлений
type NotificationsConfig struct {
	Email    EmailNotificationsConfig    `json:"email"`
	Telegram TelegramNotificationsConfig `json:"telegram"`
	Webhook  WebhookNotificationsConfig  `json:"webhook"`
	Settings NotificationSettings        `json:"settings"`
}

// EmailNotificationsConfig представляет конфигурацию email уведомлений
type EmailNotificationsConfig struct {
	Enabled  bool   `json:"enabled"`
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	UseTLS   bool   `json:"use_tls"`
}

// TelegramNotificationsConfig представляет конфигурацию Telegram уведомлений
type TelegramNotificationsConfig struct {
	Enabled               bool   `json:"enabled"`
	BotToken              string `json:"bot_token"`
	ChatID                string `json:"chat_id"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

// WebhookNotificationsConfig представляет конфигурацию webhook уведомлений
type WebhookNotificationsConfig struct {
	Enabled    bool          `json:"enabled"`
	URL        string        `json:"url"`
	Secret     string        `json:"secret"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
}

// NotificationSettings представляет настройки уведомлений
type NotificationSettings struct {
	Enabled             bool `json:"enabled"`
	PaymentReceived     bool `json:"payment_received"`
	PaymentFailed       bool `json:"payment_failed"`
	SubscriptionExpired bool `json:"subscription_expired"`
	SubscriptionRenewed bool `json:"subscription_renewed"`
	NewUser             bool `json:"new_user"`
	NewReferral         bool `json:"new_referral"`
	ServerError         bool `json:"server_error"`
}

// SecurityConfig представляет конфигурацию безопасности
type SecurityConfig struct {
	API    APISecurityConfig   `json:"api"`
	Data   DataSecurityConfig  `json:"data"`
	Access AccessControlConfig `json:"access"`
}

// APISecurityConfig представляет настройки безопасности API
type APISecurityConfig struct {
	EnableAPIKeyAuth bool     `json:"enable_api_key_auth"`
	APIKeyHeader     string   `json:"api_key_header"`
	APIKeyValue      string   `json:"api_key_value"`
	EnableCORS       bool     `json:"enable_cors"`
	CORSOrigins      []string `json:"cors_origins"`
	EnableCSRF       bool     `json:"enable_csrf"`
	CSRFSecret       string   `json:"csrf_secret"`
}

// DataSecurityConfig представляет настройки безопасности данных
type DataSecurityConfig struct {
	EnableEncryption bool   `json:"enable_encryption"`
	EncryptionKey    string `json:"encryption_key"`
	EnableHashing    bool   `json:"enable_hashing"`
	HashAlgorithm    string `json:"hash_algorithm"`
	EnableMasking    bool   `json:"enable_masking"`
	MaskingPattern   string `json:"masking_pattern"`
}

// AccessControlConfig представляет настройки контроля доступа
type AccessControlConfig struct {
	EnableIPWhitelist bool          `json:"enable_ip_whitelist"`
	IPWhitelist       []string      `json:"ip_whitelist"`
	EnableIPBlacklist bool          `json:"enable_ip_blacklist"`
	IPBlacklist       []string      `json:"ip_blacklist"`
	EnableRateLimit   bool          `json:"enable_rate_limit"`
	RateLimit         int           `json:"rate_limit"`
	RateLimitWindow   time.Duration `json:"rate_limit_window"`
}

// PerformanceConfig представляет конфигурацию производительности
type PerformanceConfig struct {
	Database DatabasePerformanceConfig `json:"database"`
	Cache    CachePerformanceConfig    `json:"cache"`
	Memory   MemoryPerformanceConfig   `json:"memory"`
	CPU      CPUPerformanceConfig      `json:"cpu"`
}

// DatabasePerformanceConfig представляет настройки производительности БД
type DatabasePerformanceConfig struct {
	MaxConnections     int           `json:"max_connections"`
	MinConnections     int           `json:"min_connections"`
	ConnectionTimeout  time.Duration `json:"connection_timeout"`
	QueryTimeout       time.Duration `json:"query_timeout"`
	EnableQueryLog     bool          `json:"enable_query_log"`
	SlowQueryThreshold time.Duration `json:"slow_query_threshold"`
}

// CachePerformanceConfig представляет настройки производительности кэша
type CachePerformanceConfig struct {
	Enabled         bool          `json:"enabled"`
	Type            string        `json:"type"`
	TTL             time.Duration `json:"ttl"`
	MaxSize         int64         `json:"max_size"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
	RedisURL        string        `json:"redis_url"`
	FileCachePath   string        `json:"file_cache_path"`
}

// MemoryPerformanceConfig представляет настройки производительности памяти
type MemoryPerformanceConfig struct {
	MaxMemoryUsage        int64         `json:"max_memory_usage"`
	GCThreshold           int64         `json:"gc_threshold"`
	EnableMemoryProfiling bool          `json:"enable_memory_profiling"`
	ProfilingInterval     time.Duration `json:"profiling_interval"`
}

// CPUPerformanceConfig представляет настройки производительности CPU
type CPUPerformanceConfig struct {
	MaxCPUUsage        float64       `json:"max_cpu_usage"`
	EnableCPUProfiling bool          `json:"enable_cpu_profiling"`
	ProfilingInterval  time.Duration `json:"profiling_interval"`
	WorkerPoolSize     int           `json:"worker_pool_size"`
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig() (*Config, error) {
	// Создаем конфигурацию по умолчанию
	config := &Config{
		BotToken:  os.Getenv("BOT_TOKEN"),
		BotURL:    os.Getenv("BOT_URL"),
		LogLevel:  "info",
		LogFormat: "json",
	}

	// Пытаемся загрузить из файла config.json
	if _, err := os.Stat("config.json"); err == nil {
		data, err := os.ReadFile("config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read config.json: %w", err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config.json: %w", err)
		}
	}

	// Переменные окружения имеют приоритет
	if botToken := os.Getenv("BOT_TOKEN"); botToken != "" {
		config.BotToken = botToken
	}
	if botURL := os.Getenv("BOT_URL"); botURL != "" {
		config.BotURL = botURL
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	// Проверяем обязательные поля
	if config.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	return config, nil
}
