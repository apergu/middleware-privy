package credential

type ApplicationParam struct {
	RecordType                     string `json:"recordtype"`
	CustRecordCustomerName         string `json:"custrecordcustomer_name"`
	CustRecordEnterpriseID         string `json:"custrecordenterprise_id"`
	CustRecordApplicationID        string `json:"custrecordapplication_id"`
	CustRecordMerchantID           string `json:"custrecordmerchant_id"`
	CustRecordPrivyCodeApplication string `json:"custrecordprivy_code_application"`
	CustRecordApplicationName      string `json:"custrecordchanel_name"`
	CustRecordAddress              string `json:"custrecordaddress"`
	CustRecordEmail                string `json:"custrecordemail"`
	CustRecordPhone                string `json:"custrecordphone"`
	CustRecordState                string `json:"custrecordstate"`
	CustRecordCity                 string `json:"custrecordcity"`
	CustRecordZip                  string `json:"custrecordzip"`
	Method                         string `json:"method"`
}

type ApplicationResponseData struct {
	RecordID int64 `json:"recordId"`
}

type ApplicationResponse struct {
	Message string                  `json:"message"`
	Data    ApplicationResponseData `json:"data"`
}

type ApplicationFailedResponse struct {
	Error string `json:"error"`
}
