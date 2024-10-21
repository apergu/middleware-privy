package usecase

import (
	"context"
	"middleware/pkg/credential"

	"middleware/internal/model"
	"middleware/internal/repository"
)

type SalesOrderUsecaseProperty struct {
	SalesOrderHeaderRepo repository.SalesOrderHeaderRepository
	SalesOrderLineRepo   repository.SalesOrderLineRepository
	SalesOrderPrivy      credential.Credential
}

type SalesOrderQueryUsecase interface {
	Find(ctx context.Context, filter repository.SalesOrderHeaderFilter, limit, skip int64) ([]model.SalesOrderResponse, interface{}, error)
	Count(ctx context.Context, filter repository.SalesOrderHeaderFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (model.SalesOrderResponse, interface{}, error)
}

type SalesOrderCommandUsecase interface {
	Create(ctx context.Context, cust model.SalesOrder) (any, interface{}, error)
	Update(ctx context.Context, id int64, cust model.SalesOrder) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
