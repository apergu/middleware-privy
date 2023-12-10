package credential

type SalesOrderParams struct {
	RecordType  string `json:"recordtype"`
	Entity      string `json:"entity"`
	TranDate    string `json:"trandate"`
	OrderStatus string `json:"orderstatus"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Memo        string `json:"memo"`
	CustBody2   string `json:"custbody2"`
}

type SalesOrderLinesParams struct {
	Merchant            string `json:"custcol_privy_merchant"`
	Channel             string `json:"custcol_privy_channel"`
	UnitPriceBeforeDisc string `json:"custcol_privy_unitprice_beforedisc"`
	Item                string `json:"item"`
	TaxCode             string `json:"taxcode"`
}

type SalesOrderResponseData struct {
	RecordID int64 `json:"recordId"`
}

type SalesOrderResponse struct {
	Message string               `json:"message"`
	Data    MerchantResponseData `json:"data"`
}

type SalesOrderFailedResponse struct {
	Error string `json:"error"`
}
