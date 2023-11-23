package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type DivissionQueryUsecaseGeneral struct {
	custRepo repository.DivissionQueryRepository
}

func NewDivissionQueryUsecaseGeneral(prop DivissionUsecaseProperty) *DivissionQueryUsecaseGeneral {
	return &DivissionQueryUsecaseGeneral{
		custRepo: prop.DivissionRepo,
	}
}

func (r *DivissionQueryUsecaseGeneral) Find(ctx context.Context, filter repository.DivissionFilter, limit, skip int64) ([]entity.Divission, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionQueryUsecaseGeneral.Find",
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
				"at":    "DivissionQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *DivissionQueryUsecaseGeneral) Count(ctx context.Context, filter repository.DivissionFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *DivissionQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Divission, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Divission{}, nil, err
	}

	return cust, nil, nil
}
