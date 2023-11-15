package model

import "github.com/go-playground/validator/v10"

type Leads struct {
	CustomerID           string `json:"customerId"`
	IsPerson             string `json:"isperson"`
	CompanyName          string `json:"company_name"`
	DefaultOrderPriority string `json:"defaultorderpriority"`
	SalesRep             string `json:"salesrep"`
	Territory            string `json:"territory"`
	Partner              string `json:"partner"`
	Email                string `json:"email"`
	Phone                string `json:"mobile"`
	Fax                  string `json:"fax"`
	CRMLeadID            string `json:"zd_lead_id"`
	EstimatedBudget      string `json:"estimatedbudget" `
	SalesReadiness       string `json:"salesreadiness" `
	BuyingReason         string `json:"buyingreason"`
	EnterpriseId         string `json:"enterprise_id" validate:"required,max=255"`
	NPWP                 string `json:"npwp"`
	CreatedBy            int64  `json:"-"`
}

func (c Leads) Validate() []map[string]interface{} {
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
