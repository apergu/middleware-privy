package usecase

import (
	"context"
	"middleware/internal/model"
	"middleware/pkg/erpprivy"
)

type ErpPrivyUsecaseProperty struct {
	ErpPrivy erpprivy.ErpPrivy
}

type ErpPrivyCommandUsecase interface {
	TopUpBalance(ctx context.Context, topUp model.TopUpBalance, xrequestid string) (int64, interface{}, error)
	CheckTopUpStatus(ctx context.Context, topUp model.CheckTopUpStatus, xrequestid string) (int64, interface{}, error)
	VoidBalance(ctx context.Context, topUp model.VoidBalance, xrequestid string) (int64, interface{}, error)
	Adendum(ctx context.Context, topUp model.Adendum, xrequestid string) (int64, interface{}, error)
	Reconcile(ctx context.Context, topUp model.Reconcile, xrequestid string) (int64, interface{}, error)
}
