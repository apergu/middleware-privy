package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/entity"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/rteja-library3/rapperror"
)

type TopUpDataCommandUsecaseGeneral struct {
	custRepo repository.TopUpDataCommandRepository
}

func NewTopUpDataCommandUsecaseGeneral(prop TopUpDataUsecaseProperty) *TopUpDataCommandUsecaseGeneral {
	return &TopUpDataCommandUsecaseGeneral{
		custRepo: prop.TopUpDataRepo,
	}
}

func (r *TopUpDataCommandUsecaseGeneral) Create(ctx context.Context, topUpData model.TopUpData) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertTopUpData := entity.TopUpData{
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
		CreatedBy:         topUpData.CreatedBy,
		CreatedAt:         tmNow,
		UpdatedBy:         topUpData.CreatedBy,
		UpdatedAt:         tmNow,
	}

	custId, err := r.custRepo.Create(ctx, insertTopUpData, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertTopUpData,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

	return custId, nil, nil
}

func (r *TopUpDataCommandUsecaseGeneral) Update(ctx context.Context, id int64, topUpData model.TopUpData) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
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

	err = r.custRepo.Update(ctx, id, updatedTopUpData, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Update",
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
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpDataCommandUsecaseGeneral.Delete",
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
