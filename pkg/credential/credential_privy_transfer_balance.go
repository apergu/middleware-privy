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

func (c *CredentialPrivy) CreateTransferBalance(ctx context.Context, param TransferBalanceParam) (TransferBalanceResponse, error) {
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
				"at":  "CredentialPrivy.CreateTransferBalance",
				"src": "JWTToken{}",
			}).
			Error(err)

		return TransferBalanceResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateTransferBalance",
			"src":                   "EnvelopeTransferBalance{}.beforeDo",
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
				"at":  "CredentialPrivy.CreateTransferBalance",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return TransferBalanceResponse{}, err
	}

	// post customer
	postTransferBalanceURL := c.host + EndpointChannel

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateTransferBalance",
			"src":  "EnvelopeTransferBalance{}.beforeDo",
			"host": postTransferBalanceURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postTransferBalanceURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "85")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeTransferBalance{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateTransferBalance",
				"src": "EnvelopeTransferBalance{}",
			}).
			Error(err)

		return TransferBalanceResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return TransferBalanceResponse{}, rapperror.ErrNotFound(
			"",
			"TransferBalance is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}
