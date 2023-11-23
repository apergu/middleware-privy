package repository

import (
	"context"

	"middleware/internal/entity"

	"github.com/jackc/pgx/v5"
)

type SalesOrderLineFilter struct {
	HeaderId int64
	Sort     string
}

type SalesOrderLineQueryRepository interface {
	Find(ctx context.Context, filter SalesOrderLineFilter, limit, skip int64, tx pgx.Tx) ([]entity.SalesOrderLine, error)
	Count(ctx context.Context, filter SalesOrderLineFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderLine, error)
}

type SalesOrderLineCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderLine, error)
	Create(ctx context.Context, SalesOrderLine entity.SalesOrderLine, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, SalesOrderLine entity.SalesOrderLine, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
	DeleteByHeader(ctx context.Context, headerId int64, tx pgx.Tx) error
}

type SalesOrderLineRepository interface {
	SalesOrderLineQueryRepository
	SalesOrderLineCommandRepository
}
