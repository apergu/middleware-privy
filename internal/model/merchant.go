package model

type Merchant struct {
	CustomerID   int64  `json:"customerId"`
	EnterpriseID string `json:"enterpriseId"`
	MerchantID   string `json:"merchantId"`
	MerchantName string `json:"merchantName"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	PhoneNo      string `json:"phoneNo"`
	State        string `json:"state"`
	City         string `json:"city"`
	ZipCode      string `json:"zip"`
	CreatedBy    int64  `json:"-"`
}

func (c Merchant) Validate() error {
	return nil
}
