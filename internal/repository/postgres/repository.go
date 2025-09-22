package postgres

import (
	"context"

	pgxTransactor "github.com/Thiht/transactor/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type Repository struct {
	pool       *pgxpool.Pool
	dbGetter   pgxTransactor.DBGetter
	transactor *pgxTransactor.Transactor
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	transactor, dbGetter := pgxTransactor.NewTransactorFromPool(pool)

	return &Repository{
		pool:       pool,
		dbGetter:   dbGetter,
		transactor: transactor,
	}
}

func (r *Repository) GetDB(ctx context.Context) DB {
	return r.dbGetter(ctx)
}

func (r *Repository) GetTransactor() *pgxTransactor.Transactor {
	return r.transactor
}
