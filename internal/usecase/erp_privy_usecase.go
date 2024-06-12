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
	TopUpBalance(ctx context.Context, topUp model.TopUpBalance) (int64, interface{}, error)
	CheckTopUpStatus(ctx context.Context, topUp model.CheckTopUpStatus) (int64, interface{}, error)
	VoidBalance(ctx context.Context, topUp model.VoidBalance) (int64, interface{}, error)
	Adendum(ctx context.Context, topUp model.Adendum) (int64, interface{}, error)
}
