package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type CustomerFilter struct {
	Sort string
}

type CustomerQueryRepository interface {
	Find(ctx context.Context, filter CustomerFilter, limit, skip int64, tx pgx.Tx) ([]entity.Customer, error)
	Count(ctx context.Context, filter CustomerFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Customer, error)
}

type CustomerCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Customer, error)
	Create(ctx context.Context, Customer entity.Customer, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Customer entity.Customer, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type CustomerRepository interface {
	CustomerQueryRepository
	CustomerCommandRepository
}
