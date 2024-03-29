package repository

import (
	"context"

	"middleware/internal/entity"

	"github.com/jackc/pgx/v5"
)

type SalesOrderHeaderFilter struct {
	Sort string
}

type SalesOrderHeaderQueryRepository interface {
	Find(ctx context.Context, filter SalesOrderHeaderFilter, limit, skip int64, tx pgx.Tx) ([]entity.SalesOrder, error)
	Count(ctx context.Context, filter SalesOrderHeaderFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrder, error)
}

type SalesOrderHeaderCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrder, error)
	Create(ctx context.Context, SalesOrderHeader entity.SalesOrder, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, SalesOrderHeader entity.SalesOrder, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type SalesOrderHeaderRepository interface {
	SalesOrderHeaderQueryRepository
	SalesOrderHeaderCommandRepository
}
