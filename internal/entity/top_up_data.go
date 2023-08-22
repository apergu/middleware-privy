package entity

type TopUpData struct {
	ID                int64  `json:"id"`
	MerchantID        int64  `json:"merchantId"`
	TransactionID     string `json:"transactionId"`
	EnterpriseID      string `json:"enterpriseId"`
	EnterpriseName    string `json:"enterpriseName"`
	OriginalServiceID string `json:"originalServiceId"`
	ServiceID         string `json:"serviceId"`
	ServiceName       string `json:"serviceName"`
	Quantity          int64  `json:"quantity"`
	TransactionDate   int64  `json:"transactionDate"`
	CreatedBy         int64  `json:"-"`
	CreatedAt         int64  `json:"createdAt"`
	UpdatedBy         int64  `json:"-"`
	UpdatedAt         int64  `json:"updatedAt"`
}
