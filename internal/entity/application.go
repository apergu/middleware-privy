package entity

type Application struct {
	ID              int64  `json:"id"`
	EnterpriseID    string `json:"enterpriseId"`
	ApplicationCode string `json:"ApplicationCode"`
	ApplicationID   string `json:"ApplicationId"`
	ApplicationName string `json:"ApplicationName"`
}

type ApplicationFind struct {
	ID              int64  `json:"id"`
	ApplicationID   string `json:"ApplicationId"`
	ApplicationName string `json:"ApplicationName"`
}
