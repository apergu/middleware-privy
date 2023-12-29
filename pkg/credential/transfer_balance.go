package credential

// CustRecordEnterpriseID string `json:"custrecordenterprise_id"`
// CustRecordDivisionID   string `json:"custrecorddivision_id"`
// CustRecordDivisionName string `json:"custrecorddivision_name"`

type TransferBalanceParam struct {
	RecordType                 string `json:"recordtype"`
	CustRecordCustomer         string `json:"custrecordcustomer"`
	CustRecordTransferDate     string `json:"custrecordtransferdate"`
	CustRecordTrxNoFrom        string `json:"custrecordtrxno_from"`
	CustRecordTrxNoTo          string `json:"custrecordtrxno_to"`
	CustRecordMerchant         string `json:"custrecordmerchant"`
	CustRecordChannel          string `json:"custrecordchannel"`
	CustRecordStartDateLayanan string `json:"custrecordstart_date_layanan"`
	CustRecordEndDateLayanan   string `json:"custrecordend_date_layanan"`
	CustRecordIsTrxIdCreated   bool   `json:"custrecordis_trxidcreated"`
	CustRecordFromQuantity     string `json:"custrecordfrom_quantity"`
}

type TransferBalanceResponseData struct {
	RecordID int64 `json:"recordId"`
}

type TransferBalanceResponse struct {
	Message string                      `json:"message"`
	Data    TransferBalanceResponseData `json:"data"`
}

type TransferBalanceFailedResponse struct {
	Error string `json:"error"`
}
