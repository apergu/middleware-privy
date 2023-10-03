package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
)

type ChannelFilter struct {
	Sort      string
	ChannelID *string
}

type ChannelQueryRepository interface {
	Find(ctx context.Context, filter ChannelFilter, limit, skip int64, tx pgx.Tx) ([]entity.Channel, error)
	Count(ctx context.Context, filter ChannelFilter, tx pgx.Tx) (int64, error)
	FindOneById(ctx context.Context, id int64, tx pgx.Tx) (entity.Channel, error)
}

type ChannelCommandRepository interface {
	Command
	FindOneByIdForUpdate(ctx context.Context, id int64, tx pgx.Tx) (entity.Channel, error)
	Create(ctx context.Context, Channel entity.Channel, tx pgx.Tx) (int64, error)
	Update(ctx context.Context, id int64, Channel entity.Channel, tx pgx.Tx) error
	Delete(ctx context.Context, id int64, tx pgx.Tx) error
}

type ChannelRepository interface {
	ChannelQueryRepository
	ChannelCommandRepository
}
