package entity

type Application struct {
	ID                    int64  `json:"id"`
	CustomerID            int64  `json:"customer_id"`
	EnterpriseID          string `json:"enterprise_id"`
	ApplicationID         string `json:"application_id"`
	ApplicationName       string `json:"application_name"`
	Address               string `json:"address"`
	Email                 string `json:"email"`
	PhoneNo               string `json:"phone_no"`
	State                 string `json:"state"`
	City                  string `json:"city"`
	ZipCode               string `json:"zip_code"`
	CreatedBy             int64  `json:"created_by"`
	CreatedAt             int64  `json:"created_at"`
	UpdatedBy             int64  `json:"updated_by"`
	UpdatedAt             int64  `json:"updated_at"`
	ApplicationCode       string `json:"application_code"`
	CustomerInternalID    int64  `json:"customer_internalid"`
	ApplicationInternalID int64  `json:"application_internalid"`
}

type ApplicationFind struct {
	ID              int64  `json:"id"`
	ApplicationID   string `json:"ApplicationId"`
	ApplicationName string `json:"ApplicationName"`
}
