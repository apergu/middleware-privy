package erpprivy

import "context"

type ErpPrivy interface {
	TopUpBalance(ctx context.Context, param TopUpBalanceParam) (TopUpBalanceResponse, error)
	CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (CheckTopUpStatusResponse, error)
	VoidBalance(ctx context.Context, param VoidBalanceParam) (VoidBalanceResponse, error)
	Adendum(ctx context.Context, param AdendumParam) (AdendumResponse, error)
	Reconcile(ctx context.Context, param ReconcileParam) (ReconcileResponse, error)
}

type Privy struct {
	ErpPrivy
}
