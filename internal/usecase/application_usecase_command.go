package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ApplicationCommandUsecaseGeneral struct {
	channelRepo  repository.ApplicationCommandRepository
	channelPrivy credential.Application
	merchantRepo repository.ApplicationQueryRepository
}

func NewApplicationCommandUsecaseGeneral(prop ApplicationUsecaseProperty) *ApplicationCommandUsecaseGeneral {
	return &ApplicationCommandUsecaseGeneral{
		channelRepo:  prop.ApplicationRepo,
		channelPrivy: prop.ApplicationPrivy,
		merchantRepo: prop.ApplicationRepo,
	}
}

func (r *ApplicationCommandUsecaseGeneral) Create(ctx context.Context, channelParam model.Application) (int64, interface{}, error) {
	tx, err := r.channelRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			r.channelRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.channelRepo.RollbackTx(ctx, tx)
		} else {
			err = r.channelRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	fmt.Println("merchantsFIND2")
	// tmNow := time.Now().UnixNano() / 1000000

	merchant_filter := repository.ApplicationFilter{
		ApplicationID: &channelParam.ApplicationID,
	}
	merchants, err := r.merchantRepo.Find(ctx, merchant_filter, 1, 0, nil)

	if len(merchants) == 0 {
		return 0, nil, rapperror.ErrUnprocessableEntity(
			"",
			"Application with ID "+channelParam.ApplicationID+" is Not Found",
			"ApplicationCommandUsecaseGeneral.Create",
			nil,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			r.channelRepo.RollbackTx(ctx, tx)
			panic(p)
		} else if err != nil {
			log.Println("Rolling back transaction due to error:", err)
			r.channelRepo.RollbackTx(ctx, tx)
		} else {
			err = r.channelRepo.CommitTx(ctx, tx)
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	// fmt.Println("respCust FIND NAME")
	respCust, _ := r.channelRepo.FindByApplicationID(ctx, channelParam.ApplicationID, tx)
	respApplication, _ := r.channelRepo.FindByApplicationID(ctx, channelParam.ApplicationID, tx)
	respApplicationName, _ := r.channelRepo.FindByName(ctx, channelParam.ApplicationName, tx)
	// respMerch, _ := r.merchantRepo.FindByEnterprisePrivyID(ctx, respCust.EnterpriseID, tx)

	// if respCust.ApplicationID != "" && respMerch.EnterpriseID == respCust.EnterpriseID {
	if respCust.ApplicationID != "" && respApplication.ApplicationID == channelParam.ApplicationID && respApplicationName.ApplicationName == channelParam.ApplicationName {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Application with ID "+channelParam.ApplicationID+"; Name "+channelParam.ApplicationName+"; Application ID "+channelParam.ApplicationID+" already exist",
			"ApplicationCommandUsecaseGeneral.Create",
			nil,
		)
	}

	// defer func() {
	// 	if p := recover(); p != nil {
	// 		r.channelRepo.RollbackTx(ctx, tx)
	// 		panic(p)
	// 	} else if err != nil {
	// 		log.Println("Rolling back transaction due to error:", err)
	// 		r.channelRepo.RollbackTx(ctx, tx)
	// 	} else {
	// 		err = r.channelRepo.CommitTx(ctx, tx)
	// 		if err != nil {
	// 			log.Println("Error committing transaction:", err)
	// 		}
	// 	}
	// }()

	// fmt.Println("respCust FIND NAME")
	// respCust2, _ := r.channelRepo.FindByApplicationID(ctx, channelParam.ApplicationID, tx)
	// fmt.Println("FINDOUT?")
	// if respCust2.ApplicationID != "" {
	// 	return 0, nil, rapperror.ErrConflict(
	// 		"",
	// 		"Application with ID "+channelParam.ApplicationID+" already exist",
	// 		"ApplicationCommandUsecaseGeneral.Create",
	// 		nil,
	// 	)
	// }

	// find merchant by merchant.EnterpriseID

	var merchant entity.Application
	if len(merchants) > 0 {
		merchant = merchants[0]
	}

	insertApplication := entity.Application{
		EnterpriseID:    merchant.EnterpriseID,
		ApplicationID:   channelParam.ApplicationID,
		ApplicationCode: channelParam.ApplicationCode,
		ApplicationName: channelParam.ApplicationName,
	}

	log.Println("insertApplication", insertApplication)

	channelId, err := r.channelRepo.Create(ctx, insertApplication, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertApplication,
			}).
			Error(err)

		return 0, nil, err
	}

	var merchantId string

	if channelParam.ApplicationID == "PUT" {
		merchantId = channelParam.ApplicationID
	} else {
		merchantId = channelParam.ApplicationID + " - " + merchant.ApplicationName
	}

	print("merchantId", merchantId)

	privyParam := credential.ApplicationParam{
		RecordType:                     "customrecord_customer_hierarchy",
		CustRecordEnterpriseID:         merchant.EnterpriseID,
		CustRecordApplicationID:        channelParam.ApplicationID,
		CustRecordPrivyCodeApplication: channelParam.ApplicationCode,
		CustRecordApplicationName:      channelParam.ApplicationName,
	}

	resp, err := r.channelPrivy.CreateApplication(ctx, privyParam)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		print("resp", resp.Message)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "channelPrivy.CreateApplication",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	// insertApplication. = resp.Data.RecordID
	// insertApplication.CustomerInternalID = merchant.CustomerInternalID
	// insertApplication.ApplicationInternalID = merchant.ApplicationInternalID

	err = r.channelRepo.Update(ctx, channelId, insertApplication, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": insertApplication,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.channelRepo.CommitTx(ctx, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ApplicationCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ApplicationCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return channelId, nil, nil
}

func (r *ApplicationCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Application) (int64, interface{}, error) {
	tx, err := r.channelRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	print("merchantsFIND2", tmNow)

	updatedApplication := entity.Application{
		ApplicationID:   merchant.ApplicationID,
		ApplicationCode: merchant.ApplicationCode,
		ApplicationName: merchant.ApplicationName,
	}

	err = r.channelRepo.Update(ctx, id, updatedApplication, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.channelRepo.CommitTx(ctx, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ApplicationCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ApplicationCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *ApplicationCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.channelRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.channelRepo.Delete(ctx, id, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ApplicationCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.channelRepo.CommitTx(ctx, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ApplicationCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ApplicationCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
