package model

type UsageModel struct {
	ServiceID       string `json:"serviceID"`
	Qty             int    `json:"qty"`
	TransactionDate string `json:"transactionDate"`
	TransactionID   string `json:"transactionID"`
}
