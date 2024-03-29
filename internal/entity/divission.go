package entity

type Divission struct {
	ID            int64  `json:"id"`
	ChannelID     int64  `json:"channelID"`
	DivissionID   string `json:"divissionId"`
	DivissionName string `json:"divissionName"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	PhoneNo       string `json:"phoneNo"`
	State         string `json:"state"`
	City          string `json:"city"`
	ZipCode       string `json:"zip"`
	CreatedBy     int64  `json:"-"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedBy     int64  `json:"-"`
	UpdatedAt     int64  `json:"updatedAt"`
}
