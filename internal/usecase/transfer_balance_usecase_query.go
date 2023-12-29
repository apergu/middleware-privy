package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type TransferBalanceQueryUsecaseGeneral struct {
	custRepo repository.TransferBalanceQueryRepository
}

func NewTransferBalanceQueryUsecaseGeneral(prop TransferBalanceUsecaseProperty) *TransferBalanceQueryUsecaseGeneral {
	return &TransferBalanceQueryUsecaseGeneral{
		custRepo: prop.TransferBalanceRepo,
	}
}

func (r *TransferBalanceQueryUsecaseGeneral) Find(ctx context.Context, filter repository.TransferBalanceFilter, limit, skip int64) ([]entity.TransferBalance, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceQueryUsecaseGeneral.Find",
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
				"at":    "TransferBalanceQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *TransferBalanceQueryUsecaseGeneral) Count(ctx context.Context, filter repository.TransferBalanceFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *TransferBalanceQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.TransferBalance, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.TransferBalance{}, nil, err
	}

	return cust, nil, nil
}
