package credential

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

func (c *CredentialPrivy) CreateTopUp(ctx context.Context, param TopUpParam) (TopUpResponse, error) {

	log.Println("TEST", param)
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

		return TopUpResponse{}, err
	}

	// TOPUP TO PRIVY

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

		return TopUpResponse{}, err
	}

	// post customer
	postSalesOrderURL := c.host + EndpointMerchant

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateMerchant",
			"src":  "EnvelopeMerchant{}.beforeDo",
			"host": postSalesOrderURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, "https://stg-b2b-api-service.privy.id/v1/orchestrator-erp-goldengate/webhook/apergu/top-up", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic YXBlcmd1OnNlY3JldA==")
	req.Header.Set("application-key", "VUNSAT9GP6e5Rc7qv8ZDnh")
	req.Header.Set("application_creds_username", "apergu")
	req.Header.Set("application_creds_password", "2dp$m48k#ut9")
	req.Header.Set("application_creds_key", "VUNSAT9GP6e5Rc7qv8ZDnh")

	// q := req.URL.Query()
	// q.Add("script", "175")
	// q.Add("deploy", "1")

	// req.URL.RawQuery = q.Encode()

	custResp := EnvelopeTopUp{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateMerchant",
				"src": "EnvelopeMerchant{}",
			}).
			Error(err)

		return TopUpResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return TopUpResponse{}, rapperror.ErrNotFound(
			"",
			"Merchant is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}
