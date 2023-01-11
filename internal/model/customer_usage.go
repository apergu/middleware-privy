package model

import "time"

type CustomerUsage struct {
	CustomerID    string    `json:"customerId"`
	CustomerName  string    `json:"customerName"`
	ProductID     string    `json:"productId"`
	ProductName   string    `json:"productName"`
	IsPerson      bool      `json:"isPerson"`
	EntityStatus  string    `json:"entityStatus"`
	URL           string    `json:"url"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	AltPhone      *string   `json:"altPhone"`
	Fax           *string   `json:"fax"`
	TransactionAt time.Time `json:"transactionAt"`
	Balance       int64     `json:"balance"`
	BalanceAmount float64   `json:"balanceAmount"`
	Usage         int64     `json:"usage"`
	UsageAmount   float64   `json:"usageAmount"`
	CreatedBy     int64     `json:"-"`
}

func (c CustomerUsage) Validate() error {
	return nil
}
