package config

import (
	"fmt"
	"os"
)

// Config конфигурация приложения
type Config struct {
	Bot     BotConfig
	DB      DBConfig
	Marzban MarzbanConfig
}

// BotConfig конфигурация Telegram бота
type BotConfig struct {
	Token string
	Debug bool
}

// DBConfig конфигурация базы данных
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// MarzbanConfig конфигурация Marzban
type MarzbanConfig struct {
	BaseURL  string
	Username string
	Password string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	return &Config{
		Bot: BotConfig{
			Token: botToken,
			Debug: os.Getenv("BOT_DEBUG") == "true",
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "3xui_bot"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Marzban: MarzbanConfig{
			BaseURL:  getEnv("MARZBAN_BASE_URL", "http://localhost:8000"),
			Username: getEnv("MARZBAN_USERNAME", ""),
			Password: getEnv("MARZBAN_PASSWORD", ""),
		},
	}, nil
}

// GetConnectionString возвращает строку подключения к БД
func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
