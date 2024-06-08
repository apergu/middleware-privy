package topuppayment

import (
	"middleware/pkg/credential"

	request "middleware/infrastructure/http/request"
)

type TopUpPaymentUsecase interface {
	TopUpPayment(payload request.CustomerDetails) (*credential.EnvelopeCustomerUsage, error)
}
