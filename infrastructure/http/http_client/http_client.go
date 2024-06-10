package http_client_infratructure

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	// "middleware/config"
	// "middleware/helper/utils"
	"middleware/infrastructure/http/exceptions"
	request "middleware/infrastructure/http/request"
	response "middleware/infrastructure/http/response"
	"middleware/infrastructure/logger"
	"middleware/infrastructure/logger/logrus"
	"middleware/internal/config"
	"net/http"
	"net/url"
	"time"
)

const (
	MaxIdleConns       int  = 100
	MaxIdleConnections int  = 100
	RequestTimeout     int  = 120
	SSL                bool = true
)

type HttpClient struct {
	cfg    *config.Config
	logger logrus.LogrusInterface
}

type HttpResponse struct {
	RespData *http.Response
	RespBody []byte
}

func NewHttpClient(cfg *config.Config, logger logrus.LogrusInterface) *HttpClient {
	return &HttpClient{
		cfg:    cfg,
		logger: logger,
	}
}

func (hc *HttpClient) CreateHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: SSL},
		MaxIdleConns:        MaxIdleConns,
		MaxIdleConnsPerHost: MaxIdleConnections,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func (hc *HttpClient) MakeAPIRequest(request request.RequestToHttpClient) (*HttpResponse, error) {

	var jsonData []byte
	var err error

	if request.Body != nil {
		jsonData, err = json.Marshal(request.Body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling body: %v", err)
		}
	}
	log := hc.logger

	req, err := http.NewRequest(request.Method, request.Url, bytes.NewBuffer(jsonData))

	if err != nil {
		log.CreateLog(&logger.Log{
			StatusCode: http.StatusInternalServerError,
			Method:     request.Method,
			Request:    request.Body,
			URL:        request.Url,
			Message:    err.Error(),
			Response:   nil,
			Service:    request.ServiceName,
		}, logger.LogError, logger.DefaultLogFileName)
		return nil, exceptions.ErrInternalServerError
	}

	q := req.URL.Query()
	for key, value := range request.Params {
		q.Add(key, value)
	}
	if request.Script != "" {
		q.Add("script", request.Script)
		q.Add("deploy", "1")
	}

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	respData, err := hc.CreateHTTPClient().Do(req)
	if err != nil {
		log.CreateLog(&logger.Log{
			StatusCode: http.StatusInternalServerError,
			Method:     request.Method,
			Request:    request.Body,
			URL:        request.Url,
			Message:    err.Error(),
			Response:   nil,
			Service:    request.ServiceName,
		}, logger.LogError, logger.DefaultLogFileName)
		return nil, exceptions.ErrInternalServerError
	}

	respBody, err := io.ReadAll(respData.Body)
	if err != nil {
		log.CreateLog(&logger.Log{
			StatusCode: http.StatusInternalServerError,
			Method:     http.MethodPost,
			Request:    request.Body,
			URL:        request.Url,
			Message:    err.Error(),
			Response:   nil,
			Service:    request.ServiceName,
		}, logger.LogError, logger.DefaultLogFileName)
		return nil, exceptions.ErrInternalServerError
	}
	io.Copy(ioutil.Discard, respData.Body) // <= NOTE
	defer respData.Body.Close()

	log.CreateLog(&logger.Log{
		StatusCode: respData.StatusCode,
		Method:     http.MethodPost,
		Request:    request.Body,
		URL:        request.Url,
		Message:    "success hit service",
		Response:   string(respBody),
		Service:    request.ServiceName,
	}, logger.LogInfo, logger.DefaultLogFileName)

	return &HttpResponse{
		RespData: respData,
		RespBody: respBody,
	}, nil
}

func (c *HttpClient) GetAccessToken(jwtToken response.JWTToken) (*response.CredentialResponse, error) {
	EndpointGetAccessToken := "/services/rest/auth/oauth2/v1/token"
	accessTokenURL := c.cfg.CredentialPrivy.Host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	req, _ := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.cfg.CredentialPrivy.Username, c.cfg.CredentialPrivy.Username)

	resp, err := c.CreateHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, exceptions.ErrInternalServerError
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := response.CredentialResponse{}
	json.Unmarshal(body, &result)

	return &result, nil
}

func (c *HttpClient) GenerateJwtTokenWithNode() (*response.JWTToken, error) {
	req, _ := http.NewRequest(http.MethodGet, "http://"+os.Getenv("JWT_MID")+":3000", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.CreateHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, exceptions.ErrInternalServerError
	}
	var jwtToken response.JWTToken
	if err := json.NewDecoder(resp.Body).Decode(&jwtToken); err != nil {
		return nil, err
	}

	return &jwtToken, nil
}

func (c *HttpClient) GenerateJwtToken() (*response.JWTToken, error) {
	return &response.JWTToken{}, nil
}
