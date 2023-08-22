package model

type Divission struct {
	ChannelID     int64  `json:"channelId"`
	DivissionID   string `json:"divissionId"`
	DivissionName string `json:"divissionName"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	PhoneNo       string `json:"phoneNo"`
	State         string `json:"state"`
	City          string `json:"city"`
	ZipCode       string `json:"zip"`
	CreatedBy     int64  `json:"-"`
}

func (c Divission) Validate() error {
	return nil
}
