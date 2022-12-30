package model

type Customer struct {
	CustomerID   string `json:"customerId"`
	CustomerType string `json:"customerType"`
	CustomerName string `json:"customerName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	PhoneNo      string `json:"phoneNo"`
	Address      string `json:"address"`
	CreatedBy    int64  `json:"-"`
}

func (c Customer) Validate() error {
	return nil
}
