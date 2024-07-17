package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"middleware/internal/entity"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type CustomerUsageCommandUsecaseGeneral struct {
	custUsageRepo      repository.CustomerUsageCommandRepository
	customerUsagePrivy credential.CustomerUsage
	custRepo           repository.CustomerQueryRepository
	merchantRepo       repository.MerchantQueryRepository
	channelRepo        repository.ChannelQueryRepository
}

func NewCustomerUsageCommandUsecaseGeneral(prop CustomerUsageUsecaseProperty) *CustomerUsageCommandUsecaseGeneral {
	return &CustomerUsageCommandUsecaseGeneral{
		custUsageRepo:      prop.CustomerUsageRepo,
		merchantRepo:       prop.MerchantRepo,
		custRepo:           prop.CustRepo,
		customerUsagePrivy: prop.CustomerPrivy,
		channelRepo:        prop.ChannelRepo,
	}
}

func (r *CustomerUsageCommandUsecaseGeneral) Create(ctx context.Context, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custUsageRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	custUsage := strings.Split(cust.TrxId, "/")

	tmNow := time.Now().UnixNano() / 1000000
	transAt := cust.TransactionAt.UnixNano() / 1000000

	customer_filter := repository.CustomerFilter{
		EnterprisePrivyID: &custUsage[0],
	}
	customers, _ := r.custRepo.Find(ctx, customer_filter, 1, 0, nil)

	if len(customers) == 0 {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Find",
				"param": customer_filter,
			}).
			Error(err)

		return 0, nil, fmt.Errorf("[err_unprocessable_entity] customer with enterprise id %s not found", custUsage[0])
	}

	var customer entity.Customer
	if len(customers) > 0 {
		customer = customers[0]
	}

	merchant_filter := repository.MerchantFilter{
		MerchantID: &custUsage[1],
	}
	merchants, err := r.merchantRepo.Find(ctx, merchant_filter, 1, 0, nil)

	if len(merchants) == 0 {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Find",
				"param": customer_filter,
			}).
			Error(err)

		return 0, nil, fmt.Errorf("[err_unprocessable_entity] merchant with merchant id %s not found", custUsage[1])
	}

	var merchant entity.Merchant
	if len(merchants) > 0 {
		merchant = merchants[0]
	}

	channel_filter := repository.ChannelFilter{
		ChannelID: &custUsage[2],
	}

	channels, err := r.channelRepo.Find(ctx, channel_filter, 1, 0, nil)

	if len(channels) == 0 {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Find",
				"param": customer_filter,
			}).
			Error(err)

		return 0, nil, fmt.Errorf("[err_unprocessable_entity] channel with channel id %s not found", custUsage[2])
	}
	var channel entity.Channel
	if len(channels) > 0 {
		channel = channels[0]
	}

	insertCustomerUsage := entity.CustomerUsage{
		CustomerID:     cust.CustomerID,
		CustomerName:   custUsage[0],
		ProductID:      cust.ProductID,
		ProductName:    cust.ProductName,
		TransactionAt:  transAt,
		Balance:        cust.Balance,
		BalanceAmount:  cust.BalanceAmount,
		Usage:          cust.Usage,
		UsageAmount:    cust.UsageAmount,
		EnterpriseID:   custUsage[0],
		EnterpriseName: cust.EnterpriseName,
		ChannelName:    custUsage[2],
		TrxId:          cust.TrxId,
		ServiceID:      cust.ServiceID,
		UnitPrice:      cust.UnitPrice,
		CreatedBy:      cust.CreatedBy,
		CreatedAt:      tmNow,
		UpdatedBy:      cust.CreatedBy,
		UpdatedAt:      tmNow,
	}

	custId, err := r.custUsageRepo.Create(ctx, insertCustomerUsage, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertCustomerUsage,
			}).
			Error(err)

		return 0, nil, err
	}

	custPrivyUsgParam := credential.CustomerUsageParam{
		RecordType:                           "customrecord_privy_integrasi_usage",
		CustrecordPrivyUsageDateIntegrasi:    cust.TransactionDate,
		CustrecordPrivyCustomerNameIntegrasi: customer.CustomerName,
		CustrecordPrivyServiceIntegrasi:      cust.ServiceID,
		CustrecordPrivyMerchantNameIntgrasi:  merchant.MerchantID + " - " + merchant.MerchantName,
		CustrecordPrivyQuantityIntegrasi:     int64(cust.Usage),
		CustrecordPrivyTypeTransIntegrasi:    false,
		CustrecordPrivyChannelNameIntgrasi:   channel.ChannelID + " - " + channel.ChannelName,
		CcustrecordPrivyTrxIdIntegrasi:       cust.TrxId,
		CustrecordEnterpriseeID:              custUsage[0],
		CustrecordServiceID:                  cust.ServiceID,
		CustrecordUnitPrice:                  cust.UnitPrice,
	}

	_, err = r.customerUsagePrivy.CreateCustomerUsage(ctx, custPrivyUsgParam)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Create",
				"src":   "customerPrivy.CreateCustomerUsage",
				"param": custPrivyUsgParam,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return custId, nil, nil
}

func (r *CustomerUsageCommandUsecaseGeneral) Update(ctx context.Context, id int64, cust model.CustomerUsage) (int64, interface{}, error) {
	tx, err := r.custUsageRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000
	transAt := cust.TransactionAt.UnixNano() / 1000000

	updatedCustomerUsage := entity.CustomerUsage{
		CustomerID:     cust.CustomerID,
		CustomerName:   cust.CustomerName,
		ProductID:      cust.ProductID,
		ProductName:    cust.ProductName,
		TransactionAt:  transAt,
		Balance:        cust.Balance,
		BalanceAmount:  cust.BalanceAmount,
		Usage:          cust.Usage,
		UsageAmount:    cust.UsageAmount,
		EnterpriseID:   cust.EnterpriseID,
		EnterpriseName: cust.EnterpriseName,
		ChannelName:    cust.ChannelName,
		TrxId:          cust.TrxId,
		ServiceID:      cust.ServiceID,
		UnitPrice:      cust.UnitPrice,
		UpdatedBy:      cust.CreatedBy,
		UpdatedAt:      tmNow,
	}

	err = r.custUsageRepo.Update(ctx, id, updatedCustomerUsage, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Update",
				"src":   "custRepo.Update",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *CustomerUsageCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.custUsageRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custUsageRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "CustomerUsageCommandUsecaseGeneral.Delete",
				"src":   "custRepo.Delete",
				"param": id,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custUsageRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custUsageRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "CustomerUsageCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"CustomerUsageCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
