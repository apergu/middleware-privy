package usecase

import (
	"context"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/sirupsen/logrus"
)

type SalesOrderQueryUsecaseGeneral struct {
	salesOrderRepo     repository.SalesOrderHeaderQueryRepository
	salesOrderLineRepo repository.SalesOrderLineQueryRepository
}

func NewSalesOrderQueryUsecaseGeneral(prop SalesOrderUsecaseProperty) *SalesOrderQueryUsecaseGeneral {
	return &SalesOrderQueryUsecaseGeneral{
		salesOrderRepo:     prop.SalesOrderHeaderRepo,
		salesOrderLineRepo: prop.SalesOrderLineRepo,
	}
}

func (r *SalesOrderQueryUsecaseGeneral) findDetail(ctx context.Context, headerId int64) ([]entity.SalesOrderLine, error) {
	filter := repository.SalesOrderLineFilter{
		HeaderId: headerId,
	}

	lines, _ := r.salesOrderLineRepo.Find(ctx, filter, 0, 0, nil)

	return lines, nil
}

func (r *SalesOrderQueryUsecaseGeneral) Find(ctx context.Context, filter repository.SalesOrderHeaderFilter, limit, skip int64) ([]model.SalesOrderResponse, interface{}, error) {
	orders, err := r.salesOrderRepo.Find(ctx, filter, limit, skip, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderQueryUsecaseGeneral.Find",
				"src":   "custRepo.Find",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	resp := make([]model.SalesOrderResponse, len(orders))
	var lastId int64

	//for i, order := range orders {
	//	lines, err := r.findDetail(ctx, order.ID)
	//	if err != nil {
	//		logrus.
	//			WithFields(logrus.Fields{
	//				"at":    "SalesOrderQueryUsecaseGeneral.Find",
	//				"src":   "r.findDetail",
	//				"param": order.ID,
	//			}).
	//			Error(err)
	//
	//		return nil, nil, err
	//	}
	//
	//	resp[i] = model.s{
	//		SalesOrderHeader: order,
	//		Lines:            lines,
	//	}
	//
	//	lastId = order.ID
	//}

	count, err := r.salesOrderRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderQueryUsecaseGeneral.Find",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return nil, nil, err
	}

	return resp, model.NewMeta(count, limit, lastId), nil
}

func (r *SalesOrderQueryUsecaseGeneral) Count(ctx context.Context, filter repository.SalesOrderHeaderFilter) (int64, interface{}, error) {
	count, err := r.salesOrderRepo.Count(ctx, filter, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderQueryUsecaseGeneral.Count",
				"src":   "custRepo.Count",
				"param": filter,
			}).
			Error(err)

		return 0, nil, err
	}

	return count, nil, nil
}

func (r *SalesOrderQueryUsecaseGeneral) FindById(ctx context.Context, id int64) (model.SalesOrderResponse, interface{}, error) {
	order, err := r.salesOrderRepo.FindOneById(ctx, id, nil)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderQueryUsecaseGeneral.FindById",
				"src":   "custRepo.FindOneById",
				"param": id,
			}).
			Error(err)

		return model.SalesOrderResponse{}, nil, err
	}

	//lines, err := r.findDetail(ctx, order.ID)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderQueryUsecaseGeneral.Find",
				"src":   "r.findDetail",
				"param": order.ID,
			}).
			Error(err)

		return model.SalesOrderResponse{}, nil, err
	}

	resp := model.SalesOrderResponse{
		SalesOrder: entity.SalesOrder{},
		Lines:      nil,
	}

	return resp, nil, nil
}
