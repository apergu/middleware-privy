package model

import "time"

type CustomerUsage struct {
	CustomerID    string    `json:"customerId"`
	CustomerName  string    `json:"customerName"`
	ProductID     string    `json:"productId"`
	ProductName   string    `json:"productName"`
	TransactionAt time.Time `json:"transactionAt"`
	Balance       int64     `json:"balance"`
	BalanceAmount float64   `json:"balanceAmount"`
	Usage         int64     `json:"usage"`
	UsageAmount   float64   `json:"usageAmount"`
	Transaction   int64     `json:"transaction"`
	CreatedBy     int64     `json:"-"`
}

func (c CustomerUsage) Validate() error {
	return nil
}
