package topuppayment

import (
	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	service "middleware/internal/services/privy"
	"middleware/pkg/credential"
)

type TopUpPaymentGatewayUsecaseGeneral struct {
	inf               *infrastructure.Infrastructure
	TopUpPaymentPrivy service.PrivyToNetsuitService
}

func NewTopUpPaymentGateWayGeneral(TopUpPaymentPrivy service.PrivyToNetsuitService, inf *infrastructure.Infrastructure) *TopUpPaymentGatewayUsecaseGeneral {
	return &TopUpPaymentGatewayUsecaseGeneral{
		TopUpPaymentPrivy: TopUpPaymentPrivy,
		inf:               inf,
	}
}

func (r *TopUpPaymentGatewayUsecaseGeneral) TopUpPaymentGateWay(payload request.PaymentGateway) (*credential.EnvelopePaymentGateway, error) {
	url := r.inf.Config.CredentialPrivy.Host + credential.EndpointPostCustomer
	resp := credential.EnvelopeCustomerUsage{}
	req := request.RequestToNetsuit{
		Request:     payload,
		Response:    resp,
		Script:      "183",
		Url:         url,
		ServiceName: "TOP_UP_PAYMENTGATEWAY",
	}

	_, err := r.TopUpPaymentPrivy.ToNetsuit(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
