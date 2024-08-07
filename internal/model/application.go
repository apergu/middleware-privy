package model

import "github.com/go-playground/validator/v10"

type Application struct {
	CustomerID      int64  `json:"customerId"`
	EnterpriseID    string `json:"enterpriseId" validate:"max=255"`
	ApplicationCode string `json:"ApplicationCode" validate:"required,max=255"`
	ApplicationID   string `json:"ApplicationId" validate:"required,max=255"`
	ApplicationName string `json:"ApplicationName" validate:"required,max=255"`
	Address         string `json:"address" validate:"max=1000"`
	Email           string `json:"email" validate:"max=255"`
	PhoneNo         string `json:"phoneNo"`
	State           string `json:"state" validate:"max=255"`
	City            string `json:"city" validate:"max=255"`
	ZipCode         string `json:"zip" validate:"max=255"`
	CreatedBy       int64  `json:"-"`
}

func (c Application) Validate() []map[string]interface{} {
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
