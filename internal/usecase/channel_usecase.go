package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type ChannelUsecaseProperty struct {
	ChannelRepo  repository.ChannelRepository
	ChannelPrivy credential.Credential
	MerchantRepo repository.MerchantRepository
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
