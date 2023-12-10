package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/privy"
)

type TopUpDataUsecaseProperty struct {
	TopUpDataRepo  repository.TopUpDataRepository
	TopUpDataPrivy privy.Privy
	CustomerRepo   repository.CustomerQueryRepository
	MerchantRepo   repository.MerchantQueryRepository
	ChannelRepo    repository.ChannelQueryRepository
}

type TopUpDataQueryUsecase interface {
	Find(ctx context.Context, filter repository.TopUpDataFilter, limit, skip int64) ([]entity.TopUpData, interface{}, error)
	Count(ctx context.Context, filter repository.TopUpDataFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.TopUpData, interface{}, error)
}

type TopUpDataCommandUsecase interface {
	Create(ctx context.Context, cust model.TopUp) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.TopUpData) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
