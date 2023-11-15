package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
)

type LeadQueryUsecaseGeneral struct {
	custRepo repository.LeadQueryRepository
}

func NewLeadQueryUsecaseGeneral(prop LeadUsecaseProperty) *LeadQueryUsecaseGeneral {
	return &LeadQueryUsecaseGeneral{
		custRepo: prop.LeadRepo,
	}
}

func (r *LeadQueryUsecaseGeneral) Find(ctx context.Context, filter repository.LeadFilter, limit, skip int64) ([]entity.Leads, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "LeadQueryUsecaseGeneral.Find",
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
				"at":    "CustomerQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *LeadQueryUsecaseGeneral) Count(ctx context.Context, filter repository.LeadFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "LeadQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *LeadQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Leads, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "LeadQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Leads{}, nil, err
	}

	return cust, nil, nil
}
