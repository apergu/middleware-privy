package usecase

import (
	"context"
	"log"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type TransferBalanceCommandUsecaseGeneral struct {
	merchantRepo  repository.TransferBalanceCommandRepository
	custRepo      repository.CustomerQueryRepository
	merchantPrivy credential.TransferBalance
}

func NewTransferBalanceCommandUsecaseGeneral(prop TransferBalanceUsecaseProperty) *TransferBalanceCommandUsecaseGeneral {
	return &TransferBalanceCommandUsecaseGeneral{
		merchantRepo:  prop.TransferBalanceRepo,
		custRepo:      prop.CustomerRepo,
		merchantPrivy: prop.TransferBalancePrivy,
	}
}

func (r *TransferBalanceCommandUsecaseGeneral) Create(ctx context.Context, merchant model.TransferBalance) (any, interface{}, error) {
	tx, err := r.merchantRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertTransferBalance := entity.TransferBalance{
		CustomerId:   merchant.CustomerId,
		TransferDate: merchant.TransferDate,
		TrxIdFrom:    merchant.TrxIdFrom,
		TrxIdTo:      merchant.TrxIdTo,
		MerchantTo:   merchant.MerchantTo,
		ChannelTo:    merchant.ChannelTo,
		StartDate:    merchant.StartDate,
		EndDate:      merchant.EndDate,
		IsTrxCreated: merchant.IsTrxCreated,
		Quantity:     merchant.Quantity,
		CreatedBy:    merchant.CreatedBy,
		CreatedAt:    tmNow,
		UpdatedBy:    merchant.CreatedBy,
		UpdatedAt:    tmNow,
	}

	// merchantId, err := r.merchantRepo.Create(ctx, insertTransferBalance, tx)
	// if err != nil {
	// 	r.merchantRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "TransferBalanceCommandUsecaseGeneral.Create",
	// 			"src":   "custRepo.Create",
	// 			"param": insertTransferBalance,
	// 		}).
	// 		Error(err)

	// 	return 0, nil, err
	// }

	// find customer by merchant.EnterpriseID
	// customer_filter := repository.CustomerFilter{
	// 	EnterprisePrivyID: &merchant.EnterpriseID,
	// }
	// customers, _ := r.custRepo.Find(ctx, customer_filter, 1, 0, nil)

	// var customer entity.Customer
	// if len(customers) > 0 {
	// 	customer = customers[0]
	// }

	// custrecordcustomer_name ambil dari customer

	privyParam := credential.TransferBalanceParam{
		RecordType:                 "customrecord_transfer_balance",
		CustRecordCustomer:         merchant.CustomerId,
		CustRecordTransferDate:     merchant.TransferDate,
		CustRecordTrxNoFrom:        merchant.TrxIdFrom,
		CustRecordTrxNoTo:          merchant.TrxIdTo,
		CustRecordMerchant:         merchant.MerchantTo,
		CustRecordChannel:          merchant.ChannelTo,
		CustRecordStartDateLayanan: merchant.StartDate,
		CustRecordEndDateLayanan:   merchant.EndDate,
		CustRecordIsTrxIdCreated:   merchant.IsTrxCreated,
		CustRecordFromQuantity:     merchant.Quantity,
	}

	log.Println("PRIVY PARAM ", privyParam)
	log.Println("PRIVY PARAM 2", r.merchantPrivy)
	resp, err := r.merchantPrivy.CreateTransferBalance(ctx, privyParam)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceCommandUsecaseGeneral.Create",
				"src":   "merchantPrivy.CreateTransferBalance",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	// insertTransferBalance.InternalId = resp.Data.RecordID

	// err = r.merchantRepo.Update(ctx, merchantId, insertTransferBalance, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": insertTransferBalance,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.merchantRepo.CommitTx(ctx, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TransferBalanceCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TransferBalanceCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return resp.Data.RecordID, nil, nil
}

func (r *TransferBalanceCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.TransferBalance) (int64, interface{}, error) {
	tx, err := r.merchantRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedTransferBalance := entity.TransferBalance{
		CustomerId:   merchant.CustomerId,
		TransferDate: merchant.TransferDate,
		TrxIdFrom:    merchant.TrxIdFrom,
		TrxIdTo:      merchant.TrxIdTo,
		MerchantTo:   merchant.MerchantTo,
		ChannelTo:    merchant.ChannelTo,
		StartDate:    merchant.StartDate,
		EndDate:      merchant.EndDate,
		IsTrxCreated: merchant.IsTrxCreated,
		Quantity:     merchant.Quantity,
		CreatedBy:    merchant.CreatedBy,
		UpdatedBy:    merchant.CreatedBy,
		UpdatedAt:    tmNow,
	}

	err = r.merchantRepo.Update(ctx, id, updatedTransferBalance, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.merchantRepo.CommitTx(ctx, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TransferBalanceCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TransferBalanceCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *TransferBalanceCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.merchantRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.merchantRepo.Delete(ctx, id, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TransferBalanceCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.merchantRepo.CommitTx(ctx, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TransferBalanceCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TransferBalanceCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
