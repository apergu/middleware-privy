package credential

type MerchantParam struct {
	RecordType             string `json:"recordtype"`
	CustRecordEnterpriseID string `json:"custrecordenterprise_id"`
	CustRecordMerchantID   string `json:"custrecordmerchant_id"`
	CustRecordMerchantName string `json:"custrecordmerchant_name"`
	CustRecordAddress      string `json:"custrecordaddress"`
	CustRecordEmail        string `json:"custrecordemail"`
	CustRecordPhone        string `json:"custrecordphone"`
	CustRecordState        string `json:"custrecordstate"`
	CustRecordCity         string `json:"custrecordcity"`
	CustRecordZip          string `json:"custrecordzip"`
}
