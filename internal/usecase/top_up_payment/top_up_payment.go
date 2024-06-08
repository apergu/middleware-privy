package topuppayment

import (
	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	service "middleware/internal/services/privy"
	"middleware/pkg/credential"
)

type TopUpPaymentUsecaseGeneral struct {
	inf               *infrastructure.Infrastructure
	TopUpPaymentPrivy service.PrivyToNetsuitService
}

func NewTopUpPaymentUsecaseGeneral(TopUpPaymentPrivy service.PrivyToNetsuitService, inf *infrastructure.Infrastructure) *TopUpPaymentUsecaseGeneral {
	return &TopUpPaymentUsecaseGeneral{
		TopUpPaymentPrivy: TopUpPaymentPrivy,
		inf:               inf,
	}
}

func (r *TopUpPaymentUsecaseGeneral) TopUpPayment(payload request.CustomerDetails) (*credential.EnvelopeCustomerUsage, error) {

	// custTopUpPayment := strings.Split(payload.TransactionID, "/")
	custPrivyReq := credential.CustomerUsageParam{
		// RecordType:                      "customrecord_privy_integrasi_top_up_payment",
		// CustrecordPrivyServiceIntegrasi: payload.ServiceID,
		// // CustrecordPrivyQuantityIntegrasi:    payload.,
		// CustrecordPrivyUsageDateIntegrasi:   payload.TransactionDate,
		// CustrecordEnterpriseeID:             custTopUpPayment[0],
		// CustrecordPrivyMerchantNameIntgrasi: custTopUpPayment[1],
		// CustrecordPrivyChannelNameIntgrasi:  custTopUpPayment[2],
	}

	script := "10"
	resp := credential.EnvelopeCustomerUsage{}
	url := r.inf.Config.CredentialPrivy.Host + credential.EndpointPostCustomer
	_, err := r.TopUpPaymentPrivy.ToNetsuit(custPrivyReq, resp, url, script, "TOP_UP_PAYMENT")
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
