package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type DivissionFilter struct {
	Sort string
}

type DivissionQueryRepository interface {
	Find(ctx context.Context, filter DivissionFilter, limit, skip int64, tx pgx.Tx) ([]entity.Divission, error)
	Count(ctx context.Context, filter DivissionFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Divission, error)
}

type DivissionCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Divission, error)
	Create(ctx context.Context, Divission entity.Divission, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Divission entity.Divission, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type DivissionRepository interface {
	DivissionQueryRepository
	DivissionCommandRepository
}
