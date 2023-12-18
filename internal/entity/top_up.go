package entity

type TopUp struct {
	TopupID     int64  `json:"topupId"`
	TopUpUUID   string `-`
	SoNo        string `json:"sono"`
	CustomerId  string `json:"customerid"`
	MerchantId  string `json:"merchantid"`
	ChannelId   string `json:"channelid"`
	StartDate   string `json:"startdate"`
	EndDate     string `json:"enddate"`
	Duration    string `json:"duration"`
	Billing     string `json:"billing"`
	ItemId      string `json:"itemid"`
	QtyBalance  int64  `json:"balance"`
	Rate        string `json:"rate"`
	Prepaid     bool   `json:"prepaid"`
	QuotationId string `json:"quotationid"`
	VoidDate    string `json:"void_date"`
	Amount      string `json:"amount"`
	CreatedBy   int64  `json:"-"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedBy   int64  `json:"-"`
	UpdatedAt   int64  `json:"updatedAt"`
}
