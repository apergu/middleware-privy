package usecase

import (
	"context"

	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
)

type ChannelUsecaseProperty struct {
	ChannelRepo  repository.ChannelRepository
	ChannelPrivy credential.Credential
}

type ChannelQueryUsecase interface {
	Find(ctx context.Context, filter repository.ChannelFilter, limit, skip int64) ([]entity.Channel, interface{}, error)
	Count(ctx context.Context, filter repository.ChannelFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Channel, interface{}, error)
}

type ChannelCommandUsecase interface {
	Create(ctx context.Context, cust model.Channel) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Channel) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
