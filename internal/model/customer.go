package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Customer struct {
	CustomerID   string  `json:"customerId" validate:"max=255"`
	CustomerType string  `json:"customerType"`
	CustomerName string  `json:"customerName" validate:"alphanum,max=255"`
	FirstName    string  `json:"firstName" validate:"alphanum,max=255"`
	LastName     string  `json:"lastName"  validate:"alphanum,max=255"`
	Email        string  `json:"email" validate:"email,max=255"`
	PhoneNo      string  `json:"phoneNo" validate:"phone,max=255"`
	Address      string  `json:"address" validate:"alphanum,max=1000"`
	IsPerson     bool    `json:"isPerson"`
	EntityStatus string  `json:"entityStatus" validate:"alphanum,max=2"`
	URL          string  `json:"url"`
	AltPhone     *string `json:"altPhone"`
	Fax          *string `json:"fax"`
	Balance      int     `json:"balanceAmount"`
	Usage        int     `json:"usageAmount"`
	CRMLeadID    string  `json:"crmLeadId" validate:"alphanum,max=255"`
	// CRMDealID         string  `json:"crmDealId" validate:"max=255"`
	EnterprisePrivyID string `json:"enterpriseId" validate:"alphanum,max=255"`
	Address1          string `json:"address1"`
	NPWP              string `json:"npwp" validate:"alphanum,max=255"`
	State             string `json:"state" validate:"alphanum,max=255"`
	City              string `json:"city"validate:"alphanum,max=255"`
	ZipCode           string `json:"zip" validate:"alphanum,max=255"`
	CreatedBy         int64  `json:"-"`
	SubIndustry       string `json:"subIndustry" validate:"alphanum,max=255"`
	RequestFrom       string `json:"requestFrom"`
}

type Lead struct {
	CustomerID   string  `json:"customerId" Validate:"max=255"`
	CustomerType string  `json:"customerType" `
	CustomerName string  `json:"customerName" validate:"required,max=255"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Email        string  `json:"email" validate:"email,max=255"`
	PhoneNo      string  `json:"phoneNo" validate:"phone,max=255"`
	Address      string  `json:"address" validate:"max=1000"`
	IsPerson     bool    `json:"isPerson"`
	EntityStatus string  `json:"entityStatus"`
	URL          string  `json:"url"`
	AltPhone     *string `json:"altPhone"`
	Fax          *string `json:"fax"`
	Balance      int     `json:"balanceAmount"`
	Usage        int     `json:"usageAmount"`
	CRMLeadID    string  `json:"crmLeadId" validate:"max=255"`
	// CRMDealID         string  `json:"crmDealId" validate:"max=255"`
	EnterprisePrivyID string `json:"enterpriseId" validate:"max=255"`
	Address1          string `json:"address1"`
	NPWP              string `json:"npwp" validate:"max=255"`
	State             string `json:"state" validate:"max=255"`
	City              string `json:"city"`
	ZipCode           string `json:"zip" validate:"max=255"`
	CreatedBy         int64  `json:"-"`
	SubIndustry       string `json:"subIndustry" validate:"max=255"`
}

func (c Customer) Validate() []map[string]interface{} {
	var validationErrors []map[string]interface{}
	v := validator.New()
	v.RegisterValidation("alphanum", func(fl validator.FieldLevel) bool {
		fmt.Println("VALIDATE", fl.Field().String())
		if fl.Field().String() == "" {
			return true
		}

		return true
	})

	v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		fmt.Println("VALIDATE", email)
		if strings.Contains(email, "@") && strings.Contains(email, ".") {
			return true
		}

		if email == "" {
			return true
		}

		return false

	})

	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		fmt.Println("VALIDATE", phone)
		if phone == "" {
			return true
		}

		pattern := "^[0-9]+$"

		// Kompilasi regex
		re, _ := regexp.Compile(pattern)

		if re.MatchString(phone) {
			return true
		} else {
			return false
		}

	})

	// return fl.Field().String() == "alphanum"

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
					"field":   err.Field(),
					"message": "is " + err.Tag(),
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
