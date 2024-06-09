package usecase

import (
	"context"
	"middleware/internal/model"
	"middleware/pkg/credential"
	"middleware/pkg/privy"
)

type ErpPrivyUsecaseProperty struct {
	ErpPrivyDataPrivy privy.Privy
	ErpPrivyPrivy     credential.Credential
}

type ErpPrivyCommandUsecase interface {
	CheckTopUpStatus(ctx context.Context, topUp model.CheckTopUpStatus) (int64, interface{}, error)
	VoidBalance(ctx context.Context, topUp model.VoidBalance) (int64, interface{}, error)
}
