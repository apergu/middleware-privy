package privy

import "context"

type TopupData interface {
	CreateTopup(ctx context.Context, param TopupCreateParam) (TopupCreateResponse, error)
}

type Privy interface {
	TopupData
}
