package postgres

import (
	"context"

	transactorPgx "github.com/Thiht/transactor/pgx"
)

type UoW struct {
	trx *transactorPgx.Transactor
}

func NewUoW(trx *transactorPgx.Transactor) *UoW {

	return &UoW{trx: trx}
}

func (u *UoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {

	return u.trx.WithinTransaction(ctx, fn)
}
