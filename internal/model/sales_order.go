package model

import "middleware/internal/entity"

type SalesOrder struct {
	ID          int
	Entity      string `json:"entity"`
	TranDate    string `json:"trandate"`
	OrderStatus string `json:"orderstatus"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Memo        string `json:"memo"`
	CustBody2   string `json:"custbody2"`
	Lines       []SalesOrderLines
}

type SalesOrderLines struct {
	Merchant            string `json:"custcol_privy_merchant"`
	Channel             string `json:"custcol_privy_channel"`
	UnitPriceBeforeDisc string `json:"custcol_privy_unitprice_beforedisc"`
	Item                string `json:"item"`
	TaxCode             string `json:"taxcode"`
}

func (c SalesOrder) Validate() error {
	return nil
}

type SalesOrderResponse struct {
	entity.SalesOrder
	Lines []entity.SalesOrderLines `json:"lines"`
}
