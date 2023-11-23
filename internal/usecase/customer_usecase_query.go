package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type CustomerQueryUsecaseGeneral struct {
	custRepo repository.CustomerQueryRepository
}

func NewCustomerQueryUsecaseGeneral(prop CustomerUsecaseProperty) *CustomerQueryUsecaseGeneral {
	return &CustomerQueryUsecaseGeneral{
		custRepo: prop.CustomerRepo,
	}
}

func (r *CustomerQueryUsecaseGeneral) Find(ctx context.Context, filter repository.CustomerFilter, limit, skip int64) ([]entity.Customer, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerQueryUsecaseGeneral.Find",
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

func (r *CustomerQueryUsecaseGeneral) Count(ctx context.Context, filter repository.CustomerFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *CustomerQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Customer, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Customer{}, nil, err
	}

	return cust, nil, nil
}
