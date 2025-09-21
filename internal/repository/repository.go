package repository

import (
	"context"

	"github.com/Thiht/transactor/pgx"
)

type DBGetter interface {
	GetDB(ctx context.Context) pgx.DB
}
