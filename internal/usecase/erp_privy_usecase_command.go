package usecase

import (
	"context"
	"middleware/pkg/erpprivy"

	"middleware/internal/model"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type ErpPrivyCommandUsecaseGeneral struct {
	ErpPrivyCred erpprivy.ErpPrivy
}

func NewErpPrivyCommandUsecaseGeneral(prop ErpPrivyUsecaseProperty) *ErpPrivyCommandUsecaseGeneral {
	return &ErpPrivyCommandUsecaseGeneral{
		ErpPrivyCred: prop.ErpPrivy,
	}
}

func (r *ErpPrivyCommandUsecaseGeneral) CheckTopUpStatus(ctx context.Context, param model.CheckTopUpStatus) (interface{}, error) {
	input := erpprivy.CheckTopUpStatusParam{
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
	input := erpprivy.VoidBalanceParam{
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
