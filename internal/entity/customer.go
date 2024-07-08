package entity

type Customer struct {
	ID                 int64  `json:"id"`
	CustomerID         string `json:"customerId"`
	CustomerType       string `json:"customerType"`
	CustomerName       string `json:"customerName"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	PhoneNo            string `json:"phoneNo"`
	Address            string `json:"address"`
	CRMLeadID          string `json:"crmLeadId"`
	EnterprisePrivyID  string `json:"enterprisePrivyId"`
	NPWP               string `json:"npwp"`
	Address1           string `json:"address1"`
	State              string `json:"state"`
	City               string `json:"city"`
	ZipCode            string `json:"zip"`
	CustomerInternalID int64  `json:"customerInternalId"`
	CreatedBy          int64  `json:"-"`
	CreatedAt          int64  `json:"createdAt"`
	UpdatedBy          int64  `json:"-"`
	UpdatedAt          int64  `json:"updatedAt"`
	EntityStatus       string `json:"entityStatus"`
	// CRMDealID          string `json:"crmDealId"`
}

type Subindustry struct {
	ID              int64  `json:"id"`
	SubindustryName string `json:"subindustry_name"`
}
