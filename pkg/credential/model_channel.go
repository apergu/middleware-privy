package credential

type ChannelParam struct {
	RecordType                 string `json:"recordtype"`
	CustRecordCustomerName     string `json:"custrecordcustomer_name"`
	CustRecordChannelID        string `json:"custrecordchannel_id"`
	CustRecordPrivyCodeChannel string `json:"custrecordprivy_code_channel"`
	CustRecordChannelName      string `json:"custrecordchannel_name"`
	CustRecordAddress          string `json:"custrecordaddress"`
	CustRecordEmail            string `json:"custrecordemail"`
	CustRecordPhone            string `json:"custrecordphone"`
	CustRecordState            string `json:"custrecordstate"`
	CustRecordCity             string `json:"custrecordcity"`
	CustRecordZip              string `json:"custrecordzip"`
}

type ChannelResponseData struct {
	RecordID int64 `json:"recordId"`
}

type ChannelResponse struct {
	Message string              `json:"message"`
	Data    ChannelResponseData `json:"data"`
}

type ChannelFailedResponse struct {
	Error string `json:"error"`
}
