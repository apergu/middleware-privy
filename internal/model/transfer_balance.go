package model

import "github.com/go-playground/validator/v10"

type TransferBalance struct {
	CustomerId   string `json:"customerId"`
	TransferDate string `json:"transferDate" validate:"required,max=255"`
	TrxIdFrom    string `json:"trxIdFrom" validate:"required,max=255"`
	TrxIdTo      string `json:"trxIdTo" validate:"required,max=255"`
	MerchantTo   string `json:"merchantTo"`
	ChannelTo    string `json:"channelTo"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	IsTrxCreated bool   `json:"isTrxCreated"`
	Quantity     string `json:"quantity"`
	CreatedBy    int64  `json:"-"`
}

func (c TransferBalance) Validate() []map[string]interface{} {
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
