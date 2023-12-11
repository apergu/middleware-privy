package credential

import "context"

type Customer interface {
	CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error)
	CreateLead(ctx context.Context, param CustomerParam) (CustomerResponse, error)
	UpdateLead(ctx context.Context, param CustomerParam) (CustomerResponse, error)
}

type Lead interface {
	CreateLead(ctx context.Context, param CustomerParam) (CustomerResponse, error)
}

type TopUp interface {
	CreateTopUp(ctx context.Context, param TopUpParam) (TopUpResponse, error)
}

type CustomerUsage interface {
	CreateCustomerUsage(ctx context.Context, param CustomerUsageParam) (CustomerUsageResponse, error)
}

type SalesOrder interface {
	CreateSalesOrder(ctx context.Context, param SalesOrderParams) (SalesOrderResponse, error)
}

type Merchant interface {
	CreateMerchant(ctx context.Context, param MerchantParam) (MerchantResponse, error)
}

type Channel interface {
	CreateChannel(ctx context.Context, param ChannelParam) (ChannelResponse, error)
}

type Credential interface {
	Lead
	Customer
	CustomerUsage
	Merchant
	Channel
}
