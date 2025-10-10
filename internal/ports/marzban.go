package ports

import (
	"context"

	"3xui-bot/internal/core"
)

// Marzban интерфейс для работы с Marzban VPN Manager
type Marzban interface {
	// Authenticate аутентифицируется в Marzban API
	Authenticate(ctx context.Context) error

	// CreateUser создает пользователя в Marzban
	CreateUser(ctx context.Context, user *core.MarzbanUserData) (*core.MarzbanUserData, error)

	// GetUser получает информацию о пользователе
	GetUser(ctx context.Context, username string) (*core.MarzbanUserData, error)

	// UpdateUser обновляет данные пользователя
	UpdateUser(ctx context.Context, username string, user *core.MarzbanUserData) (*core.MarzbanUserData, error)

	// DeleteUser удаляет пользователя
	DeleteUser(ctx context.Context, username string) error

	// ResetUserTraffic сбрасывает трафик пользователя
	ResetUserTraffic(ctx context.Context, username string) error

	// GetStats получает статистику системы
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
