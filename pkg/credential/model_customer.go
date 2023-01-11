package credential

type AddressBook struct {
	Addr1           string `json:"addr1"`
	Addr2           string `json:"addr2"`
	Addr3           string `json:"addr3"`
	Attention       string `json:"attention"`
	Override        bool   `json:"override"`
	State           string `json:"state"`
	City            string `json:"city"`
	Zip             string `json:"zip"`
	DefaultBilling  string `json:"defaultbilling"`
	DefaultShipping string `json:"defaultshipping"`
	IsResidential   string `json:"isresidential"`
}

type CustomerParam struct {
	Recordtype                     string  `json:"recordtype"`
	Customform                     string  `json:"customform"`
	EntityID                       string  `json:"entityid"`
	IsPerson                       string  `json:"isperson"`
	CompanyName                    string  `json:"companyname"`
	EntityStatus                   string  `json:"entitystatus"`
	Comments                       string  `json:"comments"`
	URL                            string  `json:"url"`
	Email                          string  `json:"email"`
	Phone                          string  `json:"phone"`
	AltPhone                       *string `json:"altphone"`
	Fax                            *string `json:"fax"`
	CustEntityPrivyCustomerBalance int     `json:"custentity_privy_customer_balance"`
	CustEntityPrivyCustomerUsage   int     `json:"custentityprivy_customer_usage"`
}

type CustomerResponseDetail struct {
	CustomerInternalID int    `json:"customer_internalid"`
	Customerid         string `json:"customerid"`
}

type CustomerResponse struct {
	Message string                 `json:"message"`
	Details CustomerResponseDetail `json:"details"`
}

type CustomerFailedResponse struct {
	Error string `json:"error"`
}
