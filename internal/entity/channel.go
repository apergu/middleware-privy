package entity

type Channel struct {
	ID                 int64  `json:"id"`
	MerchantID         string `json:"merchantID"`
	ChannelCode        string `json:"channelCode"`
	ChannelID          string `json:"channelId"`
	ChannelName        string `json:"channelName"`
	Address            string `json:"address"`
	Email              string `json:"email"`
	PhoneNo            string `json:"phoneNo"`
	State              string `json:"state"`
	City               string `json:"city"`
	ZipCode            string `json:"zip"`
	CustomerInternalID int64  `json:"customerInternalId"`
	MerchantInternalID int64  `json:"merchantInternalId"`
	ChannelInternalID  int64  `json:"channelInternalId"`
	CreatedBy          int64  `json:"-"`
	CreatedAt          int64  `json:"createdAt"`
	UpdatedBy          int64  `json:"-"`
	UpdatedAt          int64  `json:"updatedAt"`
}
