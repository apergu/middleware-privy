package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type ChannelQueryUsecaseGeneral struct {
	custRepo repository.ChannelQueryRepository
}

func NewChannelQueryUsecaseGeneral(prop ChannelUsecaseProperty) *ChannelQueryUsecaseGeneral {
	return &ChannelQueryUsecaseGeneral{
		custRepo: prop.ChannelRepo,
	}
}

func (r *ChannelQueryUsecaseGeneral) Find(ctx context.Context, filter repository.ChannelFilter, limit, skip int64) ([]entity.Channel, interface{}, error) {
	customers, err := r.custRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelQueryUsecaseGeneral.Find",
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
				"at":    "ChannelQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return customers, model.NewMeta(count, limit, lastId), nil
}

func (r *ChannelQueryUsecaseGeneral) Count(ctx context.Context, filter repository.ChannelFilter) (int64, interface{}, error) {
	count, err := r.custRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *ChannelQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (entity.Channel, interface{}, error) {
	cust, err := r.custRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return entity.Channel{}, nil, err
	}

	return cust, nil, nil
}
