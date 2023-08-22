package entity

type Merchant struct {
	ID           int64  `json:"id"`
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
	CreatedAt    int64  `json:"createdAt"`
	UpdatedBy    int64  `json:"-"`
	UpdatedAt    int64  `json:"updatedAt"`
}
