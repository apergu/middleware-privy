package repository

import (
	"context"

	"middleware/internal/entity"

	"github.com/jackc/pgx/v5"
)

type TransferBalanceFilter struct {
	Sort string
}

type TransferBalanceQueryRepository interface {
	Find(ctx context.Context, filter TransferBalanceFilter, limit, skip int64, tx pgx.Tx) ([]entity.TransferBalance, error)
	Count(ctx context.Context, filter TransferBalanceFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.TransferBalance, error)
}

type TransferBalanceCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.TransferBalance, error)
	Create(ctx context.Context, TransferBalance entity.TransferBalance, tx pgx.Tx) (any, error)
	Update(ctx context.Context, id int64, TransferBalance entity.TransferBalance, tx pgx.Tx) error
	Update2(ctx context.Context, id int64, TransferBalance entity.TransferBalance, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type TransferBalanceRepository interface {
	TransferBalanceQueryRepository
	TransferBalanceCommandRepository
}
