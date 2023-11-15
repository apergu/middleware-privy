package credential

type LeadParam struct {
	RecordType           string
	CustomerID           string `json:"custentity_privy_enterprise_id"`
	IsPerson             string `json:"isperson"`
	CompanyName          string `json:"companyname"`
	DefaultOrderPriority string `json:"defaultorderpriority"`
	SalesRep             string `json:"salesrep"`
	Territory            string `json:"territory"`
	Partner              string `json:"partner"`
	Email                string `json:"email"`
	Phone                string `json:"mobile"`
	Fax                  string `json:"fax"`
	CRMLeadID            string `json:"custentity_privy_crm_lead_id"`
	EstimatedBudget      string `json:"estimatedbudget" `
	SalesReadiness       string `json:"salesreadiness" `
	BuyingReason         string `json:"buyingreason"`
	EnterpriseId         string `json:"enterprise_id" validate:"required,max=255"`
	BankAccount          string
	NPWP                 string
}

type LeadResponse struct {
	Message string             `json:"message"`
	Details LeadResponseDetail `json:"details"`
}

type LeadResponseDetail struct {
	CustomerInternalID int64  `json:"customer_internalid"`
	Customerid         string `json:"customerid"`
}
