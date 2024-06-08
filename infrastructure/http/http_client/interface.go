package http_client_infratructure

import response "middleware/infrastructure/http/response"

type HttpClientInterface interface {
	MakeAPIRequest(url string, method string, params map[string]string, body interface{}, headers map[string]string, script string, otherData ...interface{}) (*HttpResponse, error)
	GetAccessToken(jwtToken response.JWTToken) (*response.CredentialResponse, error)
	GenerateJwtTokenWithNode() (*response.JWTToken, error)
	GenerateJwtToken() (*response.JWTToken, error)
}
