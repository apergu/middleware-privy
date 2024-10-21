package credential

type SalesOrderParams struct {
	RecordType   string                  `json:"recordtype"`
	CustomForm   string                  `json:"customform"`
	EnterpriseID string                  `json:"enterprise"`
	Entity       string                  `json:"custbody_privy_tb_cust"`
	TranDate     string                  `json:"trandate"`
	OrderStatus  string                  `json:"orderstatus"`
	StartDate    string                  `json:"startdate"`
	EndDate      string                  `json:"enddate"`
	Memo         string                  `json:"memo"`
	CustBody2    string                  `json:"custbody2"`
	Lines        []SalesOrderLinesParams `json:"lines"`
}

type SalesOrderLinesParams struct {
	Merchant            string `json:"custcol_privy_merchant"`
	Channel             string `json:"custcol_privy_channel"`
	UnitPriceBeforeDisc int    `json:"custcol_privy_unitprice_beforedisc"`
	Item                string `json:"item"`
	TaxCode             string `json:"taxcode"`
	StartDateLayanan    string `json:"custcol_privy_start_date_layanan"`
	EndDateLayanan      string `json:"custcol_privy_date_layanan"`
	Quantity            int    `json:"quantity"`
}

type SalesOrderResponseData struct {
	RecordID any `json:"recordId"`
}

type SalesOrderResponse struct {
	Message string                 `json:"message"`
	Data    SalesOrderResponseData `json:"data"`
}

type SalesOrderFailedResponse struct {
	Error string `json:"error"`
}
