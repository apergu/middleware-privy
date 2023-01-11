package credential

import "context"

type Customer interface {
	CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error)
}

type Credential interface {
	Customer
}
