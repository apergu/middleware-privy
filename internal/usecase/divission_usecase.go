package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type DivissionUsecaseProperty struct {
	DivissionRepo  repository.DivissionRepository
	DivissionPrivy credential.Credential
}

type DivissionQueryUsecase interface {
	Find(ctx context.Context, filter repository.DivissionFilter, limit, skip int64) ([]entity.Divission, interface{}, error)
	Count(ctx context.Context, filter repository.DivissionFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Divission, interface{}, error)
}

type DivissionCommandUsecase interface {
	Create(ctx context.Context, cust model.Divission) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Divission) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
