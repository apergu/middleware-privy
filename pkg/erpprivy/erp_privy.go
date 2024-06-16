package erpprivy

import "context"

type ErpPrivy interface {
	TopUpBalance(ctx context.Context, param TopUpBalanceParam) (interface{}, error)
	CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (interface{}, error)
	VoidBalance(ctx context.Context, param VoidBalanceParam) (interface{}, error)
	Adendum(ctx context.Context, param AdendumParam) (interface{}, error)
	Reconcile(ctx context.Context, param ReconcileParam) (interface{}, error)
	TransferBalanceERP(ctx context.Context, param TransferBalanceERPParam) (interface{}, error)
}

type Privy struct {
	ErpPrivy
}
