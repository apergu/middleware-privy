package model

import "github.com/go-playground/validator/v10"

type TopUp struct {
	SoNo        string `json:"sono"`
	CustomerId  string `json:"customerid"`
	MerchantId  string `json:"merchantid"`
	ChannelId   string `json:"channelid"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Duration    string `json:"duration"`
	Billing     string `json:"billing"`
	ItemId      string `json:"itemid"`
	QtyBalance  int64  `json:"balance"`
	Rate        string `json:"rate"`
	Prepaid     string `json:"prepaid"`
	QuotationId string `json:"quotationid"`
	VoidDate    string `json:"void_date"`
	Amount      string `json:"amount"`
	CreatedBy   int64  `json:"-"`
}

func (c TopUp) Validate() []map[string]interface{} {
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
