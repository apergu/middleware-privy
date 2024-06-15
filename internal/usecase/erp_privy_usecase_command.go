package usecase

import (
	"context"
	"middleware/pkg/erpprivy"
	"time"

	"middleware/internal/helper"
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
		StartPeriodDate: param.StartPeriodDate,
		EndPeriodDate:   param.EndPeriodDate,
		TransactionDate: param.TransactionDate,
	}

	startPeriodDate, _ := time.Parse(time.RFC3339, param.StartPeriodDate)
	endPeriodDate, _ := time.Parse(time.RFC3339, param.EndPeriodDate)

	if startPeriodDate.After(endPeriodDate) {
		return nil, rapperror.ErrBadRequest(
			"",
			"Start Period Date must be before End Period Date",
			"TopUpBalanceCommandUsecaseGeneral.TopUpBalance",
			nil,
		)
	}

	res, err := r.ErpPrivyCred.TopUpBalance(ctx, input)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "ErpPrivyCommandUsecaseGeneral.Create",
				"src":   "topupCred.TopUpBalance",
				"param": param,
			}).
			Error(err)

		return nil, rapperror.ErrInternalServerError(
			"",
			"Something went wrong when TopUpBalance"+err.Error(),
			"TopUpBalanceCommandUsecaseGeneral.TopUpBalance",
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
			"CheckTopUpStatusCommandUsecaseGeneral.CheckTopUpStatus",
			nil,
		)
	}

	return res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) VoidBalance(ctx context.Context, param model.VoidBalance) (map[string]interface{}, interface{}, error) {
	input := erpprivy.VoidBalanceParam{
		TopUPID: param.TopUPID,
	}

	res, err := r.ErpPrivyCred.VoidBalance(ctx, input)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Adendum(ctx context.Context, param model.Adendum) (map[string]interface{}, interface{}, error) {
	input := erpprivy.AdendumParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: param.StartPeriodDate,
		EndPeriodDate:   param.EndPeriodDate,
		Price:           param.Price,
	}

	startPeriodDate, _ := time.Parse(time.RFC3339, param.StartPeriodDate)
	endPeriodDate, _ := time.Parse(time.RFC3339, param.EndPeriodDate)

	if startPeriodDate.After(endPeriodDate) {
		err := rapperror.ErrUnprocessableEntity(
			"",
			"Start Period Date must be before End Period Date",
			"AdendumCommandUsecaseGeneral.Adendum",
			nil,
		)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), nil)
		return response, nil, err
	}

	res, err := r.ErpPrivyCred.Adendum(ctx, input)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Reconcile(ctx context.Context, param model.Reconcile) (map[string]interface{}, interface{}, error) {
	input := erpprivy.ReconcileParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: param.StartPeriodDate,
		EndPeriodDate:   param.EndPeriodDate,
		Price:           param.Price,
		Qty:             param.Qty,
	}

	startPeriodDate, _ := time.Parse(time.RFC3339, param.StartPeriodDate)
	endPeriodDate, _ := time.Parse(time.RFC3339, param.EndPeriodDate)

	if startPeriodDate.After(endPeriodDate) {
		err := rapperror.ErrUnprocessableEntity(
			"",
			"Start Period Date must be before End Period Date",
			"AdendumCommandUsecaseGeneral.Adendum",
			nil,
		)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), nil)
		return response, nil, err
	}

	res, err := r.ErpPrivyCred.Reconcile(ctx, input)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}
