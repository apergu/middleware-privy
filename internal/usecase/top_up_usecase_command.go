package usecase

import (
	"context"
	"middleware/pkg/credential"

	"middleware/internal/model"
	"middleware/pkg/privy"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type TopUpCommandUsecaseGeneral struct {
	topupPrivy privy.TopupData
	topUpCred  credential.TopUp
}

func NewTopUpCommandUsecaseGeneral(prop TopUpUsecaseProperty) *TopUpCommandUsecaseGeneral {
	return &TopUpCommandUsecaseGeneral{
		topupPrivy: prop.TopUpDataPrivy,
		topUpCred:  prop.TopUpPrivy,
	}
}

func (r *TopUpCommandUsecaseGeneral) CheckTopUpStatus(ctx context.Context, param model.CheckTopUpStatus) (interface{}, error) {
	input := credential.CheckTopUpStatusParam{
		TopUPID: param.TopUPID,
		Event:   param.Event,
	}

	res, err := r.topUpCred.CheckTopUpStatus(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "TopUpCommandUsecaseGeneral.Create",
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
