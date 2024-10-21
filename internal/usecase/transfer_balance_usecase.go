package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type TransferBalanceUsecaseProperty struct {
	TransferBalanceRepo  repository.TransferBalanceRepository
	CustomerRepo         repository.CustomerRepository
	TransferBalancePrivy credential.Credential
}

type TransferBalanceQueryUsecase interface {
	Find(ctx context.Context, filter repository.TransferBalanceFilter, limit, skip int64) ([]entity.TransferBalance, interface{}, error)
	Count(ctx context.Context, filter repository.TransferBalanceFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.TransferBalance, interface{}, error)
}

type TransferBalanceCommandUsecase interface {
	Create(ctx context.Context, cust model.TransferBalance) (any, interface{}, error)
	Update(ctx context.Context, id int64, cust model.TransferBalance) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
