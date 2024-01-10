package model

import "github.com/go-playground/validator/v10"

type Channel struct {
	MerchantID   string `json:"merchantId" validate:"max=255"`
	EnterpriseID string `json:"enterpriseId" validate:"required,max=255"`
	ChannelCode  string `json:"channelCode" `
	ChannelID    string `json:"channelId" validate:"required,max=255"`
	ChannelName  string `json:"channelName" validate:"required,max=255"`
	Address      string `json:"address" validate:"max=1000"`
	Email        string `json:"email" validate:"email,max=255"`
	PhoneNo      string `json:"phoneNo"`
	State        string `json:"state" validate:"max=255"`
	City         string `json:"city" validate:"max=255"`
	ZipCode      string `json:"zip" validate:"max=255"`
	CreatedBy    int64  `json:"-"`
	Method       string `json:"method"`
}

func (c Channel) Validate() []map[string]interface{} {
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
