package ports

import (
	"context"

	"3xui-bot/internal/core"
)

type Marzban interface {
	Authenticate(ctx context.Context) error

	GetInbounds(ctx context.Context) ([]map[string]interface{}, error)

	CreateUser(ctx context.Context, user *core.MarzbanUserData) (*core.MarzbanUserData, error)

	GetUser(ctx context.Context, username string) (*core.MarzbanUserData, error)

	UpdateUser(ctx context.Context, username string, user *core.MarzbanUserData) (*core.MarzbanUserData, error)

	DeleteUser(ctx context.Context, username string) error

	ResetUserTraffic(ctx context.Context, username string) error

	GetStats(ctx context.Context) (map[string]interface{}, error)
}
