package usecase

import (
	"context"
	"log"
	"middleware/pkg/credential"
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
	topUpCred    credential.TopUp
}

func NewTopUpDataCommandUsecaseGeneral(prop TopUpDataUsecaseProperty) *TopUpDataCommandUsecaseGeneral {
	return &TopUpDataCommandUsecaseGeneral{
		topupRepo:    prop.TopUpDataRepo,
		customerRepo: prop.CustomerRepo,
		merchantRepo: prop.MerchantRepo,
		channelRepo:  prop.ChannelRepo,
		topupPrivy:   prop.TopUpDataPrivy,
		topUpCred:    prop.TopUpPrivy,
	}
}

func (r *TopUpDataCommandUsecaseGeneral) Create(ctx context.Context, topUpData model.TopUp) (int64, interface{}, error) {
	tx, err := r.topupRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	txId := topUpData.CustomerId + "/" + topUpData.MerchantId + "/" + topUpData.ChannelId + "/" + topUpData.Amount
	splittedTxIDs := strings.Split(txId, "/")
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

	//insertTopUpData := entity.TopUpData{
	//	MerchantID:         topUpData.MerchantID,
	//	TransactionID:      topUpData.TransactionID,
	//	EnterpriseID:       topUpData.EnterpriseID,
	//	EnterpriseName:     topUpData.EnterpriseName,
	//	OriginalServiceID:  topUpData.OriginalServiceID,
	//	ServiceID:          topUpData.ServiceID,
	//	ServiceName:        topUpData.ServiceName,
	//	Quantity:           topUpData.Quantity,
	//	TransactionDate:    topUpData.TransactionDate.UnixNano() / 1000000,
	//	MerchantCode:       topUpData.MerchantCode,
	//	ChannelID:          topUpData.ChannelID,
	//	ChannelCode:        topUpData.ChannelCode,
	//	CustomerInternalID: customer.CustomerInternalID,
	//	MerchantInternalID: merchant.MerchantInternalID,
	//	ChannelInternalID:  channel.ChannelInternalID,
	//	TransactionType:    topUpData.TransactionType,
	//	TopupID:            topupIdUUID,
	//	CreatedBy:          topUpData.CreatedBy,
	//	CreatedAt:          tmNow,
	//	UpdatedBy:          topUpData.CreatedBy,
	//	UpdatedAt:          tmNow,
	//}

	insertTopUp := entity.TopUp{
		TopUpUUID:   topupIdUUID,
		SoNo:        topUpData.SoNo,
		Amount:      topUpData.Amount,
		ChannelId:   channel.ChannelID,
		ItemId:      topUpData.ItemId,
		Duration:    topUpData.Duration,
		Prepaid:     topUpData.Prepaid,
		Billing:     topUpData.Billing,
		MerchantId:  merchant.MerchantID,
		QuotationId: topUpData.QuotationId,
		VoidDate:    topUpData.VoidDate,
		QtyBalance:  topUpData.QtyBalance,
		Rate:        topUpData.Rate,
		CustomerId:  customer.CustomerID,
		StartDate:   topUpData.StartDate,
		EndDate:     topUpData.EndDate,
		CreatedBy:   1,
		CreatedAt:   tmNow,
		UpdatedBy:   1,
		UpdatedAt:   tmNow,
	}

	topupDataId, err := r.topupRepo.Create(ctx, insertTopUp, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertTopUp,
			}).
			Error(err)

		return 0, nil, err
	}

	trDate, err := time.Parse("2006-01-02 15:04:05", topUpData.StartDate)

	param := privy.TopupCreateParam{
		TransactionID:   "213/12321/123231",
		SONumber:        "",
		EnterpriseID:    topUpData.MerchantId,
		MerchantID:      merchant.MerchantID,
		ChannelID:       channel.ChannelID,
		ServiceID:       topUpData.ItemId,
		PostID:          "",
		Quantity:        topUpData.QtyBalance,
		StartPeriodDate: tmNow,
		EndPeriodDate:   tmNow,
		TransactionDate: trDate,
		Reversal:        false,
		ID:              topupIdUUID,
	}
	log.Println("topupPrivy", r.topupPrivy)
	log.Println("tup2", r.topupRepo)
	log.Println("tup3", r.channelRepo)
	log.Println("tup4", r.customerRepo)
	log.Println("tup4", r.topUpCred)
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

	privyParam := credential.TopUpParam{
		RecordType:                   "topup",
		CustRecordPrivyMbSoNo:        topUpData.SoNo,
		CustRecordPrivyMbCustomerId:  topUpData.CustomerId,
		CustRecordPrivyMbMerchantId:  topUpData.MerchantId,
		CustRecordPrivyMbChannelId:   topUpData.ChannelId,
		CustRecordPrivyMbStartDate:   topUpData.StartDate,
		CustRecordPrivyMbEndDate:     topUpData.EndDate,
		CustRecordPrivyMbDuration:    topUpData.Duration,
		CustRecordPrivyMbBilling:     topUpData.Billing,
		CustRecordPrivyMbItemId:      topUpData.ItemId,
		CustRecordPrivyMbQtyBalance:  topUpData.QtyBalance,
		CustRecordPrivyMbRate:        topUpData.Rate,
		CustRecordPrivyMbPrepaid:     topUpData.Prepaid,
		CustRecordPrivyMbQuotationId: topUpData.QuotationId,
		CustRecordPrivyMbVoidDate:    topUpData.VoidDate,
		CustRecordPrivyMbAmount:      topUpData.Amount,
	}

	log.Println("TEST", privyParam)
	log.Println("TEST2", r.topUpCred)
	resp, err := r.topUpCred.CreateTopUp(ctx, privyParam)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "merchantPrivy.CreateMerchant",
				"param": privyParam,
			}).
			Error(err)

		return 0, nil, err
	}

	insertTopUp.TopupID = resp.Data.RecordID
	//insertTopUp.TopUpInternalID = insertTopUp

	err = r.topupRepo.Update(ctx, topupDataId, insertTopUp, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "MerchantCommandUsecaseGeneral.Create",
				"src":   "custRepo.Update",
				"param": insertTopUp,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.topupRepo.CommitTx(ctx, tx)
	if err != nil {
		r.topupRepo.RollbackTx(ctx, tx)

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

	err = r.topupRepo.Update2(ctx, id, updatedTopUpData, tx)
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
