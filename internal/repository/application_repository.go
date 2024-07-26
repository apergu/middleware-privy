package repository

import (
	"context"

	"middleware/internal/entity"

	"github.com/jackc/pgx/v5"
)

type ApplicationFilter struct {
	Sort          string
	ApplicationID *string
}

type ApplicationQueryRepository interface {
	Find(ctx context.Context, filter ApplicationFilter, limit, skip int64, tx pgx.Tx) ([]entity.Application, error)
	Count(ctx context.Context, filter ApplicationFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Application, error)
	FindByMerchantID(ctx context.Context, merchantID string, tx pgx.Tx) (entity.Application, error)
}

type ApplicationCommandRepository interface {
	Command
	FindByMerchantID(ctx context.Context, merchantID string, tx pgx.Tx) (entity.Application, error)
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Application, error)
	FindByApplicationID(ctx context.Context, enterprisePrivyID string, tx pgx.Tx) (entity.ApplicationFind, error)
	FindByName(ctx context.Context, ApplicationName string, tx pgx.Tx) (entity.ApplicationFind, error)
	Create(ctx context.Context, Application entity.Application, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Application entity.Application, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type ApplicationRepository interface {
	ApplicationQueryRepository
	ApplicationCommandRepository
}
