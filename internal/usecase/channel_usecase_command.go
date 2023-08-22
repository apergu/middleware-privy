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

type ChannelCommandUsecaseGeneral struct {
	custRepo repository.ChannelCommandRepository
}

func NewChannelCommandUsecaseGeneral(prop ChannelUsecaseProperty) *ChannelCommandUsecaseGeneral {
	return &ChannelCommandUsecaseGeneral{
		custRepo: prop.ChannelRepo,
	}
}

func (r *ChannelCommandUsecaseGeneral) Create(ctx context.Context, merchant model.Channel) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertChannel := entity.Channel{
		MerchantID:  merchant.MerchantID,
		ChannelID:   merchant.ChannelID,
		ChannelName: merchant.ChannelName,
		Address:     merchant.Address,
		Email:       merchant.Email,
		PhoneNo:     merchant.PhoneNo,
		State:       merchant.State,
		City:        merchant.City,
		ZipCode:     merchant.ZipCode,
		CreatedBy:   merchant.CreatedBy,
		CreatedAt:   tmNow,
		UpdatedBy:   merchant.CreatedBy,
		UpdatedAt:   tmNow,
	}

	custId, err := r.custRepo.Create(ctx, insertChannel, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertChannel,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

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

	return custId, nil, nil
}

func (r *ChannelCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Channel) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedChannel := entity.Channel{
		MerchantID:  merchant.MerchantID,
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

	err = r.custRepo.Update(ctx, id, updatedChannel, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Update",
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
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "ChannelCommandUsecaseGeneral.Delete",
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
