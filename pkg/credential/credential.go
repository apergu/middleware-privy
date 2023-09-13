package credential

import "context"

type Customer interface {
	CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error)
}

type CustomerUsage interface {
	CreateCustomerUsage(ctx context.Context, param CustomerUsageParam) (CustomerUsageResponse, error)
}

type Merchant interface {
	CreateMerchant(ctx context.Context, param MerchantParam) error
}

type Channel interface {
	CreateChannel(ctx context.Context, param ChannelParam) error
}

type Credential interface {
	Customer
	CustomerUsage
	Merchant
	Channel
}
