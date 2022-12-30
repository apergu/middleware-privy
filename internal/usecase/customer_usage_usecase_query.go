package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
)

type CustomerUsageQueryUsecaseGeneral struct {
	custRepo repository.CustomerUsageQueryRepository
}

func NewCustomerUsageQueryUsecaseGeneral(prop CustomerUsageUsecaseProperty) *CustomerUsageQueryUsecaseGeneral {
	return &CustomerUsageQueryUsecaseGeneral{
		custRepo: prop.CustomerUsageRepo,
	}
}

func (r *CustomerUsageQueryUsecaseGeneral) Find(ctx context.Context, filter repository.CustomerUsageFilter, limit, skip int64) ([]entity.CustomerUsage, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageQueryUsecaseGeneral.Find",
				"src":   "custRepo.Find",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	var lastId int64
	if len(customers) > 0 {
		lastId = customers[len(customers)-1].ID
	}

	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *CustomerUsageQueryUsecaseGeneral) Count(ctx context.Context, filter repository.CustomerUsageFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *CustomerUsageQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.CustomerUsage, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.CustomerUsage{}, nil, err
	}

	return cust, nil, nil
}
