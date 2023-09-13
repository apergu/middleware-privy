package entity

type CustomerUsage struct {
	ID            int64   `json:"id"`
	CustomerID    string  `json:"customerId"`
	CustomerName  string  `json:"customerName"`
	ProductID     string  `json:"productId"`
	ProductName   string  `json:"productName"`
	TransactionAt int64   `json:"transactionAt"`
	Balance       int64   `json:"balance"`
	BalanceAmount float64 `json:"balanceAmount"`
	Usage         int64   `json:"usage"`
	UsageAmount   float64 `json:"usageAmount"`
	CreatedBy     int64   `json:"-"`
	CreatedAt     int64   `json:"createdAt"`
	UpdatedBy     int64   `json:"-"`
	UpdatedAt     int64   `json:"updatedAt"`

	EnterpriseID   string `json:"enterpriseId"`
	EnterpriseName string `json:"enterpriseName"`
	ChannelName    string `json:"channelName"`
	TrxId          string `json:"trxId"`
	ServiceID      string `json:"serviceId"`
	UnitPrice      string `json:"unitPrice"`
}
