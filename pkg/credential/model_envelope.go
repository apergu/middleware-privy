package credential

type Envelope interface {
	Created() int
	Failed() int
}

type MainEnvelope struct {
	TotalTransactionCreated       int `json:"total_transaction_created"`
	TotalTransactionFailedCreated int `json:"total_transaction_failed_created"`
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
