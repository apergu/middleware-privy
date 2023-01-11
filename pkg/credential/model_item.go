package credential

type Item struct {
	Recordtype   string  `json:"recordtype"`
	ItemID       string  `json:"itemid"`
	DisplayName  string  `json:"displayname"`
	UnitsType    string  `json:"unitstype"`
	SaleUnit     string  `json:"saleunit"`
	Department   *string `json:"department"`
	Class        *string `json:"class"`
	Location     *string `json:"location"`
	BaseUnit     int     `json:"baseunit"`
	SubType      string  `json:"subtype"`
	SalesTaxCode int     `json:"salestaxcode"`
}
