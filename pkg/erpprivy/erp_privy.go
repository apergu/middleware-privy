package erpprivy

import "context"

type ErpPrivy interface {
	CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (CheckTopUpStatusResponse, error)
	VoidBalance(ctx context.Context, param VoidBalanceParam) (VoidBalanceResponse, error)
	Adendum(ctx context.Context, param AdendumParam) (AdendumResponse, error)
}

type Privy struct {
	ErpPrivy
}
