package usecase

import (
	"context"

	"middleware/internal/model"
	"middleware/internal/repository"
)

type SalesOrderUsecaseProperty struct {
	SalesOrderHeaderRepo repository.SalesOrderHeaderRepository
	SalesOrderLineRepo   repository.SalesOrderLineRepository
}

type SalesOrderQueryUsecase interface {
	Find(ctx context.Context, filter repository.SalesOrderHeaderFilter, limit, skip int64) ([]model.SalesOrderHeaderResponse, interface{}, error)
	Count(ctx context.Context, filter repository.SalesOrderHeaderFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (model.SalesOrderHeaderResponse, interface{}, error)
}

type SalesOrderCommandUsecase interface {
	Create(ctx context.Context, cust model.SalesOrderHeader) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.SalesOrderHeader) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
