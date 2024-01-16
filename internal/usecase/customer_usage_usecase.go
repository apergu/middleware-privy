package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type CustomerUsageUsecaseProperty struct {
	CustomerUsageRepo repository.CustomerUsageRepository
	CustomerPrivy     credential.Credential
	MerchantRepo      repository.MerchantQueryRepository
	CustRepo          repository.CustomerQueryRepository
}

type CustomerUsageQueryUsecase interface {
	Find(ctx context.Context, filter repository.CustomerUsageFilter, limit, skip int64) ([]entity.CustomerUsage, interface{}, error)
	Count(ctx context.Context, filter repository.CustomerUsageFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.CustomerUsage, interface{}, error)
}

type CustomerUsageCommandUsecase interface {
	Create(ctx context.Context, cust model.CustomerUsage) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.CustomerUsage) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
