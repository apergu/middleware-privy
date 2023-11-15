package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type LeadFilter struct {
	Sort              string
	EnterprisePrivyID *string
	CustomerID        *string
}

type LeadQueryRepository interface {
	Find(ctx context.Context, filter LeadFilter, limit, skip int64, tx pgx.Tx) ([]entity.Leads, error)
	Count(ctx context.Context, filter LeadFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Leads, error)
}

type LeadCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Leads, error)
	Create(ctx context.Context, Customer entity.Leads, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Customer entity.Leads, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type LeadRepository interface {
	LeadQueryRepository
	LeadCommandRepository
}
