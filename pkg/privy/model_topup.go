package privy

import "time"

type TopupCreateParam struct {
	// TransactionID   string    `json:"transaction_id"`
	// SONumber        string    `json:"so_number"`
	// EnterpriseID    string    `json:"enterprise_id"`
	// MerchantID      string    `json:"merchant_id"`
	// ChannelID       string    `json:"channel_id"`
	// ServiceID       string    `json:"service_id"`
	// PostID          string    `json:"post_id"`
	// Quantity        int64     `json:"qty"`
	// StartPeriodDate int64     `json:"start_period_date"`
	// EndPeriodDate   int64     `json:"end_period_date"`
	// TransactionDate time.Time `json:"transaction_date"`
	// Reversal        bool      `json:"reversal"`
	// ID              string    `json:"id"`
	TopUpUUID       string    `json:"topup_uuid"`
	TopUpID         string    `json:"topup_id"`
	EnterpriseID    string    `json:"enterprise_id"`
	MerchantID      string    `json:"merchant_id"`
	ChannelID       string    `json:"channel_id"`
	ServiceID       string    `json:"service_id"`
	TransactionDate time.Time `json:"transaction_date"`
	StartPeriodDate time.Time `json:"start_period_date"`
	EndPeriodDate   time.Time `json:"end_period_date"`
	PostPaid        bool      `json:"post_paid"`
	Reversal        bool      `json:"reversal"`
	Qty             int64     `json:"qty"`
}

type TopupCreateResponse struct {
	QueueID string `json:"queue_id"`
}
