package topuppayment

import (
	"middleware/pkg/credential"

	request "middleware/infrastructure/http/request"
)

type TopUpPaymentGateWayUsecase interface {
	TopUpPaymentGateWay(payload request.PaymentGateway) (*credential.EnvelopePaymentGateway, error)
}
