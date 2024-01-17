package credential

// CustRecordEnterpriseID string `json:"custrecordenterprise_id"`
// CustRecordDivisionID   string `json:"custrecorddivision_id"`
// CustRecordDivisionName string `json:"custrecorddivision_name"`

type MerchantParam struct {
	RecordType                  string `json:"recordtype"`
	CustRecordCustomerName      int64  `json:"custrecordcustomer_name"`
	CustName                    string `json:"custname"`
	CustRecordEnterpriseID      string `json:"custrecordenterprise_id"`
	CustRecordMerchantID        string `json:"custrecordmerchant_id"`
	CustRecordPrivyCodeMerchant string `json:"custrecordprivy_code_merchant"`
	CustRecordMerchantName      string `json:"custrecordmerchant_name"`
	CustRecordAddress           string `json:"custrecordaddress"`
	CustRecordEmail             string `json:"custrecordemail"`
	CustRecordPhone             string `json:"custrecordphone"`
	CustRecordState             string `json:"custrecordstate"`
	CustRecordCity              string `json:"custrecordcity"`
	CustRecordZip               string `json:"custrecordzip"`
	Method                      string `json:"method"`
}

type MerchantResponseData struct {
	RecordID int64 `json:"recordId"`
}

type MerchantResponse struct {
	Message string               `json:"message"`
	Data    MerchantResponseData `json:"data"`
}

type MerchantFailedResponse struct {
	Error string `json:"error"`
}
