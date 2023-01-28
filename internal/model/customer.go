package model

type Customer struct {
	CustomerID        string  `json:"customerId"`
	CustomerType      string  `json:"customerType"`
	CustomerName      string  `json:"customerName"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Email             string  `json:"email"`
	PhoneNo           string  `json:"phoneNo"`
	Address           string  `json:"address"`
	IsPerson          bool    `json:"isPerson"`
	EntityStatus      string  `json:"entityStatus"`
	URL               string  `json:"url"`
	AltPhone          *string `json:"altPhone"`
	Fax               *string `json:"fax"`
	Balance           int     `json:"balanceAmount"`
	Usage             int     `json:"usageAmount"`
	CRMLeadID         string  `json:"crmLeadId"`
	EnterprisePrivyID string  `json:"enterprisePrivyId"`
	CreatedBy         int64   `json:"-"`
}

func (c Customer) Validate() error {
	return nil
}
