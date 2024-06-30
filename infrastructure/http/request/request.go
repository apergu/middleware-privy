package http_request_infrastructure

import "time"

type CustomerDetails struct {
	CustomForm                 string     `json:"customform" validate:"required"`
	CustBodyPrivySoCustID      string     `json:"custbody_privy_so_custid" validate:"required"`
	Entity                     string     `json:"entity" validate:"required"`
	StartDate                  string     `json:"startdate" validate:"required"`
	EndDate                    string     `json:"enddate" validate:"required"`
	CustBodyPrivyTermOfPayment string     `json:"custbody_privy_termofpayment" validate:"required"`
	OtherRefNum                string     `json:"otherrefnum" validate:"required"`
	CustBodyPrivyBilling       string     `json:"custbody_privy_billing" validate:"required"`
	CustBodyPrivyIntegrasi     string     `json:"custbody_privy_integrasi" validate:"required"`
	Memo                       string     `json:"memo,omitempty"` // Optional
	CustBodyPrivyBDA           string     `json:"custbody_privy_bda" validate:"required"`
	CustBodyPrivyBDM           string     `json:"custbody_privy_bdm" validate:"required"`
	CustBodyPrivySalesSupport  string     `json:"custbody_privy_salessupport" validate:"required"`
	CustBody10                 string     `json:"custbody10" validate:"required"`
	CustBody9                  string     `json:"custbody9" validate:"required"`
	CustBody7                  string     `json:"custbody7" validate:"required"`
	Lines                      []LineItem `json:"lines" validate:"required,dive,required"`
}

type LineItem struct {
	Item                            string `json:"item" validate:"required"`
	CustColPrivyMerchant            string `json:"custcol_privy_merchant" validate:"required"`
	CustColPrivyChannel             string `json:"custcol_privy_channel" validate:"required"`
	CustColPrivyUnitPriceBeforeDisc string `json:"custcol_privy_unitprice_beforedisc" validate:"required"`
	TaxCode                         string `json:"taxcode" validate:"required"`
	CustColPrivyMainProduct         string `json:"custcol_privy_mainproduct,omitempty"`        // Optional
	CustColPrivySubProduct          string `json:"custcol_privy_subproduct,omitempty"`         // Optional
	Description                     string `json:"description,omitempty"`                      // Optional
	Quantity                        int    `json:"quantity,omitempty"`                         // Optional
	CustColPrivyStartDateLayanan    string `json:"custcol_privy_start_date_layanan,omitempty"` // Optional
	CustColPrivyDateLayanan         string `json:"custcol_privy_date_layanan,omitempty"`       // Optional
	CustColPrivyTrxID               string `json:"custcol_privy_trxid" validate:"required"`
	CustColPaymentGatewayFee        string `json:"custcolprivy_paymentgatewayfee" ` // Optional
	Amount                          string `json:"amount" validate:"required"`
	CustColAmountBeforeDisc         string `json:"custcol_privy_amountbeforediscount" `
	// Optional
}

type RequestToNetsuit struct {
	Request     interface{}
	Response    interface{}
	Url         string
	Script      string
	ServiceName string
}

type RequestToHttpClient struct {
	Body        interface{}
	Url         string
	Method      string
	Params      map[string]string
	Headers     map[string]string
	Script      string
	ServiceName string
}

type PaymentGateway struct {
	RecordType                       string     `json:"recordtype" validate:"required"`
	CustomForm                       string     `json:"customform" validate:"required"`
	CustBodyPrivySoCustID            string     `json:"custbody_privy_so_custid" validate:"required"`
	Entity                           string     `json:"entity" validate:"required"`
	StartDate                        string     `json:"startdate" validate:"required"`
	EndDate                          string     `json:"enddate" validate:"required"`
	CustBodyPrivyTermOfPayment       string     `json:"custbody_privy_termofpayment" validate:"required"`
	OtherRefNum                      string     `json:"otherrefnum" validate:"required"`
	CustBodyPrivyBilling             string     `json:"custbody_privy_billing" validate:"required"`
	CustBodyPrivyIntegrasi           string     `json:"custbody_privy_integrasi" validate:"required"`
	Memo                             string     `json:"memo,omitempty"` // Optional
	CustBodyPrivyBDA                 string     `json:"custbody_privy_bda" validate:"required"`
	CustBodyPrivyBDM                 string     `json:"custbody_privy_bdm" validate:"required"`
	CustBodyPrivySalesSupport        string     `json:"custbody_privy_salessupport" validate:"required"`
	CustBodyPrivySalesSupportManager string     `json:"custbodyprivy_salessupportmanager" validate:"required"`
	CustBody10                       string     `json:"custbody10" validate:"required"`
	CustBody9                        string     `json:"custbody9" validate:"required"`
	CustBody7                        string     `json:"custbody7" validate:"required"`
	Lines                            []LineItem `json:"lines" validate:"required,dive,required"`
}

type ResPaymentGateway struct {
	Success  bool `json:"success"`
	RecordID int  `json:"recordId"`
	Data     struct {
		TopupID         string    `json:"topup_id"`
		EnterpriseID    string    `json:"enterprise_id"`
		MerchantID      string    `json:"merchant_id"`
		ChannelID       string    `json:"channel_id"`
		ServiceID       string    `json:"service_id"`
		PostPaid        bool      `json:"post_paid"`
		Qty             int       `json:"qty"`
		UnitPrice       int       `json:"unit_price"`
		StartPeriodDate time.Time `json:"start_period_date"`
		EndPeriodDate   time.Time `json:"end_period_date"`
		TransactionDate time.Time `json:"transaction_date"`
	} `json:"data"`
	Message string `json:"message"`
}
