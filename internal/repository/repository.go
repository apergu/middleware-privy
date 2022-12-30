package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Command interface {
	BeginTx(ctx context.Context) (tx pgx.Tx, err error)
	CommitTx(ctx context.Context, tx pgx.Tx) (err error)
	RollbackTx(ctx context.Context, tx pgx.Tx) (err error)
}
