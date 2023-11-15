package entity

type Leads struct {
	ID                   int64  `json:"id"`
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
	NPWP                 string `json:"npwp"`
	EstimatedBudget      string `json:"estimatedbudget" `
	SalesReadiness       string `json:"salesreadiness" `
	BuyingReason         string `json:"buyingreason"`
	EnterpriseId         string `json:"enterprise_id"`
	BankAccount          string `-`
}
