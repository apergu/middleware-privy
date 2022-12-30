package usecase

import (
	"context"

	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
)

type CustomerUsecaseProperty struct {
	CustomerRepo repository.CustomerRepository
}

type CustomerQueryUsecase interface {
	Find(ctx context.Context, filter repository.CustomerFilter, limit, skip int64) ([]entity.Customer, interface{}, error)
	Count(ctx context.Context, filter repository.CustomerFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Customer, interface{}, error)
}

type CustomerCommandUsecase interface {
	Create(ctx context.Context, cust model.Customer) (int64, interface{}, error)
	Update(ctx context.Context, id int64, cust model.Customer) (int64, interface{}, error)
	Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
