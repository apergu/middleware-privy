package entity

type TopUpData struct {
	ID                 int64  `json:"id"`
	MerchantID         string `json:"merchantId"`
	TransactionID      string `json:"transactionId"`
	EnterpriseID       string `json:"enterpriseId"`
	EnterpriseName     string `json:"enterpriseName"`
	OriginalServiceID  string `json:"originalServiceId"`
	ServiceID          string `json:"serviceId"`
	ServiceName        string `json:"serviceName"`
	Quantity           int64  `json:"quantity"`
	TransactionDate    int64  `json:"transactionDate"`
	MerchantCode       string `json:"merchantCode"`
	ChannelID          string `json:"channelId"`
	ChannelCode        string `json:"channelCode"`
	CustomerInternalID int64  `json:"customerInternalId"`
	MerchantInternalID int64  `json:"merchantInternalId"`
	ChannelInternalID  int64  `json:"channelInternalId"`
	TransactionType    int8   `json:"transactionType"`
	TopupID            string `json:"topupId"`
	CreatedBy          int64  `json:"-"`
	CreatedAt          int64  `json:"createdAt"`
	UpdatedBy          int64  `json:"-"`
	UpdatedAt          int64  `json:"updatedAt"`
}
