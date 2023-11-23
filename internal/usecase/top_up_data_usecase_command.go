package usecase

import (
	"context"
	"strings"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/privy"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type TopUpDataCommandUsecaseGeneral struct {
	topupRepo    repository.TopUpDataCommandRepository
	customerRepo repository.CustomerQueryRepository
	merchantRepo repository.MerchantQueryRepository
	channelRepo  repository.ChannelQueryRepository
	topupPrivy   privy.TopupData
}

func NewTopUpDataCommandUsecaseGeneral(prop TopUpDataUsecaseProperty) *TopUpDataCommandUsecaseGeneral {
	return &TopUpDataCommandUsecaseGeneral{
		topupRepo:    prop.TopUpDataRepo,
		customerRepo: prop.CustomerRepo,
		merchantRepo: prop.MerchantRepo,
		channelRepo:  prop.ChannelRepo,
		topupPrivy:   prop.TopUpDataPrivy,
	}
}

func (r *TopUpDataCommandUsecaseGeneral) Create(ctx context.Context, topUpData model.TopUpData) (int64, interface{}, error) {
	tx, err := r.topupRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	splittedTxIDs := strings.Split(topUpData.TransactionID, "/")
	if len(splittedTxIDs) != 4 {
		return 0, nil, rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid transcation ID format",
			"TopUpDataCommandUsecaseGeneral.Create",
			nil,
		)
	}

	var customer entity.Customer
	var merchant entity.Merchant
	var channel entity.Channel

	customers, _ := r.customerRepo.Find(ctx, repository.CustomerFilter{CustomerID: &splittedTxIDs[0]}, 1, 0, tx)
	if len(customers) > 0 {
		customer = customers[0]
	}

	merchants, _ := r.merchantRepo.Find(ctx, repository.MerchantFilter{MerchantID: &splittedTxIDs[1]}, 1, 0, tx)
	if len(merchants) > 0 {
		merchant = merchants[0]
	}

	channels, _ := r.channelRepo.Find(ctx, repository.ChannelFilter{ChannelID: &splittedTxIDs[2]}, 1, 0, tx)
	if len(channels) > 0 {
		channel = channels[0]
	}

	topupIdUUID := uuid.New().String()

	insertTopUpData := entity.TopUpData{
		MerchantID:         topUpData.MerchantID,
		TransactionID:      topUpData.TransactionID,
		EnterpriseID:       topUpData.EnterpriseID,
		EnterpriseName:     topUpData.EnterpriseName,
		OriginalServiceID:  topUpData.OriginalServiceID,
		ServiceID:          topUpData.ServiceID,
		ServiceName:        topUpData.ServiceName,
		Quantity:           topUpData.Quantity,
		TransactionDate:    topUpData.TransactionDate.UnixNano() / 1000000,
		MerchantCode:       topUpData.MerchantCode,
		ChannelID:          topUpData.ChannelID,
		ChannelCode:        topUpData.ChannelCode,
		CustomerInternalID: customer.CustomerInternalID,
		MerchantInternalID: merchant.MerchantInternalID,
		ChannelInternalID:  channel.ChannelInternalID,
		TransactionType:    topUpData.TransactionType,
		TopupID:            topupIdUUID,
		CreatedBy:          topUpData.CreatedBy,
		CreatedAt:          tmNow,
		UpdatedBy:          topUpData.CreatedBy,
		UpdatedAt:          tmNow,
	}

	topupDataId, err := r.topupRepo.Create(ctx, insertTopUpData, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertTopUpData,
			}).
			Error(err)

		return 0, nil, err
	}

	param := privy.TopupCreateParam{
		TransactionID:   topUpData.TransactionID,
		SONumber:        "",
		EnterpriseID:    topUpData.EnterpriseID,
		MerchantID:      topUpData.MerchantID,
		ChannelID:       topUpData.ChannelID,
		ServiceID:       topUpData.ServiceID,
		PostID:          "",
		Quantity:        topUpData.Quantity,
		StartPeriodDate: tmNow,
		EndPeriodDate:   tmNow,
		TransactionDate: topUpData.TransactionDate,
		Reversal:        false,
		ID:              topupIdUUID,
	}
	_, err = r.topupPrivy.CreateTopup(ctx, param)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Create",
				"src":   "topupPrivy.CreateTopup",
				"param": param,
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when CreateTopup",
			"TopUpDataCommandUsecaseGeneral.Create",
			nil,
		)
	}

	err = r.topupRepo.CommitTx(ctx, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TopUpDataCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TopUpDataCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return topupDataId, nil, nil
}

func (r *TopUpDataCommandUsecaseGeneral) Update(ctx context.Context, id int64, topUpData model.TopUpData) (int64, interface{}, error) {
	tx, err := r.topupRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedTopUpData := entity.TopUpData{
		MerchantID:        topUpData.MerchantID,
		TransactionID:     topUpData.TransactionID,
		EnterpriseID:      topUpData.EnterpriseID,
		EnterpriseName:    topUpData.EnterpriseName,
		OriginalServiceID: topUpData.OriginalServiceID,
		ServiceID:         topUpData.ServiceID,
		ServiceName:       topUpData.ServiceName,
		Quantity:          topUpData.Quantity,
		TransactionDate:   topUpData.TransactionDate.UnixNano() / 1000000,
		MerchantCode:      topUpData.MerchantCode,
		ChannelID:         topUpData.ChannelID,
		ChannelCode:       topUpData.ChannelCode,
		UpdatedBy:         topUpData.CreatedBy,
		UpdatedAt:         tmNow,
	}

	err = r.topupRepo.Update(ctx, id, updatedTopUpData, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.topupRepo.CommitTx(ctx, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TopUpDataCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TopUpDataCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *TopUpDataCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.topupRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.topupRepo.Delete(ctx, id, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.topupRepo.CommitTx(ctx, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "TopUpDataCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"TopUpDataCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
