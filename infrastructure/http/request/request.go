package http_request_infrastructure

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
type LineItemReq struct {
	Item                            string `json:"item" validate:"required"`
	CustColPrivyMerchant            string `json:""merchantID"" validate:"required"`
	CustColPrivyChannel             string `json:"channelID" validate:"required"`
	CustColPrivyUnitPriceBeforeDisc string `json:"unitePriceBeforeDiscount" validate:"required"`
	TaxCode                         string `json:"taxCode" validate:"required"`
	CustColPrivyMainProduct         string `json:"mainProductID,omitempty"`    // Optional
	CustColPrivySubProduct          string `json:"subProductID,omitempty"`     // Optional
	Description                     string `json:"description,omitempty"`      // Optional
	Quantity                        int    `json:"quantity,omitempty"`         // Optional
	CustColPrivyStartDateLayanan    string `json:"startDateLayanan,omitempty"` // Optional
	CustColPrivyDateLayanan         string `json:"dateLayanan,omitempty"`      // Optional
	CustColPrivyTrxID               string `json:"transactionID" validate:"required"`
	CustColPaymentGatewayFee        string `json:"paymentGatewayFee" ` // Optional
	Amount                          string `json:"amount" validate:"required"`
	CustColAmountBeforeDisc         string `json:"amountBeforeDiscount" `
	// Optional
}

type LineItemPG struct {
	ServiceID               string `json:"serviceID" validate:"required,alphanumSymbol"`
	MerchantID              string `json:"merchantID" validate:"required,alphanumSymbol"`
	ChannelID               string `json:"channelID" validate:"required,alphanumSymbol"`
	UnitPriceBeforeDiscount int    `json:"unitPriceBeforeDiscount" validate:"required,min=0"`
	ItemDiscount            int    `json:"itemDiscount" validate:"required,min=0"`
	StartDateLayanan        string `json:"startDateLayanan" validate:"required,datetime=02/01/2006"`
	EndDateLayanan          string `json:"endDateLayanan" validate:"required,datetime=02/01/2006"`
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
	EnterpriseID string       `json:"enterPriseID" validate:"required,alphanumSymbol"`
	CustomerID   string       `json:"customerID" validate:"required,alphanumSymbol"`
	RequestDate  string       `json:"requestDate" validate:"required,datetime=02/01/2006"`
	StartDate    string       `json:"startdate" validate:"required,datetime=02/01/2006"`
	Lines        []LineItemPG `json:"lines" validate:"required,dive,required"`
}
type PaymentGatewayReq struct {
	// RecordType                       string     `json:"recordtype" validate:"required"`
	// CustomForm                       string     `json:"customform" validate:"required"`
	CustBodyPrivySoCustID            string        `json:"soCustomerID" validate:"required"`
	Entity                           string        `json:"customerID" validate:"required"`
	StartDate                        string        `json:"startDate" validate:"required,datetime=02/01/2006"`
	EndDate                          string        `json:"endDate" validate:"required,datetime=02/01/2006"`
	CustBodyPrivyTermOfPayment       string        `json:"termOfPayment" validate:"required"`
	OtherRefNum                      string        `json:"otherRefNum" validate:"required"`
	CustBodyPrivyBilling             string        `json:"billType" validate:"required"`
	CustBodyPrivyIntegrasi           string        `json:"integrationType" validate:"required"`
	Memo                             string        `json:"memo,omitempty"` // Optional
	CustBodyPrivyBDA                 string        `json:"bda" validate:"required"`
	CustBodyPrivyBDM                 string        `json:"bdm" validate:"required"`
	CustBodyPrivySalesSupport        string        `json:"salesSupport" validate:"required"`
	CustBodyPrivySalesSupportManager string        `json:"salesSupportManager" validate:"required"`
	CustBody10                       string        `json:"invoiceType" validate:"required"`
	CustBody9                        string        `json:"taxReport" validate:"required"`
	CustBody7                        string        `json:"agreementNpwp" validate:"required"`
	Lines                            []LineItemReq `json:"lines" validate:"required,dive,required"`
}

type ResPaymentGateway struct {
	Code         int    `json:"code"`
	EnterPriseID int    `json:"enterPriseID"`
	CustomerID   int    `json:"customerID"`
	RequestDate  string `json:"requestDate"`
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Data         []struct {
		TransactionID           string `json:"transactionID"`
		ServiceID               string `json:"serviceID"`
		MerchantID              string `json:"merchantID"`
		ChannelID               string `json:"channelID"`
		UnitPriceBeforeDiscount int    `json:"unitPriceBeforeDiscount"`
		ItemDiscount            int    `json:"itemDiscount"`
		StartDateLayanan        string `json:"startDateLayanan"`
		EndDateLayanan          string `json:"endDateLayanan"`
	} `json:"data"`
}
