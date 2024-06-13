package model

import (
	"time"
)

type TopUpBalance struct {
	TopUPID         string `json:"topup_id" validate:"required"`
	EnterpriseId    string `json:"enterprise_id" validate:"required"`
	MerchantId      string `json:"merchant_id"`
	ChannelId       string `json:"channel_id"`
	ServiceId       string `json:"service_id" validate:"required"`
	PostPaid        bool   `json:"post_paid" validate:"required"`
	Qty             int    `json:"qty" validate:"required"`
	UnitPrice       int    `json:"unit_price"`
	EndPeriodDate   string `json:"end_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	StartPeriodDate string `json:"start_period_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TransactionDate string `json:"transaction_date" validate:"required"`
}

type CheckTopUpStatus struct {
	TopUPID string `json:"topup_id" validate:"required"`
	Event   string `json:"event" validate:"required"`
}

type VoidBalance struct {
	TopUPID string `json:"topup_id" validate:"required"`
}
type Adendum struct {
	TopUPID         string `json:"topup_id" validate:"required"`
	StartPeriodDate string `json:"start_period_date" validate:"required"`
	EndPeriodDate   string `json:"end_period_date" validate:"required"`
	Price           int    `json:"price" validate:"required"`
}

type Reconcile struct {
	TopUPID         string    `json:"topup_id" validate:"required"`
	StartPeriodDate time.Time `json:"start_period_date" validate:"required"`
	EndPeriodDate   time.Time `json:"end_period_date"`
	Price           int       `json:"price"`
	Qty             int       `json:"qty"`
}
