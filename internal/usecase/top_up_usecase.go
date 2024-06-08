package usecase

import (
	"context"
	"middleware/internal/model"
	"middleware/pkg/credential"
	"middleware/pkg/privy"
)

type TopUpUsecaseProperty struct {
	TopUpDataPrivy privy.Privy
	TopUpPrivy     credential.Credential
}

type TopUpCommandUsecase interface {
	CheckTopUpStatus(ctx context.Context, topUp model.CheckTopUpStatus) (int64, interface{}, error)
}
