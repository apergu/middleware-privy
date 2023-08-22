package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type TopUpDataFilter struct {
	Sort string
}

type TopUpDataQueryRepository interface {
	Find(ctx context.Context, filter TopUpDataFilter, limit, skip int64, tx pgx.Tx) ([]entity.TopUpData, error)
	Count(ctx context.Context, filter TopUpDataFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.TopUpData, error)
}

type TopUpDataCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.TopUpData, error)
	Create(ctx context.Context, TopUpData entity.TopUpData, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, TopUpData entity.TopUpData, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type TopUpDataRepository interface {
	TopUpDataQueryRepository
	TopUpDataCommandRepository
}
