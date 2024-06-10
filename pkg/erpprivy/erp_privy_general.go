package erpprivy

type ErpPrivyProperty struct {
	Host           string
	Username       string
	Password       string
	ApplicationKey string
	RequestId      string
}

type CredentialERPPrivy struct {
	host           string
	username       string
	password       string
	applicationkey string
	requestid      string
}

func NewCredentialERPPrivy(prop ErpPrivyProperty) *CredentialERPPrivy {
	return &CredentialERPPrivy{
		host:           prop.Host,
		username:       prop.Username,
		password:       prop.Password,
		applicationkey: prop.ApplicationKey,
		requestid:      prop.RequestId,
	}
}
