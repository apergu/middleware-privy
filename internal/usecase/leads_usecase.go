package usecase

import (
	"context"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
)

type LeadUsecaseProperty struct {
	LeadRepo  repository.LeadRepository
	LeadPrivy credential.Credential
}

type LeadQueryUsecase interface {
	Find(ctx context.Context, filter repository.LeadFilter, limit, skip int64) ([]entity.Leads, interface{}, error)
	Count(ctx context.Context, filter repository.LeadFilter) (int64, interface{}, error)
	FindById(ctx context.Context, id int64) (entity.Leads, interface{}, error)
}

type LeadCommandUsecase interface {
	Create(ctx context.Context, cust model.Leads) (int64, interface{}, error)
	//Update(ctx context.Context, id int64, cust model.Leads) (int64, interface{}, error)
	//Delete(ctx context.Context, id int64) (int64, interface{}, error)
}
