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

type DivissionCommandUsecaseGeneral struct {
	custRepo repository.DivissionCommandRepository
}

func NewDivissionCommandUsecaseGeneral(prop DivissionUsecaseProperty) *DivissionCommandUsecaseGeneral {
	return &DivissionCommandUsecaseGeneral{
		custRepo: prop.DivissionRepo,
	}
}

func (r *DivissionCommandUsecaseGeneral) Create(ctx context.Context, merchant model.Divission) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	insertDivission := entity.Divission{
		ChannelID:     merchant.ChannelID,
		DivissionID:   merchant.DivissionID,
		DivissionName: merchant.DivissionName,
		Address:       merchant.Address,
		Email:         merchant.Email,
		PhoneNo:       merchant.PhoneNo,
		State:         merchant.State,
		City:          merchant.City,
		ZipCode:       merchant.ZipCode,
		CreatedBy:     merchant.CreatedBy,
		CreatedAt:     tmNow,
		UpdatedBy:     merchant.CreatedBy,
		UpdatedAt:     tmNow,
	}

	custId, err := r.custRepo.Create(ctx, insertDivission, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionCommandUsecaseGeneral.Create",
				"src":   "custRepo.Create",
				"param": insertDivission,
			}).
			Error(err)

		return 0, nil, err
	}

	err = r.custRepo.CommitTx(ctx, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":  "DivissionCommandUsecaseGeneral.Create",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"DivissionCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return custId, nil, nil
}

func (r *DivissionCommandUsecaseGeneral) Update(ctx context.Context, id int64, merchant model.Divission) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	tmNow := time.Now().UnixNano() / 1000000

	updatedDivission := entity.Divission{
		ChannelID:     merchant.ChannelID,
		DivissionID:   merchant.DivissionID,
		DivissionName: merchant.DivissionName,
		Address:       merchant.Address,
		Email:         merchant.Email,
		PhoneNo:       merchant.PhoneNo,
		State:         merchant.State,
		City:          merchant.City,
		ZipCode:       merchant.ZipCode,
		UpdatedBy:     merchant.CreatedBy,
		UpdatedAt:     tmNow,
	}

	err = r.custRepo.Update(ctx, id, updatedDivission, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionCommandUsecaseGeneral.Update",
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
				"at":  "DivissionCommandUsecaseGeneral.Update",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"DivissionCommandUsecaseGeneral.Update",
			nil,
		)
	}

	return id, nil, nil
}

func (r *DivissionCommandUsecaseGeneral) Delete(ctx context.Context, id int64) (int64, interface{}, error) {
	tx, err := r.custRepo.BeginTx(ctx)
	if err != nil {
		return 0, nil, err
	}

	err = r.custRepo.Delete(ctx, id, tx)
	if err != nil {
		r.custRepo.RollbackTx(ctx, tx)

		logrus.
			WithFields(logrus.Fields{
				"at":    "DivissionCommandUsecaseGeneral.Delete",
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
				"at":  "DivissionCommandUsecaseGeneral.Delete",
				"src": "custRepo.CommitTx",
			}).
			Error(err)

		return 0, nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when commit",
			"DivissionCommandUsecaseGeneral.Delete",
			nil,
		)
	}

	return id, nil, nil
}
