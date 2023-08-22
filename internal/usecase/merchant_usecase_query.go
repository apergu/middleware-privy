package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
)

type MerchantQueryUsecaseGeneral struct {
	custRepo repository.MerchantQueryRepository
}

func NewMerchantQueryUsecaseGeneral(prop MerchantUsecaseProperty) *MerchantQueryUsecaseGeneral {
	return &MerchantQueryUsecaseGeneral{
		custRepo: prop.MerchantRepo,
	}
}

func (r *MerchantQueryUsecaseGeneral) Find(ctx context.Context, filter repository.MerchantFilter, limit, skip int64) ([]entity.Merchant, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantQueryUsecaseGeneral.Find",
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
				"at":    "MerchantQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *MerchantQueryUsecaseGeneral) Count(ctx context.Context, filter repository.MerchantFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *MerchantQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Merchant, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Merchant{}, nil, err
	}

	return cust, nil, nil
}
