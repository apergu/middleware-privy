package model

import "time"

type TopUpData struct {
	MerchantID        int64     `json:"merchantId"`
	TransactionID     string    `json:"transactionId"`
	EnterpriseID      string    `json:"enterpriseId"`
	EnterpriseName    string    `json:"enterpriseName"`
	OriginalServiceID string    `json:"originalServiceId"`
	ServiceID         string    `json:"serviceId"`
	ServiceName       string    `json:"serviceName"`
	Quantity          int64     `json:"quantity"`
	TransactionDate   time.Time `json:"transactionDate"`
	CreatedBy         int64     `json:"-"`
}

func (c TopUpData) Validate() error {
	return nil
}
