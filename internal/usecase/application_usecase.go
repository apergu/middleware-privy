package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type ApplicationUsecaseProperty struct {
	// ApplicationRepo  repository.ApplicationRepository
	// ApplicationPrivy credential.Credential
	// MerchantRepo     repository.MerchantRepository
	CustomerRepo     repository.CustomerRepository
	ApplicationRepo  repository.ApplicationRepository
	ApplicationPrivy credential.Credential
}

type ApplicationQueryUsecase interface {
	Find(ctx context.Context, filter repository.ApplicationFilter, limit, skip int64) ([]entity.Application, interface{}, error)
	Count(ctx context.Context, filter repository.ApplicationFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Application, interface{}, error)
}

type ApplicationCommandUsecase interface {
	Create(ctx context.Context, cust model.Application) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Application) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
