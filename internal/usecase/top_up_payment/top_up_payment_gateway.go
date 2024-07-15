package topuppayment

import (
	"encoding/json"
	"errors"
	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	service "middleware/internal/services/privy"
	"middleware/pkg/credential"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
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

func (r *TopUpPaymentGatewayUsecaseGeneral) TopUpPaymentGateWay(payload request.PaymentGateway) (interface{}, error) {
	url := r.inf.Config.CredentialPrivy.Host + credential.EndpointPostCustomer
	resp := credential.EnvelopeCustomerUsage{}
	req := request.RequestToNetsuit{
		Request:     payload,
		Response:    resp,
		Script:      "183",
		Url:         url,
		ServiceName: "TOP_UP_PAYMENTGATEWAY",
	}

	result, err := r.TopUpPaymentPrivy.ToNetsuit(req)
	if err != nil {
		return nil, err
	}

	var resTopUpmodel request.ResPaymentGateway

	resByte, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resByte, &resTopUpmodel)
	if err != nil {
		return nil, err
	}

	if !resTopUpmodel.Success {
		if resTopUpmodel.Message == "" {
			err = rapperror.ErrInternalServerError(
				rapperror.AppErrorCodeInternalServerError,
				"no have message response erp payment gateway",
				"usecase.TopUpPaymentGateWay",
				nil,
			)

			logrus.WithFields(logrus.Fields{
				"action": "no have message response erp payment gateway",
				"at":     "usecase.TopUpPaymentGateWay",
				"src":    "toNetSuit.TopUpPaymentGateWay",
			}).
				Error(err)

			return nil, errors.New("failed to top up, please try again")
		}

		return nil, errors.New(resTopUpmodel.Message)
	}

	return result, nil
}
