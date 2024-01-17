package model

import (
	"github.com/go-playground/validator/v10"
)

type Customer struct {
	CustomerID   string `json:"customerId" validate:"max=255"`
	CustomerType string `json:"customerType" `
	CustomerName string `json:"customerName" validate:"required,max=255"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email" validate:"max=255"`
	PhoneNo      string `json:"phoneNo" validate:"max=255"`
	Address      string `json:"address" validate:"max=1000"`
	IsPerson     bool   `json:"isPerson"`
	// EntityStatus      string  `json:"entityStatus"`
	URL               string  `json:"url"`
	AltPhone          *string `json:"altPhone"`
	Fax               *string `json:"fax"`
	Balance           int     `json:"balanceAmount"`
	Usage             int     `json:"usageAmount"`
	CRMLeadID         string  `json:"crmLeadId" validate:"max=255"`
	EnterprisePrivyID string  `json:"enterprisePrivyId" validate:"required,max=255"`
	Address1          string  `json:"address1"`
	NPWP              string  `json:"npwp" validate:"max=255"`
	State             string  `json:"state" validate:"max=255"`
	City              string  `json:"city"`
	ZipCode           string  `json:"zip" validate:"max=255"`
	CreatedBy         int64   `json:"-"`
}

type Lead struct {
	CustomerID   string `json:"customerId" validate:"max=255"`
	CustomerType string `json:"customerType" `
	CustomerName string `json:"customerName" validate:"required,max=255"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email" validate:"max=255"`
	PhoneNo      string `json:"phoneNo" validate:"max=255"`
	Address      string `json:"address" validate:"max=1000"`
	IsPerson     bool   `json:"isPerson"`
	// EntityStatus      string  `json:"entityStatus"`
	URL               string  `json:"url"`
	AltPhone          *string `json:"altPhone"`
	Fax               *string `json:"fax"`
	Balance           int     `json:"balanceAmount"`
	Usage             int     `json:"usageAmount"`
	CRMLeadID         string  `json:"crmLeadId" validate:"max=255"`
	EnterprisePrivyID string  `json:"enterprisePrivyId" validate:"max=255"`
	Address1          string  `json:"address1"`
	NPWP              string  `json:"npwp" validate:"max=255"`
	State             string  `json:"state" validate:"max=255"`
	City              string  `json:"city"`
	ZipCode           string  `json:"zip" validate:"max=255"`
	CreatedBy         int64   `json:"-"`
}

func (c Customer) Validate() []map[string]interface{} {
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

func (c Lead) ValidateLead() []map[string]interface{} {
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
