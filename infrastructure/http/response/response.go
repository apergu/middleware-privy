package http_response_infrastructure

type CredentialResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
}

type JWTToken struct {
	SignedJWT string `json:"signedJWT"`
}
