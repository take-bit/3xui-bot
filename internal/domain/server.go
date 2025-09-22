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
	Token          string       `json:"token"`
	Status         ServerStatus `json:"status"`
	MaxClients     int          `json:"max_clients"`
	CurrentClients int          `json:"current_clients"`
	Priority       int          `json:"priority"`    // приоритет сервера (1 = высший)
	Region         string       `json:"region"`      // регион сервера (EU, US, ASIA, etc.)
	Description    string       `json:"description"` // описание сервера
	Enabled        bool         `json:"enabled"`     // включен ли сервер
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// IsAvailable проверяет, доступен ли сервер для новых клиентов
func (s *Server) IsAvailable() bool {
	return s.Status == ServerStatusActive && s.Enabled && s.CurrentClients < s.MaxClients
}

// GetLoadPercentage возвращает процент загрузки сервера
func (s *Server) GetLoadPercentage() float64 {
	if s.MaxClients == 0 {
		return 0
	}
	return float64(s.CurrentClients) / float64(s.MaxClients) * 100
}

// GetLoadRatio возвращает коэффициент загрузки сервера (0.0-1.0)
func (s *Server) GetLoadRatio() float64 {
	if s.MaxClients == 0 {
		return 0
	}
	return float64(s.CurrentClients) / float64(s.MaxClients)
}

// IsOverloaded проверяет, перегружен ли сервер
func (s *Server) IsOverloaded(threshold float64) bool {
	return s.GetLoadRatio() > threshold
}

// GetAvailableSlots возвращает количество доступных слотов
func (s *Server) GetAvailableSlots() int {
	if s.MaxClients == 0 {
		return 0
	}
	available := s.MaxClients - s.CurrentClients
	if available < 0 {
		return 0
	}
	return available
}

// IsHealthy проверяет, здоров ли сервер
func (s *Server) IsHealthy() bool {
	return s.Status == ServerStatusActive && s.Enabled
}
