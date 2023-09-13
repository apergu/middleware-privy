package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/pkg/credential"
	"gitlab.com/rteja-library3/rapperror"
)

type MerchantCommandUsecaseGeneral struct {
	custRepo      repository.MerchantCommandRepository
	merchantPrivy credential.Merchant
}

func NewMerchantCommandUsecaseGeneral(prop MerchantUsecaseProperty) *MerchantCommandUsecaseGeneral {
	return &MerchantCommandUsecaseGeneral{
		custRepo:      prop.MerchantRepo,
		merchantPrivy: prop.MerchantPrivy,
	}
}

func (r *MerchantCommandUsecaseGeneral) Create(ctx context.Context, merchant model.Merchant) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertMerchant := entity.Merchant{
		CustomerID:   merchant.CustomerID,
		EnterpriseID: merchant.EnterpriseID,
		MerchantCode: merchant.MerchantCode,
		MerchantID:   merchant.MerchantID,
		MerchantName: merchant.MerchantName,
		Address:      merchant.Address,
		Email:        merchant.Email,
		PhoneNo:      merchant.PhoneNo,
		State:        merchant.State,
		City:         merchant.City,
		ZipCode:      merchant.ZipCode,
		CreatedBy:    merchant.CreatedBy,
		CreatedAt:    tmNow,
		UpdatedBy:    merchant.CreatedBy,
		UpdatedAt:    tmNow,
	}

	custId, err := r.custRepo.Create(ctx, insertMerchant, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertMerchant,
			}).
			Error(err)

		return 0, nil, err
	}

	// privyParam := credential.MerchantParam{
	// 	RecordType:             "customrecord_customer_hierarchy",
	// 	CustRecordEnterpriseID: merchant.EnterpriseID,
	// 	CustRecordMerchantID:   merchant.MerchantID,
	// 	CustRecordMerchantName: merchant.MerchantName,
	// 	CustRecordAddress:      merchant.Address,
	// 	CustRecordEmail:        merchant.Email,
	// 	CustRecordPhone:        merchant.PhoneNo,
	// 	CustRecordState:        merchant.State,
	// 	CustRecordCity:         merchant.City,
	// 	CustRecordZip:          merchant.ZipCode,
	// }

	// err = r.merchantPrivy.CreateMerchant(ctx, privyParam)
	// if err != nil {
	// 	r.custRepo.RollbackTx(ctx, tx)

	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"at":    "MerchantCommandUsecaseGeneral.Create",
	// 			"src":   "merchantPrivy.CreateMerchant",
	// 			"param": privyParam,
	// 		}).
	// 		Error(err)

	// 	return 0, nil, err
	// }

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "MerchantCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"MerchantCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return custId, nil, nil
}

func (r *MerchantCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Merchant) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedMerchant := entity.Merchant{
		MerchantID:   merchant.MerchantID,
		MerchantName: merchant.MerchantName,
		Address:      merchant.Address,
		Email:        merchant.Email,
		PhoneNo:      merchant.PhoneNo,
		State:        merchant.State,
		City:         merchant.City,
		ZipCode:      merchant.ZipCode,
		UpdatedBy:    merchant.CreatedBy,
		UpdatedAt:    tmNow,
	}

	err = r.custRepo.Update(ctx, id, updatedMerchant, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "MerchantCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"MerchantCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *MerchantCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "MerchantCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"MerchantCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
