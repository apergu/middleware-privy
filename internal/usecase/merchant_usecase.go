package usecase

import (
	"context"

	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
)

type MerchantUsecaseProperty struct {
	MerchantRepo  repository.MerchantRepository
	MerchantPrivy credential.Credential
}

type MerchantQueryUsecase interface {
	Find(ctx context.Context, filter repository.MerchantFilter, limit, skip int64) ([]entity.Merchant, interface{}, error)
	Count(ctx context.Context, filter repository.MerchantFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Merchant, interface{}, error)
}

type MerchantCommandUsecase interface {
	Create(ctx context.Context, cust model.Merchant) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Merchant) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
