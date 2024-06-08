package credential

type TopUpParam struct {
	RecordType                   string `json:"recordtype"`
	CustRecordPrivyMbSoNo        string `json:"custrecord_privy_mb_sono"`
	CustRecordPrivyMbCustomerId  string `json:"custrecord_privy_mb_customerid"`
	CustRecordPrivyMbMerchantId  string `json:"custrecord_privy_mb_merchantid"`
	CustRecordPrivyMbChannelId   string `json:"custrecord_privy_mb_channelid"`
	CustRecordPrivyMbStartDate   string `json:"custrecord_privy_mb_startdate"`
	CustRecordPrivyMbEndDate     string `json:"custrecord_privy_mb_enddate"`
	CustRecordPrivyMbDuration    string `json:"cust_record_privy_mb_duration"`
	CustRecordPrivyMbBilling     string `json:"cust_record_privy_mb_billing"`
	CustRecordPrivyMbItemId      string `json:"cust_record_privy_mb_itemid"`
	CustRecordPrivyMbQtyBalance  int64  `json:"cust_record_privy_mb_qty_balance"`
	CustRecordPrivyMbRate        string `json:"cust_record_privy_mb_rate"`
	CustRecordPrivyMbPrepaid     bool   `json:"cust_record_privy_mb_prepaid"`
	CustRecordPrivyMbQuotationId string `json:"cust_record_privy_mb_quotationid"`
	CustRecordPrivyMbVoidDate    string `json:"cust_record_privy_mb_void_date"`
	CustRecordPrivyMbAmount      string `json:"cust_record_privy_mb_amount"`
}

type TopUpResponseData struct {
	RecordID int64 `json:"recordId"`
}

type TopUpResponse struct {
	Message string               `json:"message"`
	Data    MerchantResponseData `json:"data"`
}

type TopUpFailedResponse struct {
	Error string `json:"error"`
}

type CheckTopUpStatusResponse struct {
	Code int `json:"code"`
	Data struct {
		Entity int    `json:"entity"`
		Status string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

type CheckTopUpStatusFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CheckTopUpStatusBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}

type CheckTopUpStatusParam struct {
	TopUPID string `json:"topup_id"`
	Event   string `json:"event"`
}
