package credential

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *CredentialPrivy) CreateMerchant(ctx context.Context, param MerchantParam) (MerchantResponse, error) {
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
				"at":  "CredentialPrivy.CreateMerchant",
				"src": "JWTToken{}",
			}).
			Error(err)

		return MerchantResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateMerchant",
			"src":                   "EnvelopeMerchant{}.beforeDo",
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
				"at":  "CredentialPrivy.CreateMerchant",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return MerchantResponse{}, err
	}

	// post customer
	postMerchantURL := c.host + EndpointMerchant

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateMerchant",
			"src":  "EnvelopeMerchant{}.beforeDo",
			"host": postMerchantURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postMerchantURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "126")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeMerchant{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateMerchant",
				"src": "EnvelopeMerchant{}",
			}).
			Error(err)

		return MerchantResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return MerchantResponse{}, nil
	}
	// print("Cust Resp", &custResp.)

	return custResp.SuccessTransaction[0], nil
}
