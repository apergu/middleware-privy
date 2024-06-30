package topuppayment

import (
	request "middleware/infrastructure/http/request"
)

type TopUpPaymentGateWayUsecase interface {
	TopUpPaymentGateWay(payload request.PaymentGateway) (interface{}, error)
}
