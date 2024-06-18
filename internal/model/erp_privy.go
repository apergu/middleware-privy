package model

type TopUpBalance struct {
	TopUPID         string `json:"topup_id" validate:"required,formatTopUpID"`
	EnterpriseId    string `json:"enterprise_id" validate:"required"`
	MerchantId      string `json:"merchant_id"`
	ChannelId       string `json:"channel_id"`
	ServiceId       string `json:"service_id" validate:"required"`
	PostPaid        bool   `json:"post_paid"`
	Qty             int    `json:"qty" validate:"required,min=1,max=2147483647"`
	UnitPrice       int    `json:"unit_price" validate:"max=2147483647"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=02/01/2006"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=02/01/2006"`
	TransactionDate string `json:"transaction_date" validate:"required,datetime=02/01/2006"`
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
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=02/01/2006"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=02/01/2006"`
	Price           int    `json:"price" validate:"required"`
}

type Reconcile struct {
	TopUPID         string `json:"topup_id" validate:"required,formatTopUpID"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=02/01/2006"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=02/01/2006"`
	Price           int    `json:"price"`
	Qty             int    `json:"qty"`
}

type TransferBalanceERP struct {
	Origin struct {
		TopUPID   string `json:"topup_id" validate:"required,formatTopUpID"`
		ServiceID string `json:"service_id" validate:"required,min=1,max=20"`
	} `json:"origin"`
	Destinations []struct {
		TopUPID      string `json:"topup_id" validate:"required,formatTopUpID"`
		EnterpriseId string `json:"enterprise_id" validate:"required,min=1,max=100"`
		MerchantId   string `json:"merchant_id" validate:"omitempty,min=1,max=100"`
		ChannelId    string `json:"channel_id" validate:"omitempty,min=1,max=100"`
		Qty          int    `json:"qty" validate:"required,min=1,max=2147483647"`
	} `json:"destinations" validate:"required,min=1,dive"`
}
