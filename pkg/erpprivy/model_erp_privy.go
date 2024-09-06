package erpprivy

type TopUpBalanceResponse struct {
	Code int `json:"code"`
	Data struct {
		Entity int    `json:"entity"`
		Status string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

type TopUpBalanceFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TopUpBalanceBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}

type TopUpBalanceParam struct {
	TopUPID         string `json:"topup_id"`
	EnterpriseId    string `json:"enterprise_id"`
	ApplicationId   string `json:"application_id"`
	ChannelId       string `json:"channel_id"`
	ServiceId       string `json:"service_id"`
	PostPaid        bool   `json:"post_paid"`
	Qty             int    `json:"qty"`
	UnitPrice       int    `json:"unit_price"`
	StartPeriodDate string `json:"start_period_date"`
	EndPeriodDate   string `json:"end_period_date"`
	TransactionDate string `json:"transaction_date"`
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

type VoidBalanceParam struct {
	TopUPID string `json:"topup_id"`
}

type VoidBalanceResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type VoidBalanceFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type VoidBalanceBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}

type AdendumParam struct {
	TopUPID         string `json:"topup_id"`
	StartPeriodDate string `json:"start_period_date"`
	EndPeriodDate   string `json:"end_period_date"`
	Price           int    `json:"price"`
}

type AdendumResponse struct {
	Code    int    `json:"code"`
	Data    int    `json:"data"`
	Message string `json:"message"`
}

type AdendumFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AdendumBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}

type ReconcileParam struct {
	TopUPID         string `json:"topup_id"`
	StartPeriodDate string `json:"start_period_date"`
	EndPeriodDate   string `json:"end_period_date"`
	Price           int    `json:"price"`
	Qty             int    `json:"qty"`
}

type ReconcileResponse struct {
	Code    int    `json:"code"`
	Data    int    `json:"data"`
	Message string `json:"message"`
}

type ReconcileFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ReconcileBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}

type TransferBalanceERPParam struct {
	Origin struct {
		TopUPID   string `json:"topup_id"`
		ServiceID string `json:"service_id"`
	} `json:"origin"`
	Destinations []struct {
		TopUPID       string `json:"topup_id"`
		EnterpriseId  string `json:"enterprise_id"`
		ApplicationID string `json:"application_id"`
		ChannelId     string `json:"channel_id"`
		Qty           int    `json:"qty"`
	} `json:"destinations"`
}

type TransferBalanceERPResponse struct {
	Code    int    `json:"code"`
	Data    int    `json:"data"`
	Message string `json:"message"`
}

type TransferBalanceERPAlreadySuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TransferBalanceERPFailedResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TransferBalanceERPBadRequestResponse struct {
	Code   int `json:"code"`
	Errors []struct {
		Field       string `json:"field"`
		Description string `json:"description"`
	} `json:"errors"`
	Message string `json:"message"`
}
