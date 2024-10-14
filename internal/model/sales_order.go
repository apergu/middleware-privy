package model

import "middleware/internal/entity"

type SalesOrder struct {
	ID           int
	EnterpriseID string `json:"enterpriseId"`
	TranDate     string `json:"requestDate"`
	OrderStatus  string `json:"orderstatus"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Lines        []SalesOrderLines
}

type SalesOrderLines struct {
	Merchant            string `json:"merchantId"`
	Channel             string `json:"channelId"`
	UnitPriceBeforeDisc int    `json:"unitPriceBeforeDiscount"`
	Item                string `json:"serviceId"`
	StartDateLayanan    string `json:"startDateLayanan"`
	EndDateLayanan      string `json:"endDateLayanan"`
	Quantity            int    `json:"qty"`
}

func (c SalesOrder) Validate() error {
	return nil
}

type SalesOrderResponse struct {
	entity.SalesOrder
	Lines []entity.SalesOrderLines `json:"lines"`
}
