package http_client_infratructure

import (
	request "middleware/infrastructure/http/request"
	response "middleware/infrastructure/http/response"
)

type HttpClientInterface interface {
	MakeAPIRequest(req request.RequestToHttpClient) (*HttpResponse, error)
	GetAccessToken(jwtToken response.JWTToken) (*response.CredentialResponse, error)
	GenerateJwtTokenWithNode() (*response.JWTToken, error)
	GenerateJwtToken() (*response.JWTToken, error)
}
