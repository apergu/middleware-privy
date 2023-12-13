package usecase

import (
	"context"
	"log"
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
	orderUsecase := &SalesOrderCommandUsecaseGeneral{
		orderRepo:  prop.SalesOrderHeaderRepo,
		lineRepo:   prop.SalesOrderLineRepo,
		orderPrivy: prop.SalesOrderPrivy,
	}

	if orderUsecase.orderPrivy == nil {
		log.Println("Warning: orderPrivy is nil in NewSalesOrderCommandUsecaseGeneral")
	}

	log.Println("orderUsecase ", orderUsecase.orderPrivy)
	return orderUsecase
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

func (r *SalesOrderCommandUsecaseGeneral) insertDetail(ctx context.Context, lines []model.SalesOrderLines, orderId, createdBy, tm int64, tx pgx.Tx) error {
	for _, line := range lines {
		lineIns := entity.SalesOrderLine{
			ID:                 0,
			SalesOrderHeaderId: 0,
			ProductID:          "",
			ProductName:        line.Item,
			Quantity:           0,
			RateItem:           0,
			TaxRate:            0,
			Subtotal:           0,
			Grandtotal:         0,
			CreatedBy:          0,
			CreatedAt:          0,
			UpdatedBy:          0,
			UpdatedAt:          0,
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

	//insertOrder := entity.SalesOrder{
	//	Entity:      order.Entity,
	//	TranDate:    order.TranDate,
	//	OrderStatus: order.OrderStatus,
	//	StartDate:   order.StartDate,
	//	EndDate:     order.EndDate,
	//	Memo:        order.Memo,
	//	CustBody2:   order.CustBody2,
	//}
	//
	//orderId, err := r.orderRepo.Create(ctx, insertOrder, tx)
	//if err != nil {
	//	r.orderRepo.RollbackTx(ctx, tx)
	//
	//	logrus.
	//		WithFields(logrus.Fields{
	//			"at":    "SalesOrderCommandUsecaseGeneral.Create",
	//			"src":   "orderRepo.Create",
	//			"param": insertOrder,
	//		}).
	//		Error(err)
	//
	//	return 0, nil, err
	//}

	// Find customer by order.Entity
	//customerFilter := repository.CustomerFilter{
	//	EnterprisePrivyID: &order.Entity,
	//}
	//customers, _ := r.custRepo.Find(ctx, customerFilter, 1, 0, nil)
	//
	//var customer entity.Customer
	//if len(customers) > 0 {
	//	customer = customers[0]
	//}

	privyParam := credential.SalesOrderParams{
		RecordType:  "salesord",
		Entity:      order.Entity,
		TranDate:    order.TranDate,
		OrderStatus: order.OrderStatus,
		StartDate:   order.StartDate,
		EndDate:     order.EndDate,
		Memo:        order.Memo,
		CustBody2:   order.CustBody2,
		//CustRecordCustomerName: strconv.Itoa(int(customer.CustomerInternalID)),
		// Add other parameters as needed
	}

	resp, err := r.orderPrivy.CreateSalesOrder(ctx, privyParam)
	log.Println("RESP ", resp)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "SalesOrderCommandUsecaseGeneral.Create",
				"src":   "orderPrivy.CreateSalesOrder",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	//insertOrder.ID = resp.Data.RecordID
	//insertOrder. = customer.CustomerInternalID

	//err = r.orderRepo.Update(ctx, orderId, insertOrder, tx)
	//if err != nil {
	//	r.orderRepo.RollbackTx(ctx, tx)
	//
	//	logrus.
	//		WithFields(logrus.Fields{
	//			"at":    "SalesOrderCommandUsecaseGeneral.Create",
	//			"src":   "orderRepo.Update",
	//			"param": insertOrder,
	//		}).
	//		Error(err)
	//
	//	return 0, nil, err
	//}

	err = r.orderRepo.CommitTx(ctx, tx)
	if err != nil {
		r.orderRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "SalesOrderCommandUsecaseGeneral.Create",
				"src": "orderRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"SalesOrderCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return 12, nil, nil
}

func (r *SalesOrderCommandUsecaseGeneral) Update(ctx context.Context, id int64, order model.SalesOrder) (int64, interface{}, error) {
	tx, err := r.orderRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	//var subtotal float64
	//var taxes float64

	//for _, line := range order.Lines {
	//	subtotal += line.RateItem * float64(line.Quantity)
	//	taxes += line.TaxRate * float64(line.Quantity)
	//}
	//grandtotal := subtotal + taxes

	updatedSalesOrder := entity.SalesOrder{
		ID:          0,
		Entity:      "",
		TranDate:    "",
		OrderStatus: "",
		StartDate:   "",
		EndDate:     "",
		Memo:        "",
		CustBody2:   "",
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

	err = r.insertDetail(ctx, order.Lines, id, 12, tmNow, tx)
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
