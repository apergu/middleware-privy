package middleware

import (
	"encoding/json"
	"middleware/infrastructure"
	response "middleware/infrastructure/http/response"
	"net/http"
)

type ToNetsuitService struct {
	inf *infrastructure.Infrastructure
}

func NewToNetsuitService(inf *infrastructure.Infrastructure) *ToNetsuitService {
	return &ToNetsuitService{inf: inf}
}

func (ds *ToNetsuitService) ToNetsuit(req, responseStruct interface{}, url, script, serviceName string) (interface{}, error) {

	// get jwt
	isNode := true
	jwtToken := &response.JWTToken{}
	var err error

	if isNode {
		jwtToken, err = ds.inf.HttpClient.GenerateJwtTokenWithNode()

	} else {
		jwtToken, err = ds.inf.HttpClient.GenerateJwtToken()
	}
	if err != nil {

		return nil, err
	}
	accessToken, err := ds.inf.HttpClient.GetAccessToken(*jwtToken)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"Authorization": accessToken.TokenType + " " + accessToken.AccessToken}

	result, err := ds.inf.HttpClient.MakeAPIRequest(url, http.MethodPost, nil, req, headers, script, serviceName)

	if err != nil {
		return nil, err
	}

	switch result.RespData.StatusCode {
	case http.StatusCreated:
		json.Unmarshal(result.RespBody, &responseStruct)
		return responseStruct, nil
	case http.StatusOK:
		json.Unmarshal(result.RespBody, &responseStruct)
		return responseStruct, nil
	default:
		return nil, err
	}
}
