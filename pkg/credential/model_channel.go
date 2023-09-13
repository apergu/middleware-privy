package credential

type ChannelParam struct {
	RecordType            string `json:"recordtype"`
	CustRecordMerchantID  string `json:"custrecordmerchant_id"`
	CustRecordChannelID   string `json:"custrecordchannel_id"`
	CustRecordChannelName string `json:"custrecordchannel_name"`
	CustRecordAddress     string `json:"custrecordaddress"`
	CustRecordEmail       string `json:"custrecordemail"`
	CustRecordPhone       string `json:"custrecordphone"`
	CustRecordState       string `json:"custrecordstate"`
	CustRecordCity        string `json:"custrecordcity"`
	CustRecordZip         string `json:"custrecordzip"`
}
