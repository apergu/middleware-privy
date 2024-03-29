package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type CustomerUsecaseProperty struct {
	CustomerRepo  repository.CustomerRepository
	CustomerPrivy credential.Credential
}

type CustomerQueryUsecase interface {
	Find(ctx context.Context, filter repository.CustomerFilter, limit, skip int64) ([]entity.Customer, interface{}, error)
	Count(ctx context.Context, filter repository.CustomerFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Customer, interface{}, error)
}

type CustomerCommandUsecase interface {
	Create(ctx context.Context, cust model.Customer) (int64, interface{}, error)
	CreateLead(ctx context.Context, cust model.Lead) (int64, interface{}, error)
	CreateLead2(ctx context.Context, cust model.Customer) (int64, interface{}, error)
	UpdateLead(ctx context.Context, id string, cust model.Lead) (any, interface{}, error)
	UpdateLead2(ctx context.Context, id int64, cust model.Lead) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Customer) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
