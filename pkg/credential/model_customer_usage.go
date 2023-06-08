package credential

type CustomerUsageParam struct {
	RecordType string `json:"recordtype"`
	// CustrecordPrivySoTransaction    int    `json:"custrecord_privy_so_transaction"`
	// CustrecordPrivyCustomerName     string `json:"custrecord_privy_customer_name"`
	// CustrecordPrivyIdProduct        int    `json:"custrecord_privy_id_product"`
	// CustrecordPrivyProductName      string `json:"custrecord_privy_product_name"`
	// CustrecordPrivyTransactionUsage string `json:"custrecord_privy_transaction_usage"` // 18/01/2023
	// CustrecordPrivyQuantityUsage    int64  `json:"custrecord_privy_quantity_usage"`
	// CustrecordPrivyAmount           int64  `json:"custrecord_privy_amount"`

	CustrecordPrivyUsageDateIntegrasi    string `json:"custrecord_privy_usagedate_integrasi"`    // CustrecordPrivyTransactionUsage
	CustrecordPrivyCustomerNameIntegrasi string `json:"custrecord_privy_customername_integrasi"` // CustrecordPrivyCustomerName
	CustrecordPrivyMerchantNameIntgrasi  string `json:"custrecord_privy_merchantname_intgrasi"`
	CustrecordPrivyServiceIntegrasi      string `json:"custrecord_privy_service_integrasi"`  // CustrecordPrivyProductName
	CustrecordPrivyQuantityIntegrasi     int64  `json:"custrecord_privy_quantity_integrasi"` // CustrecordPrivyQuantityUsage
	CustrecordPrivyTypeTransIntegrasi    int64  `json:"custrecord_privy_typetrans_integrasi"`
}

type CustomerUsageResponse struct {
	Message string                      `json:"message"`
	Details CustomerUsageResponseDetail `json:"details"`
}

type CustomerUsageResponseDetail struct {
	CustomerInternalID int    `json:"customerusage_internalid"`
	CustomerName       string `json:"customername"`
}

type CustomerUsageFailedResponse struct {
	Error string `json:"error"`
}
