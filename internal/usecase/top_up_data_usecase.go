package usecase

import (
	"context"

	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
)

type TopUpDataUsecaseProperty struct {
	TopUpDataRepo  repository.TopUpDataRepository
	TopUpDataPrivy credential.Credential
}

type TopUpDataQueryUsecase interface {
	Find(ctx context.Context, filter repository.TopUpDataFilter, limit, skip int64) ([]entity.TopUpData, interface{}, error)
	Count(ctx context.Context, filter repository.TopUpDataFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.TopUpData, interface{}, error)
}

type TopUpDataCommandUsecase interface {
	Create(ctx context.Context, cust model.TopUpData) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.TopUpData) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
