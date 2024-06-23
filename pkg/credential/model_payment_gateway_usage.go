package credential

type PaymentGatewayParam struct {
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
	CustrecordPrivyChannelNameIntgrasi   string `json:"custrecord_privy_channelname_intgrasi"`
	CustrecordPrivyTypeTransIntegrasi    bool   `json:"custrecord_privy_typetrans_integrasi"`
	CcustrecordPrivyTrxIdIntegrasi       string `json:"custrecord_privy_trxid_integrasi"`
	CustrecordEnterpriseeID              string `json:"custrecordenterprisee_id"`
	CustrecordServiceID                  string `json:"custrecordservice_id"`
	CustrecordUnitPrice                  string `json:"custrecordunit_price"`
}

type PaymentGatewayResponse struct {
	Message string                       `json:"message"`
	Details PaymentGatewayResponseDetail `json:"details"`
}

type PaymentGatewayResponseDetail struct {
	CustomerInternalID int    `json:"customerusage_internalid"`
	CustomerName       string `json:"customername"`
}

type PaymentGatewayFailedResponse struct {
	Error string `json:"error"`
}
