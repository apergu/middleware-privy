package credential

type Envelope interface {
	Created() int
	Failed() int
}

type MainEnvelope struct {
	TotalTransactionCreated       int                    `json:"total_transaction_created"`
	TotalTransactionFailedCreated int                    `json:"total_transaction_failed_created"`
	Error                         map[string]interface{} `json:"error"`
}

func (m MainEnvelope) Created() int {
	return m.TotalTransactionCreated
}

func (m MainEnvelope) Failed() int {
	return m.TotalTransactionFailedCreated
}

type EnvelopeCustomer struct {
	MainEnvelope
	SuccessTransaction []CustomerResponse       `json:"success_transaction"`
	FailedTransaction  []CustomerFailedResponse `json:"failed_transaction"`
}

type EnvelopeCustomerUsage struct {
	MainEnvelope
	SuccessTransaction []CustomerUsageResponse       `json:"success_transaction"`
	FailedTransaction  []CustomerUsageFailedResponse `json:"failed_transaction"`
}

type EnvelopeMerchant struct {
	MainEnvelope
	SuccessTransaction []MerchantResponse       `json:"success_transaction"`
	FailedTransaction  []MerchantFailedResponse `json:"failed_transaction"`
}

type EnvelopeChannel struct {
	MainEnvelope
	SuccessTransaction []ChannelResponse       `json:"success_transaction"`
	FailedTransaction  []ChannelFailedResponse `json:"failed_transaction"`
}
