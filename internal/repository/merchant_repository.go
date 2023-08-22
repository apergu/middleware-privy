package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type MerchantFilter struct {
	Sort string
}

type MerchantQueryRepository interface {
	Find(ctx context.Context, filter MerchantFilter, limit, skip int64, tx pgx.Tx) ([]entity.Merchant, error)
	Count(ctx context.Context, filter MerchantFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Merchant, error)
}

type MerchantCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Merchant, error)
	Create(ctx context.Context, Merchant entity.Merchant, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Merchant entity.Merchant, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type MerchantRepository interface {
	MerchantQueryRepository
	MerchantCommandRepository
}
