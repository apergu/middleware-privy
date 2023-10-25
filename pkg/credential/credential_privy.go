package credential

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

type CredentialPrivyProperty struct {
	Host     string
	Username string
	Password string
	Client   *http.Client
}

type CredentialPrivy struct {
	host      string
	username  string
	password  string
	requester *requester
}

func NewCredentialPrivy(prop CredentialPrivyProperty) *CredentialPrivy {
	if prop.Client == nil {
		prop.Client = http.DefaultClient
	}

	r := &requester{
		hc: prop.Client,
	}

	return &CredentialPrivy{
		host:      prop.Host,
		username:  prop.Username,
		password:  prop.Password,
		requester: r,
	}
}

func (c *CredentialPrivy) GenerateJwtTokenWithNode(ctx context.Context) (JWTToken, error) {
	req, _ := http.NewRequest(http.MethodGet, "http://project-privy_nodejs_1:3000", nil)
	req.Header.Set("Content-Type", "application/json")

	jwtToken := JWTToken{}
	err := c.requester.Do(ctx, req, &jwtToken)
	if err != nil {
		return JWTToken{}, err
	}

	return jwtToken, nil
}

func (c *CredentialPrivy) GenerateJwtToken(ctx context.Context) (JWTToken, error) {
	return JWTToken{}, nil
}

func (c *CredentialPrivy) CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error) {
	// get jwt
	isNode := true
	jwtToken := JWTToken{}
	var err error

	if isNode {
		jwtToken, err = c.GenerateJwtTokenWithNode(ctx)
	} else {
		jwtToken, err = c.GenerateJwtToken(ctx)
	}

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "JWTToken{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateCustomer",
			"src":                   "EnvelopeCustomer{}.beforeDo",
			"grant_type":            "client_credentials",
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      jwtToken.SignedJWT,
		}).
		Info(accessTokenURL)

	req, _ := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.username, c.password)

	credential := CredentialResponse{}
	err = c.requester.Do(ctx, req, &credential)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	// post customer
	postCustomerURL := c.host + EndpointPostCustomer

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateCustomer",
			"src":  "EnvelopeCustomer{}.beforeDo",
			"host": postCustomerURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postCustomerURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "9")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeCustomer{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "EnvelopeCustomer{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return CustomerResponse{}, rapperror.ErrNotFound(
			"",
			"Customer is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}

func (c *CredentialPrivy) CreateCustomerUsage(ctx context.Context, param CustomerUsageParam) (CustomerUsageResponse, error) {
	// get jwt
	isNode := true
	jwtToken := JWTToken{}
	var err error

	if isNode {
		jwtToken, err = c.GenerateJwtTokenWithNode(ctx)
	} else {
		jwtToken, err = c.GenerateJwtToken(ctx)
	}

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "JWTToken{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateCustomer",
			"src":                   "EnvelopeCustomer{}.beforeDo",
			"grant_type":            "client_credentials",
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      jwtToken.SignedJWT,
		}).
		Info(accessTokenURL)

	req, _ := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.username, c.password)

	credential := CredentialResponse{}
	err = c.requester.Do(ctx, req, &credential)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	// post customer
	postCustomerURL := c.host + EndpointPostCustomer

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateCustomer",
			"src":  "EnvelopeCustomer{}.beforeDo",
			"host": postCustomerURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postCustomerURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "10")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeCustomerUsage{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "EnvelopeCustomer{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return CustomerUsageResponse{}, rapperror.ErrNotFound(
			"",
			"Customer is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}
