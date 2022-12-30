package model

type SalesOrderLine struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int64   `json:"quantity"`
	RateItem    float64 `json:"rateItem"`
	TaxRate     float64 `json:"taxRate"`
}

func (c SalesOrderLine) Validate() error {
	return nil
}
