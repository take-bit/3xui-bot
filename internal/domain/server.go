package domain

import (
	"time"
)

// ServerStatus представляет статус сервера
type ServerStatus string

const (
	ServerStatusActive   ServerStatus = "active"
	ServerStatusInactive ServerStatus = "inactive"
	ServerStatusError    ServerStatus = "error"
)

// Server представляет сервер 3X-UI
type Server struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Host           string       `json:"host"`
	Port           int          `json:"port"`
	Username       string       `json:"username"`
	Password       string       `json:"password"`
	Status         ServerStatus `json:"status"`
	MaxClients     int          `json:"max_clients"`
	CurrentClients int          `json:"current_clients"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// IsAvailable проверяет, доступен ли сервер для новых клиентов
func (s *Server) IsAvailable() bool {
	return s.Status == ServerStatusActive && s.CurrentClients < s.MaxClients
}

// GetLoadPercentage возвращает процент загрузки сервера
func (s *Server) GetLoadPercentage() float64 {
	if s.MaxClients == 0 {
		return 0
	}
	return float64(s.CurrentClients) / float64(s.MaxClients) * 100
}
