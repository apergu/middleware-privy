package usecase

import (
	"context"
	"middleware/pkg/erpprivy"
	"time"

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

func (r *ErpPrivyCommandUsecaseGeneral) TopUpBalance(ctx context.Context, param model.TopUpBalance) (interface{}, error) {
	input := erpprivy.TopUpBalanceParam{
		TopUPID:         param.TopUPID,
		EnterpriseId:    param.EnterpriseId,
		MerchantId:      param.MerchantId,
		ChannelId:       param.ChannelId,
		ServiceId:       param.ServiceId,
		PostPaid:        param.PostPaid,
		Qty:             param.Qty,
		UnitPrice:       param.UnitPrice,
		StartPeriodDate: param.StartPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   param.EndPeriodDate.Format(time.RFC3339),
		TransactionDate: param.TransactionDate,
	}

	res, err := r.ErpPrivyCred.TopUpBalance(ctx, input)
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
			"Something went wrong when TopUpBalance"+err.Error(),
			"TopUpBalanceCommandUsecaseGeneral.Create",
			nil,
		)
	}

	return res, nil
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
			"Something went wrong when CheckTopUpStatus"+err.Error(),
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
				"at":    "ErpPrivyCommandUsecaseGeneral.VoidBalance" + err.Error(),
				"src":   "ErpPrivyCred.VoidBalance",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when VoidBalance"+err.Error(),
			"VoidBalanceCommandUsecaseGeneral.VoidBalance",
			nil,
		)
	}

	return res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Adendum(ctx context.Context, param model.Adendum) (interface{}, error) {
	input := erpprivy.AdendumParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: param.StartPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   param.EndPeriodDate.Format(time.RFC3339),
		Price:           param.Price,
	}

	res, err := r.ErpPrivyCred.Adendum(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ErpPrivyCommandUsecaseGeneral.Adendedum",
				"src":   "ErpPrivyCred.Adendedum",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when Adendedum"+err.Error(),
			"AdendedumCommandUsecaseGeneral.Adendedum",
			"konz",
		)
	}

	return res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Reconcile(ctx context.Context, param model.Reconcile) (interface{}, error) {
	input := erpprivy.ReconcileParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: param.StartPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   param.EndPeriodDate.Format(time.RFC3339),
		Price:           param.Price,
		Qty:             param.Qty,
	}

	res, err := r.ErpPrivyCred.Reconcile(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ErpPrivyCommandUsecaseGeneral.Reconcile",
				"src":   "ErpPrivyCred.Reconcile",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when Reconcile"+err.Error(),
			"ReconcileCommandUsecaseGeneral.Reconcile",
			"konz",
		)
	}

	return res, nil
}
