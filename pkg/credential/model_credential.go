package credential

type JWTToken struct {
	SignedJWT string `json:"signedJWT"`
}

func (m JWTToken) Created() int {
	return 0
}

func (m JWTToken) Failed() int {
	return 0
}

type CredentialResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
}

func (m CredentialResponse) Created() int {
	return 0
}

func (m CredentialResponse) Failed() int {
	return 0
}
