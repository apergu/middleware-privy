package credential

import "context"

type Customer interface {
	CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error)
}

type CustomerUsage interface {
	CreateCustomerUsage(ctx context.Context, param CustomerUsageParam) (CustomerUsageResponse, error)
}

type Credential interface {
	Customer
	CustomerUsage
}
