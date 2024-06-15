package model

type TopUpBalance struct {
	TopUPID         string `json:"topup_id" validate:"required,formatTopUpID"`
	EnterpriseId    string `json:"enterprise_id" validate:"required"`
	MerchantId      string `json:"merchant_id"`
	ChannelId       string `json:"channel_id"`
	ServiceId       string `json:"service_id" validate:"required"`
	PostPaid        bool   `json:"post_paid" validate:"required"`
	Qty             int    `json:"qty" validate:"required,min=1,max=2147483647"`
	UnitPrice       int    `json:"unit_price" validate:"max=2147483647"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TransactionDate string `json:"transaction_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type CheckTopUpStatus struct {
	TopUPID string `json:"topup_id" validate:"required,formatTopUpID"`
	Event   string `json:"event" validate:"required"`
}

type VoidBalance struct {
	TopUPID string `json:"topup_id" validate:"required,formatTopUpID"`
}
type Adendum struct {
	TopUPID         string `json:"topup_id" validate:"required,formatTopUpID"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Price           int    `json:"price" validate:"required"`
}

type Reconcile struct {
	TopUPID         string `json:"topup_id" validate:"required,formatTopUpID"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Price           int    `json:"price"`
	Qty             int    `json:"qty"`
}
