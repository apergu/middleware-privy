package usecase

import (
	"context"
	"middleware/pkg/erpprivy"
	"time"

	"middleware/internal/helper"
	"middleware/internal/model"

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

func (r *ErpPrivyCommandUsecaseGeneral) TopUpBalance(ctx context.Context, param model.TopUpBalance, xrequestid string) (map[string]interface{}, interface{}, error) {
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
		err := rapperror.ErrUnprocessableEntity(
			"",
			"Start Period Date must be before End Period Date",
			"TopUpBalanceCommandUsecaseGeneral.TopUpBalance",
			nil,
		)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), nil)
		return response, nil, err
	}

	res, err := r.ErpPrivyCred.TopUpBalance(ctx, input, xrequestid)
	if err != nil {
		err := rapperror.ErrUnprocessableEntity(
			"",
			"Start Period Date must be before End Period Date",
			"CheckTopUpStatusCommandUsecaseGeneral.CheckTopUpStatus",
			nil,
		)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), nil)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) CheckTopUpStatus(ctx context.Context, param model.CheckTopUpStatus, xrequestid string) (map[string]interface{}, interface{}, error) {
	input := erpprivy.CheckTopUpStatusParam{
		TopUPID: param.TopUPID,
		Event:   param.Event,
	}

	res, err := r.ErpPrivyCred.CheckTopUpStatus(ctx, input, xrequestid)
	if err != nil {
		err := rapperror.ErrUnprocessableEntity(
			"",
			"Start Period Date must be before End Period Date",
			"CheckTopUpStatusCommandUsecaseGeneral.CheckTopUpStatus",
			nil,
		)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), nil)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) VoidBalance(ctx context.Context, param model.VoidBalance, xrequestid string) (map[string]interface{}, interface{}, error) {
	input := erpprivy.VoidBalanceParam{
		TopUPID: param.TopUPID,
	}

	res, err := r.ErpPrivyCred.VoidBalance(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Adendum(ctx context.Context, param model.Adendum, xrequestid string) (map[string]interface{}, interface{}, error) {
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

	res, err := r.ErpPrivyCred.Adendum(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Reconcile(ctx context.Context, param model.Reconcile, xrequestid string) (map[string]interface{}, interface{}, error) {
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

	res, err := r.ErpPrivyCred.Reconcile(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) TransferBalance(ctx context.Context, param model.TransferBalanceERP, xrequestid string) (map[string]interface{}, interface{}, error) {
	input := erpprivy.TransferBalanceERPParam{
		Origin: struct {
			TopUPID   string "json:\"topup_id\""
			ServiceID string "json:\"service_id\""
		}(param.Origin),
		Destinations: []struct {
			TopUPID      string "json:\"topup_id\""
			EnterpriseId string "json:\"enterprise_id\""
			MerchantId   string "json:\"merchant_id\""
			ChannelId    string "json:\"channel_id\""
			Qty          int    "json:\"qty\""
		}(param.Destinations),
	}

	res, err := r.ErpPrivyCred.TransferBalanceERP(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}
