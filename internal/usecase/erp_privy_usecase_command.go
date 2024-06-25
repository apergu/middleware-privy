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
	startPeriodDate, _ := time.Parse("02/01/2006", param.StartPeriodDate)
	endPeriodDate, _ := time.Parse("02/01/2006", param.EndPeriodDate)
	transDate, _ := time.Parse("02/01/2006", param.TransactionDate)

	now := time.Now()

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

	startPeriodDate = time.Date(startPeriodDate.Year(), startPeriodDate.Month(), startPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))
	endPeriodDate = time.Date(endPeriodDate.Year(), endPeriodDate.Month(), endPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))
	transDate = time.Date(transDate.Year(), transDate.Month(), transDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))

	input := erpprivy.TopUpBalanceParam{
		TopUPID:         param.TopUPID,
		EnterpriseId:    param.EnterpriseId,
		MerchantId:      param.MerchantId,
		ChannelId:       param.ChannelId,
		ServiceId:       param.ServiceId,
		PostPaid:        param.PostPaid,
		Qty:             param.Qty,
		UnitPrice:       param.UnitPrice,
		StartPeriodDate: startPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   endPeriodDate.Format(time.RFC3339),
		TransactionDate: transDate.Format(time.RFC3339),
	}

	res, err := r.ErpPrivyCred.TopUpBalance(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
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
	startPeriodDate, _ := time.Parse("02/01/2006", param.StartPeriodDate)
	endPeriodDate, _ := time.Parse("02/01/2006", param.EndPeriodDate)

	now := time.Now()

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

	startPeriodDate = time.Date(startPeriodDate.Year(), startPeriodDate.Month(), startPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))
	endPeriodDate = time.Date(endPeriodDate.Year(), endPeriodDate.Month(), endPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))

	input := erpprivy.AdendumParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: startPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   endPeriodDate.Format(time.RFC3339),
		Price:           param.Price,
	}

	res, err := r.ErpPrivyCred.Adendum(ctx, input, xrequestid)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), res)
		return response, nil, err
	}

	return map[string]interface{}{}, res, nil
}

func (r *ErpPrivyCommandUsecaseGeneral) Reconcile(ctx context.Context, param model.Reconcile, xrequestid string) (map[string]interface{}, interface{}, error) {
	startPeriodDate, _ := time.Parse("02/01/2006", param.StartPeriodDate)
	endPeriodDate, _ := time.Parse("02/01/2006", param.EndPeriodDate)

	now := time.Now()

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

	startPeriodDate = time.Date(startPeriodDate.Year(), startPeriodDate.Month(), startPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))
	endPeriodDate = time.Date(endPeriodDate.Year(), endPeriodDate.Month(), endPeriodDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("WIB", 7*60*60))

	input := erpprivy.ReconcileParam{
		TopUPID:         param.TopUPID,
		StartPeriodDate: startPeriodDate.Format(time.RFC3339),
		EndPeriodDate:   endPeriodDate.Format(time.RFC3339),
		Price:           param.Price,
		Qty:             param.Qty,
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
