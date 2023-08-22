package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
)

type TopUpDataQueryUsecaseGeneral struct {
	custRepo repository.TopUpDataQueryRepository
}

func NewTopUpDataQueryUsecaseGeneral(prop TopUpDataUsecaseProperty) *TopUpDataQueryUsecaseGeneral {
	return &TopUpDataQueryUsecaseGeneral{
		custRepo: prop.TopUpDataRepo,
	}
}

func (r *TopUpDataQueryUsecaseGeneral) Find(ctx context.Context, filter repository.TopUpDataFilter, limit, skip int64) ([]entity.TopUpData, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataQueryUsecaseGeneral.Find",
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
				"at":    "TopUpDataQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *TopUpDataQueryUsecaseGeneral) Count(ctx context.Context, filter repository.TopUpDataFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *TopUpDataQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.TopUpData, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.TopUpData{}, nil, err
	}

	return cust, nil, nil
}
