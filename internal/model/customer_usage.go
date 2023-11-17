package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type CustomerUsage struct {
	CustomerID          string    `json:"customerId"`
	CustomerName        string    `json:"customerName"`
	ProductID           string    `json:"productId"`
	ProductName         string    `json:"productName"`
	TransactionAt       time.Time `json:"transactionAt"`
	TransactionDate     string    `json:"transactionDate" validate:"required"`
	Balance             int64     `json:"balance"`
	BalanceAmount       float64   `json:"balanceAmount"`
	Usage               int64     `json:"qty" validate:"required,max=100"`
	UsageAmount         float64   `json:"usageAmount"`
	SalesOrderReference int64     `json:"salesOrderReference"`
	MerchantName        string    `json:"merchantID" validate:"required"`
	EnterpriseID        string    `json:"enterpriseId" validate:"required,max=255"`
	EnterpriseName      string    `json:"enterpriseName" validate:"required,max=255"`
	ChannelName         string    `json:"channelID" validate:"required,max=255"`
	TrxId               string    `json:"transactionID" validate:"required,max=255"`
	ServiceID           string    `json:"serviceId" validate:"required,max=255"`
	UnitPrice           string    `json:"unitPrice"`
	TypeTrans           int64     `json:"typeTrans"`
	CreatedBy           int64     `json:"-"`
}

func (c CustomerUsage) Validate() []map[string]interface{} {
	var validationErrors []map[string]interface{}
	v := validator.New()

	// Create a user with invalid data
	// Validate the user struct
	err := v.Struct(c)
	if err != nil {
		// Validation failed, print the error messages
		for _, err := range err.(validator.ValidationErrors) {
			//fmt.Println(err)
			//return err
			validationErrors = append(validationErrors,
				map[string]interface{}{
					"field":       err.Field(),
					"Description": err.Tag(),
				})
		}
	}

	return validationErrors
}
