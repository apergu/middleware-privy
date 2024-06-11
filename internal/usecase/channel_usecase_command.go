package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ChannelCommandUsecaseGeneral struct {
	channelRepo  repository.ChannelCommandRepository
	channelPrivy credential.Channel
	merchantRepo repository.MerchantQueryRepository
}

func NewChannelCommandUsecaseGeneral(prop ChannelUsecaseProperty) *ChannelCommandUsecaseGeneral {
	return &ChannelCommandUsecaseGeneral{
		channelRepo:  prop.ChannelRepo,
		channelPrivy: prop.ChannelPrivy,
		merchantRepo: prop.MerchantRepo,
	}
}

func (r *ChannelCommandUsecaseGeneral) Create(ctx context.Context, channelParam model.Channel) (int64, interface{}, error) {
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
	tmNow := time.Now().UnixNano() / 1000000

	merchant_filter := repository.MerchantFilter{
		MerchantID: &channelParam.MerchantID,
	}
	merchants, err := r.merchantRepo.Find(ctx, merchant_filter, 1, 0, nil)

	if len(merchants) == 0 {
		return 0, nil, rapperror.ErrNotFound(
			"",
			"Merchant with ID "+channelParam.MerchantID+" is Not Found",
			"ChannelCommandUsecaseGeneral.Create",
			nil,
		)
	}

	fmt.Println("respCust FIND NAME")
	respCust, _ := r.channelRepo.FindByName(ctx, channelParam.ChannelName, tx)

	if respCust.ChannelName != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Channel with name "+channelParam.ChannelName+" already exist",
			"ChannelCommandUsecaseGeneral.Create",
			nil,
		)
	}

	fmt.Println("respCust FIND NAME")
	respCust2, _ := r.channelRepo.FindByChannelID(ctx, channelParam.ChannelID, tx)
	fmt.Println("FINDOUT?")
	if respCust2.ChannelID != "" {
		return 0, nil, rapperror.ErrConflict(
			"",
			"Channel with ID "+channelParam.ChannelID+" already exist",
			"ChannelCommandUsecaseGeneral.Create",
			nil,
		)
	}

	// find merchant by merchant.EnterpriseID

	var merchant entity.Merchant
	if len(merchants) > 0 {
		merchant = merchants[0]
	}

	insertChannel := entity.Channel{
		EnterpriseID: merchant.EnterpriseID,
		MerchantID:   channelParam.MerchantID,
		ChannelCode:  channelParam.ChannelCode,
		ChannelID:    channelParam.ChannelID,
		ChannelName:  channelParam.ChannelName,
		Address:      channelParam.Address,
		Email:        channelParam.Email,
		PhoneNo:      channelParam.PhoneNo,
		State:        channelParam.State,
		City:         channelParam.City,
		ZipCode:      channelParam.ZipCode,
		CreatedBy:    channelParam.CreatedBy,
		CreatedAt:    tmNow,
		UpdatedBy:    channelParam.CreatedBy,
		UpdatedAt:    tmNow,
	}

	log.Println("insertChannel", insertChannel)

	channelId, err := r.channelRepo.Create(ctx, insertChannel, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertChannel,
			}).
			Error(err)

		return 0, nil, err
	}

	var merchantId string

	if channelParam.MerchantID == "PUT" {
		merchantId = channelParam.MerchantID
	} else {
		merchantId = channelParam.MerchantID + " - " + merchant.MerchantName

	}

	privyParam := credential.ChannelParam{
		RecordType:                 "customrecord_customer_hierarchy",
		CustRecordCustomerName:     strconv.Itoa(int(merchant.CustomerInternalID)),
		CustRecordEnterpriseID:     merchant.EnterpriseID,
		CustRecordChannelID:        channelParam.ChannelID,
		CustRecordMerchantID:       merchantId,
		CustRecordPrivyCodeChannel: channelParam.ChannelCode,
		CustRecordChannelName:      channelParam.ChannelName,
		CustRecordAddress:          channelParam.Address,
		CustRecordEmail:            channelParam.Email,
		CustRecordPhone:            channelParam.PhoneNo,
		CustRecordState:            channelParam.State,
		CustRecordCity:             channelParam.City,
		CustRecordZip:              channelParam.ZipCode,
		Method:                     channelParam.Method,
	}

	resp, err := r.channelPrivy.CreateChannel(ctx, privyParam)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Create",
				"src":   "channelPrivy.CreateChannel",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertChannel.ChannelInternalID = resp.Data.RecordID
	insertChannel.CustomerInternalID = merchant.CustomerInternalID
	insertChannel.MerchantInternalID = merchant.MerchantInternalID

	err = r.channelRepo.Update(ctx, channelId, insertChannel, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": insertChannel,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.channelRepo.CommitTx(ctx, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "ChannelCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ChannelCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return channelId, nil, nil
}

func (r *ChannelCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Channel) (int64, interface{}, error) {
	tx, err := r.channelRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedChannel := entity.Channel{
		MerchantID:  merchant.MerchantID,
		ChannelCode: merchant.ChannelCode,
		ChannelID:   merchant.ChannelID,
		ChannelName: merchant.ChannelName,
		Address:     merchant.Address,
		Email:       merchant.Email,
		PhoneNo:     merchant.PhoneNo,
		State:       merchant.State,
		City:        merchant.City,
		ZipCode:     merchant.ZipCode,
		UpdatedBy:   merchant.CreatedBy,
		UpdatedAt:   tmNow,
	}

	err = r.channelRepo.Update(ctx, id, updatedChannel, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Update",
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
				"at":  "ChannelCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ChannelCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *ChannelCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.channelRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.channelRepo.Delete(ctx, id, tx)
	if err != nil {
		r.channelRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Delete",
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
				"at":  "ChannelCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"ChannelCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
