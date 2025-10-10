package core

import (
	"fmt"
	"time"
)

// VPNConnection представляет VPN подключение пользователя в Marzban
type VPNConnection struct {
	ID              string    `json:"id" db:"id"`
	TelegramUserID  int64     `json:"telegram_user_id" db:"telegram_user_id"`
	MarzbanUsername string    `json:"marzban_username" db:"marzban_username"` // Username в Marzban
	Name            string    `json:"name" db:"name"`                         // Локальное имя
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Данные из Marzban API (не хранятся в БД)
	ExpireAt       *time.Time             `json:"expire_at,omitempty"`
	DataLimitBytes *int64                 `json:"data_limit_bytes,omitempty"`
	DataUsedBytes  *int64                 `json:"data_used_bytes,omitempty"`
	Status         string                 `json:"status,omitempty"`
	ProtocolConfig map[string]interface{} `json:"protocol_config,omitempty"`
}

// GetDisplayName возвращает отображаемое имя подключения
func (v *VPNConnection) GetDisplayName() string {
	if v.Name != "" {
		return v.Name
	}
	return fmt.Sprintf("VPN Connection %s", v.MarzbanUsername)
}

// IsExpired проверяет, истекло ли подключение
func (v *VPNConnection) IsExpired() bool {
	if v.ExpireAt == nil {
		return false // Безлимитное подключение
	}
	return v.ExpireAt.Before(time.Now())
}

// IsDataLimitReached проверяет, достигнут ли лимит трафика
func (v *VPNConnection) IsDataLimitReached() bool {
	if v.DataLimitBytes == nil || *v.DataLimitBytes == 0 {
		return false // Безлимитное подключение
	}
	if v.DataUsedBytes == nil {
		return false
	}
	return *v.DataUsedBytes >= *v.DataLimitBytes
}

// IsValid проверяет валидность подключения
func (v *VPNConnection) IsValid() bool {
	return v.MarzbanUsername != "" && v.TelegramUserID > 0
}

// GetStatusText возвращает статус подключения
func (v *VPNConnection) GetStatusText() string {
	if v.IsExpired() {
		return "Истекло"
	}
	if v.IsDataLimitReached() {
		return "Лимит трафика"
	}
	switch v.Status {
	case "active":
		return "Активно"
	case "disabled":
		return "Отключено"
	case "expired":
		return "Истекло"
	case "limited":
		return "Лимит трафика"
	default:
		return "Неизвестно"
	}
}

// GetDataUsageText возвращает текст использования трафика
func (v *VPNConnection) GetDataUsageText() string {
	if v.DataLimitBytes == nil || *v.DataLimitBytes == 0 {
		if v.DataUsedBytes == nil {
			return "Использовано: 0 B"
		}
		return fmt.Sprintf("Использовано: %s", formatBytes(*v.DataUsedBytes))
	}

	if v.DataUsedBytes == nil {
		usedBytes := int64(0)
		v.DataUsedBytes = &usedBytes
	}

	usedGB := float64(*v.DataUsedBytes) / (1024 * 1024 * 1024)
	limitGB := float64(*v.DataLimitBytes) / (1024 * 1024 * 1024)
	percent := (float64(*v.DataUsedBytes) / float64(*v.DataLimitBytes)) * 100

	return fmt.Sprintf("Использовано: %.2f GB / %.2f GB (%.1f%%)", usedGB, limitGB, percent)
}

// formatBytes форматирует байты в читаемый вид
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// MarzbanUserData представляет данные пользователя из Marzban API
type MarzbanUserData struct {
	Username         string                 `json:"username"`
	Expire           *int64                 `json:"expire"`     // UTC timestamp
	DataLimit        *int64                 `json:"data_limit"` // в байтах
	DataUsed         *int64                 `json:"data_used"`  // в байтах
	Status           string                 `json:"status"`     // active, disabled, expired
	Proxies          map[string]interface{} `json:"proxies"`    // Конфигурация протоколов
	Inbounds         map[string][]string    `json:"inbounds"`   // Входящие соединения
	Note             string                 `json:"note"`
	SubUpdatedAt     *string                `json:"sub_updated_at"`
	SubLastUserAgent *string                `json:"sub_last_user_agent"`
	OnlineAt         *string                `json:"online_at"`
	OnHoldTimeout    *string                `json:"on_hold_timeout"`
	CreatedAt        *string                `json:"created_at"`
	Links            []string               `json:"links"`            // Ссылки для подключения
	SubscriptionURL  string                 `json:"subscription_url"` // URL подписки
}

// IsExpired проверяет, истекло ли подключение
func (m *MarzbanUserData) IsExpired() bool {
	if m.Expire == nil || *m.Expire == 0 {
		return false // Безлимитное подключение
	}
	return *m.Expire < time.Now().Unix()
}

// IsDataLimitReached проверяет, достигнут ли лимит трафика
func (m *MarzbanUserData) IsDataLimitReached() bool {
	if m.DataLimit == nil || *m.DataLimit == 0 {
		return false // Безлимитное подключение
	}
	if m.DataUsed == nil {
		return false
	}
	return *m.DataUsed >= *m.DataLimit
}

// GetStatusText возвращает текстовое описание статуса
func (m *MarzbanUserData) GetStatusText() string {
	if m.IsExpired() {
		return "Истекло"
	}
	if m.IsDataLimitReached() {
		return "Лимит трафика"
	}
	switch m.Status {
	case "active":
		return "Активно"
	case "disabled":
		return "Отключено"
	case "expired":
		return "Истекло"
	case "limited":
		return "Лимит трафика"
	default:
		return "Неизвестно"
	}
}

// GetDataUsageText возвращает текст использования трафика
func (m *MarzbanUserData) GetDataUsageText() string {
	if m.DataLimit == nil || *m.DataLimit == 0 {
		if m.DataUsed == nil {
			return "Использовано: 0 B"
		}
		return fmt.Sprintf("Использовано: %s", formatBytes(*m.DataUsed))
	}

	if m.DataUsed == nil {
		usedBytes := int64(0)
		m.DataUsed = &usedBytes
	}

	usedGB := float64(*m.DataUsed) / (1024 * 1024 * 1024)
	limitGB := float64(*m.DataLimit) / (1024 * 1024 * 1024)
	percent := (float64(*m.DataUsed) / float64(*m.DataLimit)) * 100

	return fmt.Sprintf("Использовано: %.2f GB / %.2f GB (%.1f%%)", usedGB, limitGB, percent)
}

// GetExpireText возвращает текст с датой истечения
func (m *MarzbanUserData) GetExpireText() string {
	if m.Expire == nil || *m.Expire == 0 {
		return "Безлимитно"
	}
	return time.Unix(*m.Expire, 0).Format("02.01.2006 15:04")
}

// ExpireAt возвращает время истечения
func (m *MarzbanUserData) ExpireAt() *time.Time {
	if m.Expire == nil || *m.Expire == 0 {
		return nil
	}
	expireTime := time.Unix(*m.Expire, 0)
	return &expireTime
}
