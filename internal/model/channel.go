package model

type Channel struct {
	MerchantID  int64  `json:"merchantId"`
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNo     string `json:"phoneNo"`
	State       string `json:"state"`
	City        string `json:"city"`
	ZipCode     string `json:"zip"`
	CreatedBy   int64  `json:"-"`
}

func (c Channel) Validate() error {
	return nil
}
