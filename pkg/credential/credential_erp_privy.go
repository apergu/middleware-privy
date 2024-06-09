package credential

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func (c *CredentialPrivy) CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (CheckTopUpStatusResponse, error) {
	checkTopUpStatusURL := c.host + EndpointCheckTopUpStatus

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CheckTopUpStatus",
			"src":  "CheckTopUpStatus{}.beforeDo",
			"host": checkTopUpStatusURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, checkTopUpStatusURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", "7c7bab49266d3529254f2532fe7cff8e")
	req.Header.Set("Application-Key", "VUNSAT9GP6e5Rc7qv8ZDnh")
	req.SetBasicAuth(c.username, c.password)

	resp := CheckTopUpStatusResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CheckTopUpStatus",
				"src": "CheckTopUpStatus{}",
			}).
			Error(err)

		return CheckTopUpStatusResponse{}, err
	}

	if res.StatusCode != 200 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "CredentialPrivy.CheckTopUpStatus",
					"src": "CheckTopUpStatusFailedResponse{}",
				}).
				Error(err)

			return CheckTopUpStatusResponse{}, errors.New("request privy unauthorized")
		case 422:
			var resp CheckTopUpStatusBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "CredentialPrivy.CheckTopUpStatus",
						"src": "CheckTopUpStatusBadRequestResponse{}",
					}).
					Error(err)
			}

			if resp.Errors == nil {
				return CheckTopUpStatusResponse{}, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			return CheckTopUpStatusResponse{}, errors.New(strErr)
		default:
			var resp CheckTopUpStatusFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "CredentialPrivy.CheckTopUpStatus",
						"src": "CheckTopUpStatusBadRequestResponse{}",
					}).
					Error(err)
				return CheckTopUpStatusResponse{}, err
			}

			return CheckTopUpStatusResponse{}, errors.New("something went wrong")
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CheckTopUpStatus",
				"src": "CheckTopUpStatus{}",
			}).
			Error(err)

		return CheckTopUpStatusResponse{}, err
	}

	return resp, nil
}

func (c *CredentialPrivy) VoidBalance(ctx context.Context, param VoidBalanceParam) (VoidBalanceResponse, error) {
	VoidBalanceURL := c.host + EndpointVoidBalance

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.VoidBalance",
			"src":  "VoidBalance{}.beforeDo",
			"host": VoidBalanceURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, VoidBalanceURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", "7c7bab49266d3529254f2532fe7cff8e")
	req.Header.Set("Application-Key", "VUNSAT9GP6e5Rc7qv8ZDnh")
	req.SetBasicAuth(c.username, c.password)

	resp := VoidBalanceResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.VoidBalance",
				"src": "VoidBalance{}",
			}).
			Error(err)

		return VoidBalanceResponse{}, err
	}

	if res.StatusCode != 200 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "CredentialPrivy.VoidBalance",
					"src": "VoidBalanceFailedResponse{}",
				}).
				Error(err)

			return VoidBalanceResponse{}, errors.New("request privy unauthorized")
		case 422:
			var resp VoidBalanceBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "CredentialPrivy.VoidBalance",
						"src": "VoidBalanceBadRequestResponse{}",
					}).
					Error(err)
			}

			if resp.Errors == nil {
				return VoidBalanceResponse{}, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			return VoidBalanceResponse{}, errors.New(strErr)
		default:
			var resp VoidBalanceFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "CredentialPrivy.VoidBalance",
						"src": "VoidBalanceBadRequestResponse{}",
					}).
					Error(err)
				return VoidBalanceResponse{}, err
			}

			return VoidBalanceResponse{}, errors.New("something went wrong")
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.VoidBalance",
				"src": "VoidBalance{}",
			}).
			Error(err)

		return VoidBalanceResponse{}, err
	}

	return resp, nil
}
