package entity

type Customer struct {
	ID           int64  `json:"id"`
	CustomerID   string `json:"customerId"`
	CustomerType string `json:"customerType"`
	CustomerName string `json:"customerName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	PhoneNo      string `json:"phoneNo"`
	Address      string `json:"address"`
	CreatedBy    int64  `json:"-"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedBy    int64  `json:"-"`
	UpdatedAt    int64  `json:"updatedAt"`
}
