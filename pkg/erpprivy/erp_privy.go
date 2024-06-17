package erpprivy

import "context"

type ErpPrivy interface {
	TopUpBalance(ctx context.Context, param TopUpBalanceParam, xrequestid string) (interface{}, error)
	CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam, xrequestid string) (interface{}, error)
	VoidBalance(ctx context.Context, param VoidBalanceParam, xrequestid string) (interface{}, error)
	Adendum(ctx context.Context, param AdendumParam, xrequestid string) (interface{}, error)
	Reconcile(ctx context.Context, param ReconcileParam, xrequestid string) (interface{}, error)
	TransferBalanceERP(ctx context.Context, param TransferBalanceERPParam, xrequestid string) (interface{}, error)
}

type Privy struct {
	ErpPrivy
}
