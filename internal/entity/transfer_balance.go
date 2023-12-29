package entity

type TransferBalance struct {
	ID           int64  `json:"id"`
	InternalId   int64  `json:"-"`
	CustomerId   string `json:"customerId"`
	TransferDate string `json:"transferDate"`
	TrxIdFrom    string `json:"trxIdFrom"`
	TrxIdTo      string `json:"trxIdTo"`
	MerchantTo   string `json:"merchantTo"`
	ChannelTo    string `json:"channelTo"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	IsTrxCreated bool   `json:"isTrxCreated"`
	Quantity     int64  `json:"quantity"`
	CreatedBy    int64  `json:"-"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedBy    int64  `json:"-"`
	UpdatedAt    int64  `json:"updatedAt"`
}
