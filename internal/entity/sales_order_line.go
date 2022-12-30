package entity

type SalesOrderLine struct {
	ID                 int64   `json:"id"`
	SalesOrderHeaderId int64   `json:"salesOrderHeaderId"`
	ProductID          string  `json:"productId"`
	ProductName        string  `json:"productName"`
	Quantity           int64   `json:"quantity"`
	RateItem           float64 `json:"rateItem"`
	TaxRate            float64 `json:"taxRate"`
	Subtotal           float64 `json:"subtotal"`
	Grandtotal         float64 `json:"grandtotal"`
	CreatedBy          int64   `json:"-"`
	CreatedAt          int64   `json:"createdAt"`
	UpdatedBy          int64   `json:"-"`
	UpdatedAt          int64   `json:"updatedAt"`
}
