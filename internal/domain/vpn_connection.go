package domain

import "time"

// VPNConnection представляет VPN подключение пользователя
type VPNConnection struct {
	UserID       int64     `json:"user_id"`
	ServerID     int64     `json:"server_id"`
	XUIInboundID int       `json:"xui_inbound_id"`
	XUIClientID  string    `json:"xui_client_id"`
	UUID         string    `json:"uuid"`
	Email        string    `json:"email"`
	ConfigURL    string    `json:"config_url"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsActive проверяет, активно ли VPN подключение
func (v *VPNConnection) IsActive() bool {
	return time.Now().Before(v.ExpiresAt)
}

// IsExpired проверяет, истекло ли VPN подключение
func (v *VPNConnection) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// GetRemainingDays возвращает количество оставшихся дней
func (v *VPNConnection) GetRemainingDays() int {
	if v.IsExpired() {
		return 0
	}

	remaining := time.Until(v.ExpiresAt)
	return int(remaining.Hours() / 24)
}

// GetStatus возвращает статус подключения
func (v *VPNConnection) GetStatus() string {
	if v.IsExpired() {
		return "expired"
	}
	if v.IsActive() {
		return "active"
	}
	return "inactive"
}
