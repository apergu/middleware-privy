package usecase

import (
	"context"
	"middleware/internal/model"
	"middleware/pkg/credential"
)

type UsageUsecases interface {
	Create(ctx context.Context, cust model.UsageModel) (*credential.CustomerUsageResponse, error)
}
