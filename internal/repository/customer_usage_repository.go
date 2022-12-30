package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type CustomerUsageFilter struct {
	Sort string
}

type CustomerUsageQueryRepository interface {
	Find(ctx context.Context, filter CustomerUsageFilter, limit, skip int64, tx pgx.Tx) ([]entity.CustomerUsage, error)
	Count(ctx context.Context, filter CustomerUsageFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.CustomerUsage, error)
}

type CustomerUsageCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.CustomerUsage, error)
	Create(ctx context.Context, CustomerUsage entity.CustomerUsage, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, CustomerUsage entity.CustomerUsage, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type CustomerUsageRepository interface {
	CustomerUsageQueryRepository
	CustomerUsageCommandRepository
}
