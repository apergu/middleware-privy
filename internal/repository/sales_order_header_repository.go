package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type SalesOrderHeaderFilter struct {
	Sort string
}

type SalesOrderHeaderQueryRepository interface {
	Find(ctx context.Context, filter SalesOrderHeaderFilter, limit, skip int64, tx pgx.Tx) ([]entity.SalesOrderHeader, error)
	Count(ctx context.Context, filter SalesOrderHeaderFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderHeader, error)
}

type SalesOrderHeaderCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.SalesOrderHeader, error)
	Create(ctx context.Context, SalesOrderHeader entity.SalesOrderHeader, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, SalesOrderHeader entity.SalesOrderHeader, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type SalesOrderHeaderRepository interface {
	SalesOrderHeaderQueryRepository
	SalesOrderHeaderCommandRepository
}
