package usecase

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"
)

type MerchantCommandUsecaseGeneral struct {
	merchantRepo  repository.MerchantCommandRepository
	custRepo      repository.CustomerQueryRepository
	merchantPrivy credential.Merchant
}

func NewMerchantCommandUsecaseGeneral(prop MerchantUsecaseProperty) *MerchantCommandUsecaseGeneral {
	return &MerchantCommandUsecaseGeneral{
		merchantRepo:  prop.MerchantRepo,
		merchantPrivy: prop.MerchantPrivy,
		custRepo:      prop.CustomerRepo,
	}
}

func (r *MerchantCommandUsecaseGeneral) Create(ctx context.Context, merchant model.Merchant) (int64, interface{}, error) {
	tx, err := r.merchantRepo.BeginTx(ctx)
	log.Println("custRepo", r.custRepo)
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			r.merchantRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.merchantRepo.RollbackTx(ctx, tx)
		} else {
			err = r.merchantRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	tmNow := time.Now().UnixNano() / 1000000

	// respCust, _ := r.merchantRepo.FindByName(ctx, merchant.MerchantName, tx)

	// if respCust.MerchantName != "" {
	// 	return 0, nil, rapperror.ErrConflict(
	// 		"",
	// 		"Merchant with name "+merchant.MerchantName+" already exist",
	// 		"MerchantCommandUsecaseGeneral.Create",
	// 		nil,
	// 	)
	// }

	// defer func() {
	// 	if p := recover(); p != nil {
	// 		r.merchantRepo.RollbackTx(ctx, tx)
	// 		panic(p)
	// 	} else if err != nil {
	// 		log.Println("Rolling back transaction due to error:", err)
	// 		r.merchantRepo.RollbackTx(ctx, tx)
	// 	} else {
	// 		err = r.merchantRepo.CommitTx(ctx, tx)
	// 		if err != nil {
	// 			log.Println("Error committing transaction:", err)
	// 		}
	// 	}
	// }()

	respCust2, _ := r.merchantRepo.FindByMerchantID(ctx, merchant.MerchantID, tx)

	if respCust2.MerchantID != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Merchant with Merchant ID "+merchant.MerchantID+" already exist",
			"MerchantCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.merchantRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.merchantRepo.RollbackTx(ctx, tx)
		} else {
			err = r.merchantRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

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

	merchantId, err := r.merchantRepo.Create(ctx, insertMerchant, tx)
	print("merchantId", merchantId)
	if err != nil {
		r.merchantRepo.CommitTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertMerchant,
			}).
			Error(err)

		data := map[string]interface{}{
			"merchant": insertMerchant,
			"message":  err.Error(),
		}

		return 0, data, nil
	}

	defer func() {
		if p := recover(); p != nil {
			r.merchantRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.merchantRepo.RollbackTx(ctx, tx)
		} else {
			err = r.merchantRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	// find customer by merchant.EnterpriseID
	customer_filter := repository.CustomerFilter{
		EnterprisePrivyID: &merchant.EnterpriseID,
	}
	customers, _ := r.custRepo.Find(ctx, customer_filter, 1, 0, nil)

	defer func() {
		if p := recover(); p != nil {
			r.merchantRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.merchantRepo.CommitTx(ctx, tx)
		} else {
			err = r.merchantRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	var customer entity.Customer
	if len(customers) > 0 {
		customer = customers[0]
	}

	// if customer.CustomerInternalID == 0 {

	// }

	// custrecordcustomer_name ambil dari customer

	privyParam := credential.MerchantParam{
		RecordType:                  "customrecord_customer_hierarchy",
		CustRecordCustomerName:      customer.CustomerInternalID,
		CustRecordEnterpriseID:      merchant.EnterpriseID,
		CustRecordMerchantID:        merchant.MerchantID,
		CustRecordPrivyCodeMerchant: merchant.MerchantCode,
		CustRecordMerchantName:      merchant.MerchantName,
		CustRecordAddress:           merchant.Address,
		CustRecordEmail:             merchant.Email,
		CustRecordPhone:             merchant.PhoneNo,
		CustRecordState:             merchant.State,
		CustRecordCity:              merchant.City,
		CustRecordZip:               merchant.ZipCode,
		Method:                      "POST",
	}

	resp, err := r.merchantPrivy.CreateMerchant(ctx, privyParam)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "merchantPrivy.CreateMerchant",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertMerchant.MerchantInternalID = resp.Data.RecordID
	insertMerchant.CustomerInternalID = customer.CustomerInternalID

	err = r.merchantRepo.Update(ctx, merchantId, insertMerchant, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": insertMerchant,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.merchantRepo.CommitTx(ctx, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

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

	return merchantId, nil, nil
}

func (r *MerchantCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Merchant) (int64, interface{}, error) {
	tx, err := r.merchantRepo.BeginTx(ctx)
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

	err = r.merchantRepo.Update(ctx, id, updatedMerchant, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Update",
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
	tx, err := r.merchantRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.merchantRepo.Delete(ctx, id, tx)
	if err != nil {
		r.merchantRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Delete",
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
