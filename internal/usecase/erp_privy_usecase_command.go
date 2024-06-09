package usecase

import (
	"context"
	"middleware/pkg/credential"

	"middleware/internal/model"
	"middleware/pkg/privy"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ErpPrivyCommandUsecaseGeneral struct {
	ErpPrivy     privy.TopupData
	ErpPrivyCred credential.ErpPrivy
}

func NewErpPrivyCommandUsecaseGeneral(prop ErpPrivyUsecaseProperty) *ErpPrivyCommandUsecaseGeneral {
	return &ErpPrivyCommandUsecaseGeneral{
		ErpPrivy:     prop.ErpPrivyDataPrivy,
		ErpPrivyCred: prop.ErpPrivyPrivy,
	}
}

func (r *ErpPrivyCommandUsecaseGeneral) CheckTopUpStatus(ctx context.Context, param model.CheckTopUpStatus) (interface{}, error) {
	input := credential.CheckTopUpStatusParam{
		TopUPID: param.TopUPID,
		Event:   param.Event,
	}

	res, err := r.ErpPrivyCred.CheckTopUpStatus(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ErpPrivyCommandUsecaseGeneral.Create",
				"src":   "topupCred.CreateTopup",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when CheckTopUpStatus",
			"CheckTopUpStatusCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) VoidBalance(ctx context.Context, param model.VoidBalance) (interface{}, error) {
	input := credential.VoidBalanceParam{
		TopUPID: param.TopUPID,
	}

	res, err := r.ErpPrivyCred.VoidBalance(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ErpPrivyCommandUsecaseGeneral.VoidBalance",
				"src":   "ErpPrivyCred.VoidBalance",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when VoidBalance",
			"VoidBalanceCommandUsecaseGeneral.VoidBalance",
			nil,
		)
	}

	return res, nil
}
