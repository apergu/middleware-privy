package model

import "gitlab.com/mohamadikbal/project-privy/internal/entity"

type SalesOrderHeader struct {
	OrderNumber  string           `json:"orderNumber"`
	CustomerID   string           `json:"customerId"`
	CustomerName string           `json:"customerName"`
	Lines        []SalesOrderLine `json:"lines"`
	CreatedBy    int64            `json:"-"`
}

func (c SalesOrderHeader) Validate() error {
	for _, line := range c.Lines {
		err := line.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type SalesOrderHeaderResponse struct {
	entity.SalesOrderHeader
	Lines []entity.SalesOrderLine `json:"lines"`
}
