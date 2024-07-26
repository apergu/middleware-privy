package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type ApplicationQueryUsecaseGeneral struct {
	custRepo repository.ApplicationQueryRepository
}

func NewApplicationQueryUsecaseGeneral(prop ApplicationUsecaseProperty) *ApplicationQueryUsecaseGeneral {
	return &ApplicationQueryUsecaseGeneral{
		custRepo: prop.ApplicationRepo,
	}
}

func (r *ApplicationQueryUsecaseGeneral) Find(ctx context.Context, filter repository.ApplicationFilter, limit, skip int64) ([]entity.Application, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationQueryUsecaseGeneral.Find",
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
				"at":    "ApplicationQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *ApplicationQueryUsecaseGeneral) Count(ctx context.Context, filter repository.ApplicationFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *ApplicationQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Application, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Application{}, nil, err
	}

	return cust, nil, nil
}
