package usecase

import (
	"context"
	"middleware/pkg/credential"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type SalesOrderCommandUsecaseGeneral struct {
	orderRepo  repository.SalesOrderHeaderCommandRepository
	lineRepo   repository.SalesOrderLineCommandRepository
	orderPrivy credential.SalesOrder
}

func NewSalesOrderCommandUsecaseGeneral(prop SalesOrderUsecaseProperty) *SalesOrderCommandUsecaseGeneral {
	return &SalesOrderCommandUsecaseGeneral{
		orderRepo: prop.SalesOrderHeaderRepo,
		lineRepo:  prop.SalesOrderLineRepo,
	}
}

func (r *SalesOrderCommandUsecaseGeneral) deleteDetail(ctx context.Context, orderId int64, tx pgx.Tx) error {
	err := r.lineRepo.DeleteByHeader(ctx, orderId, tx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.deleteDetail",
				"src":   "lineRepo.DeleteByHeader",
				"param": orderId,
			}).
			Error(err)

		return err
	}

	return nil
}

func (r *SalesOrderCommandUsecaseGeneral) insertDetail(ctx context.Context, lines []model.SalesOrderLine, orderId, createdBy, tm int64, tx pgx.Tx) error {
	for _, line := range lines {
		lineIns := entity.SalesOrderLine{
			SalesOrderHeaderId: orderId,
			ProductID:          line.ProductID,
			ProductName:        line.ProductName,
			Quantity:           line.Quantity,
			RateItem:           line.RateItem,
			TaxRate:            line.TaxRate,
			Subtotal:           line.RateItem * float64(line.Quantity),
			Grandtotal:         line.RateItem*float64(line.Quantity) + line.TaxRate,
			CreatedBy:          createdBy,
			CreatedAt:          tm,
			UpdatedBy:          createdBy,
			UpdatedAt:          tm,
		}

		_, err := r.lineRepo.Create(ctx, lineIns, tx)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"at":    "SalesOrderCommandUsecaseGeneral.insertDetail",
					"src":   "lineRepo.Create",
					"param": lineIns,
				}).
				Error(err)

			return err
		}
	}

	return nil
}

func (r *SalesOrderCommandUsecaseGeneral) Create(ctx context.Context, order model.SalesOrder) (int64, interface{}, error) {
	tx, err := r.orderRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	//tmNow := time.Now().UnixNano() / 1000000

	//var subtotal float64
	//var taxes float64

	//for _, line := range order.Lines {
	//	subtotal += line.RateItem * float64(line.Quantity)
	//	taxes += line.TaxRate * float64(line.Quantity)
	//}
	//grandtotal := subtotal + taxes
	//
	//insertSalesOrder := entity.SalesOrderHeader{
	//	OrderNumber:  order.OrderNumber,
	//	CustomerID:   order.CustomerID,
	//	CustomerName: order.CustomerName,
	//	Subtotal:     subtotal,
	//	Tax:          taxes,
	//	Grandtotal:   grandtotal,
	//	CreatedBy:    order.CreatedBy,
	//	CreatedAt:    tmNow,
	//	UpdatedBy:    order.CreatedBy,
	//	UpdatedAt:    tmNow,
	//}

	insertSalesOrder := entity.SalesOrder{
		Entity:      order.Entity,
		TranDate:    order.TranDate,
		OrderStatus: order.OrderStatus,
		StartDate:   order.StartDate,
		EndDate:     order.EndDate,
		Memo:        order.Memo,
		CustBody2:   order.CustBody2,
	}

	orderId, err := r.orderRepo.Create(ctx, insertSalesOrder, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Create",
				"src":   "orderRepo.Create",
				"param": insertSalesOrder,
			}).
			Error(err)

		return 0, nil, err
	}

	//err = r.insertDetail(ctx, order.Lines, orderId, order.CreatedBy, tmNow, tx)
	//if err != nil {
	//	r.orderRepo.RollbackTx(ctx, tx)
	//
	//	logrus.
	//		WithFields(logrus.Fields{
	//			"at":  "SalesOrderCommandUsecaseGeneral.Create",
	//			"src": "SalesOrderCommandUsecaseGeneral.insertDetail",
	//		}).
	//		Error(err)
	//
	//	return 0, nil, err
	//}

	custPrivyUsgParam := credential.SalesOrderParams{
		RecordType:  "salesord",
		Entity:      order.Entity,
		TranDate:    order.TranDate,
		OrderStatus: order.OrderStatus,
		StartDate:   order.StartDate,
		EndDate:     order.EndDate,
		Memo:        order.Memo,
		CustBody2:   order.CustBody2,
	}

	_, err = r.orderPrivy.CreateSalesOrder(ctx, custPrivyUsgParam)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomerUsage",
				"param": custPrivyUsgParam,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.orderRepo.CommitTx(ctx, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "SalesOrderCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"SalesOrderCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return orderId, nil, nil
}

func (r *SalesOrderCommandUsecaseGeneral) Update(ctx context.Context, id int64, order model.SalesOrderHeader) (int64, interface{}, error) {
	tx, err := r.orderRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	var subtotal float64
	var taxes float64

	for _, line := range order.Lines {
		subtotal += line.RateItem * float64(line.Quantity)
		taxes += line.TaxRate * float64(line.Quantity)
	}
	grandtotal := subtotal + taxes

	updatedSalesOrder := entity.SalesOrderHeader{
		OrderNumber:  order.OrderNumber,
		CustomerID:   order.CustomerID,
		CustomerName: order.CustomerName,
		Subtotal:     subtotal,
		Tax:          taxes,
		Grandtotal:   grandtotal,
		CreatedBy:    order.CreatedBy,
		CreatedAt:    tmNow,
		UpdatedBy:    order.CreatedBy,
		UpdatedAt:    tmNow,
	}

	err = r.orderRepo.Update(ctx, id, updatedSalesOrder, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.deleteDetail(ctx, id, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Update",
				"src":   "r.deleteDetail",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.insertDetail(ctx, order.Lines, id, order.CreatedBy, tmNow, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "SalesOrderCommandUsecaseGeneral.Create",
				"src": "SalesOrderCommandUsecaseGeneral.insertDetail",
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.orderRepo.CommitTx(ctx, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "SalesOrderCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"SalesOrderCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *SalesOrderCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.orderRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.deleteDetail(ctx, id, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Update",
				"src":   "r.deleteDetail",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.orderRepo.Delete(ctx, id, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.orderRepo.CommitTx(ctx, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "SalesOrderCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"SalesOrderCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
