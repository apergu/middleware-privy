package privy

import "net/http"

type PrivyProperty struct {
	Host     string
	Username string
	Password string
	Client   *http.Client
}

type PrivyGeneral struct {
	host      string
	username  string
	password  string
	requester *requester
}

func NewPrivyGeneral(prop PrivyProperty) *PrivyGeneral {
	if prop.Client == nil {
		prop.Client = http.DefaultClient
	}

	r := &requester{
		hc: prop.Client,
	}

	return &PrivyGeneral{
		host:      prop.Host,
		username:  prop.Username,
		password:  prop.Password,
		requester: r,
	}
}
