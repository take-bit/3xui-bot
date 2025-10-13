package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Bot       BotConfig       `json:"bot"`
	DB        DBConfig        `json:"db"`
	Marzban   MarzbanConfig   `json:"marzban"`
	Payment   PaymentConfig   `json:"payment"`
	Scheduler SchedulerConfig `json:"scheduler"`
	Logging   LoggingConfig   `json:"logging"`
}

type BotConfig struct {
	Token              string  `env:"BOT_TOKEN,required"`
	Debug              bool    `json:"debug"`
	Timeout            int     `json:"timeout"`
	UpdatesChannelSize int     `json:"updates_channel_size"`
	MaxConcurrent      int     `json:"max_concurrent"`
	AdminIDs           []int64 `json:"admin_ids"`
	SupportUsername    string  `json:"support_username"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

type MarzbanConfig struct {
	BaseURL  string `json:"base_url"`
	Username string `env:"MARZBAN_USERNAME,required"`
	Password string `env:"MARZBAN_PASSWORD,required"`
}

type PaymentConfig struct {
	// TODO: пока mock, позже можно добавить ключи API и указать тэги env
}

type SchedulerConfig struct {
	Enabled bool `json:"enabled"`
}

type LoggingConfig struct {
	Level string `json:"level"`
}

func Load(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	cfg := &Config{}

	if err := json.NewDecoder(configFile).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	normalize(cfg)

	if err := validateRequired(cfg); err != nil {
		return nil, err
	}

	applyDefaults(cfg)

	return cfg, nil
}

func normalize(cfg *Config) {
	cfg.Marzban.BaseURL = strings.TrimSpace(cfg.Marzban.BaseURL)
	cfg.Marzban.BaseURL = strings.TrimRight(cfg.Marzban.BaseURL, "/")

	cfg.DB.Host = strings.TrimSpace(cfg.DB.Host)
	cfg.DB.Port = strings.TrimSpace(cfg.DB.Port)
	cfg.DB.Database = strings.TrimSpace(cfg.DB.Database)
	cfg.DB.SSLMode = strings.TrimSpace(strings.ToLower(cfg.DB.SSLMode))

	cfg.Logging.Level = strings.TrimSpace(strings.ToLower(cfg.Logging.Level))
	cfg.Bot.SupportUsername = strings.TrimSpace(cfg.Bot.SupportUsername)
}

func validateRequired(cfg *Config) error {
	var errs []string

	if cfg.Marzban.BaseURL == "" {
		errs = append(errs, "marzban.base_url is required (set in JSON)")
	}

	if cfg.DB.Host == "" {
		errs = append(errs, "db.host is required (set in JSON)")
	}
	if cfg.DB.Database == "" {
		errs = append(errs, "db.database is required (set in JSON)")
	}

	if len(errs) > 0 {
		return errors.New("invalid config: " + strings.Join(errs, "; "))
	}
	return nil
}

func applyDefaults(cfg *Config) {
	if cfg.Bot.Timeout == 0 {
		cfg.Bot.Timeout = 30
	}
	if cfg.Bot.UpdatesChannelSize == 0 {
		cfg.Bot.UpdatesChannelSize = 100
	}
	if cfg.Bot.MaxConcurrent == 0 {
		cfg.Bot.MaxConcurrent = runtime.NumCPU()
	}

	if cfg.DB.Host == "" {
		cfg.DB.Host = "localhost"
	}
	if cfg.DB.Port == "" {
		cfg.DB.Port = "5432"
	}
	if cfg.DB.Database == "" {
		cfg.DB.Database = "postgres"
	}
	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
}
