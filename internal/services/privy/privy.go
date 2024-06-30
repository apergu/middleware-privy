package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	response "middleware/infrastructure/http/response"
	"net/http"
)

type ToNetsuitService struct {
	inf *infrastructure.Infrastructure
}

func NewToNetsuitService(inf *infrastructure.Infrastructure) *ToNetsuitService {
	return &ToNetsuitService{inf: inf}
}

func (ds *ToNetsuitService) ToNetsuit(req request.RequestToNetsuit) (interface{}, error) {

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
	log.Println("accessToken", accessToken)

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": accessToken.TokenType + " " + accessToken.AccessToken,
		"Content-Type":  "application/json",
	}
	reqq := request.RequestToHttpClient{
		Body:        req.Request,
		Headers:     headers,
		Url:         req.Url,
		Method:      http.MethodPost,
		Script:      req.Script,
		ServiceName: req.ServiceName,
	}
	result, err := ds.inf.HttpClient.MakeAPIRequest(reqq)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(result.RespBody))
	fmt.Println(result.RespData.StatusCode)

	switch result.RespData.StatusCode {
	case http.StatusCreated:
		json.Unmarshal(result.RespBody, &req.Response)
		return req.Response, nil
	case http.StatusOK:
		json.Unmarshal(result.RespBody, &req.Response)
		return req.Response, nil
	default:
		return nil, err
	}
}
