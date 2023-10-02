package privy

import "time"

type TopupCreateParam struct {
	TransactionID   string    `json:"transaction_id"`
	SONumber        string    `json:"so_number"`
	EnterpriseID    string    `json:"enterprise_id"`
	MerchantID      string    `json:"merchant_id"`
	ChannelID       string    `json:"channel_id"`
	ServiceID       string    `json:"service_id"`
	PostID          string    `json:"post_id"`
	Quantity        int64     `json:"qty"`
	StartPeriodDate int64     `json:"start_period_date"`
	EndPeriodDate   int64     `json:"end_period_date"`
	TransactionDate time.Time `json:"transaction_date"`
	Reversal        bool      `json:"reversal"`
	ID              string    `json:"id"`
}

type TopupCreateResponse struct {
	QueueID string `json:"queue_id"`
}
