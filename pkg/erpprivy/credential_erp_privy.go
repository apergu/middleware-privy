package erpprivy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *CredentialERPPrivy) CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (CheckTopUpStatusResponse, error) {
	checkTopUpStatusURL := c.host + EndpointCheckTopUpStatus

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.CheckTopUpStatus",
			"src":  "CheckTopUpStatus{}.beforeDo",
			"host": checkTopUpStatusURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, checkTopUpStatusURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := CheckTopUpStatusResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.CheckTopUpStatus",
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
					"at":  "ERPPrivy.CheckTopUpStatus",
					"src": "CheckTopUpStatusFailedResponse{}",
				}).
				Error(err)

			return CheckTopUpStatusResponse{}, errors.New("request erp privy unauthorized")
		case 422:
			var resp CheckTopUpStatusBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.CheckTopUpStatus",
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
						"at":  "ERPPrivy.CheckTopUpStatus",
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
				"at":  "ERPPrivy.CheckTopUpStatus",
				"src": "CheckTopUpStatus{}",
			}).
			Error(err)

		return CheckTopUpStatusResponse{}, err
	}

	return resp, nil
}

func (c *CredentialERPPrivy) VoidBalance(ctx context.Context, param VoidBalanceParam) (VoidBalanceResponse, error) {
	VoidBalanceURL := c.host + EndpointVoidBalance

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.VoidBalance",
			"src":  "VoidBalance{}.beforeDo",
			"host": VoidBalanceURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, VoidBalanceURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := VoidBalanceResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.VoidBalance",
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
					"at":  "ERPPrivy.VoidBalance",
					"src": "VoidBalanceFailedResponse{}",
				}).
				Error(err)

			return VoidBalanceResponse{}, errors.New("request erp privy unauthorized")
		case 422:
			var resp VoidBalanceBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.VoidBalance",
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
						"at":  "ERPPrivy.VoidBalance",
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
				"at":  "ERPPrivy.VoidBalance",
				"src": "VoidBalance{}",
			}).
			Error(err)

		return VoidBalanceResponse{}, err
	}

	return resp, nil
}
