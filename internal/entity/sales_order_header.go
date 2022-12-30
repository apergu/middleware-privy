package entity

type SalesOrderHeader struct {
	ID           int64   `json:"id"`
	OrderNumber  string  `json:"orderNumber"`
	CustomerID   string  `json:"customerId"`
	CustomerName string  `json:"customerName"`
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Grandtotal   float64 `json:"grandtotal"`
	CreatedBy    int64   `json:"-"`
	CreatedAt    int64   `json:"createdAt"`
	UpdatedBy    int64   `json:"-"`
	UpdatedAt    int64   `json:"updatedAt"`
}
