package entity

type SalesOrder struct {
	Entity      string `json:"entity"`
	TranDate    string `json:"trandate"`
	OrderStatus string `json:"orderstatus"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Memo        string `json:"memo"`
	CustBody2   string `json:"custbody2"`
}

type SalesOrderLines struct {
	Merchant            string `json:"custcol_privy_merchant"`
	Channel             string `json:"custcol_privy_channel"`
	UnitPriceBeforeDisc string `json:"custcol_privy_unitprice_beforedisc"`
	Item                string `json:"item"`
	TaxCode             string `json:"taxcode"`
}
