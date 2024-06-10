package erpprivy

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
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type VoidBalanceFailedResponse struct {
	Code    string `json:"code"`
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
