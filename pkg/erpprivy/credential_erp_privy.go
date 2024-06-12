package erpprivy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *CredentialERPPrivy) TopUpBalance(ctx context.Context, param TopUpBalanceParam) (TopUpBalanceResponse, error) {
	TopUpBalanceURL := c.host + EndpointTopUpBalance

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.TopUpBalance",
			"src":  "TopUpBalance{}.beforeDo",
			"host": TopUpBalanceURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, TopUpBalanceURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := TopUpBalanceResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.TopUpBalance",
				"src": "TopUpBalance{}",
			}).
			Error(err)

		return TopUpBalanceResponse{}, err
	}

	if res.StatusCode != 200 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.TopUpBalance",
					"src": "TopUpBalanceFailedResponse{}",
				}).
				Error(err)

			return TopUpBalanceResponse{}, errors.New("request erp privy unauthorized")
		case 422:
			var resp TopUpBalanceBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.TopUpBalance",
						"src": "TopUpBalanceBadRequestResponse{}",
					}).
					Error(err)
			}

			if resp.Errors == nil {
				return TopUpBalanceResponse{}, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			return TopUpBalanceResponse{}, errors.New(strErr)
		default:
			var resp TopUpBalanceFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.TopUpBalance",
						"src": "TopUpBalanceBadRequestResponse{}",
					}).
					Error(err)
				return TopUpBalanceResponse{}, err
			}

			return TopUpBalanceResponse{}, errors.New("something went wrong")
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.TopUpBalance",
				"src": "TopUpBalance{}",
			}).
			Error(err)

		return TopUpBalanceResponse{}, err
	}

	return resp, nil
}

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

func (c *CredentialERPPrivy) Adendum(ctx context.Context, param AdendumParam) (AdendumResponse, error) {
	AdendumURL := c.host + EndpointAdendum

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.Adendum",
			"src":  "Adendum{}.beforeDo",
			"host": AdendumURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPatch, AdendumURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := AdendumResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.Adendum",
				"src": "Adendum{}",
			}).
			Error(err)

		return AdendumResponse{}, err
	}

	if res.StatusCode != 200 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.Adendum",
					"src": "AdendumFailedResponse{}",
				}).
				Error(err)

			return AdendumResponse{}, errors.New("request erp privy unauthorized")
		case 422:
			var resp AdendumBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.Adendum",
						"src": "AdendumBadRequestResponse{}",
					}).
					Error(err)
			}

			if resp.Errors == nil {
				return AdendumResponse{}, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			return AdendumResponse{}, errors.New(strErr)
		default:
			var resp AdendumFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.Adendum",
						"src": "AdendumBadRequestResponse{}",
					}).
					Error(err)
				return AdendumResponse{}, err
			}

			return AdendumResponse{}, errors.New("something went wrong")
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.Adendum",
				"src": "AdendumResponse{}",
			}).
			Error(err)

		return AdendumResponse{}, err
	}

	return resp, nil
}

func (c *CredentialERPPrivy) Reconcile(ctx context.Context, param ReconcileParam) (ReconcileResponse, error) {
	ReconcileURL := c.host + EndpointReconcile

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.Reconcile",
			"src":  "Reconcile{}.beforeDo",
			"host": ReconcileURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPatch, ReconcileURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := ReconcileResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.Reconcile",
				"src": "Reconcile{}",
			}).
			Error(err)

		return ReconcileResponse{}, err
	}

	if res.StatusCode != 200 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.Reconcile",
					"src": "ReconcileFailedResponse{}",
				}).
				Error(err)

			return ReconcileResponse{}, errors.New("request erp privy unauthorized")
		case 422:
			var resp ReconcileBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.Reconcile",
						"src": "ReconcileBadRequestResponse{}",
					}).
					Error(err)
			}

			if resp.Errors == nil {
				return ReconcileResponse{}, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			return ReconcileResponse{}, errors.New(strErr)
		default:
			var resp ReconcileFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.Reconcile",
						"src": "ReconcileBadRequestResponse{}",
					}).
					Error(err)
				return ReconcileResponse{}, err
			}

			return ReconcileResponse{}, errors.New("something went wrong")
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.Reconcile",
				"src": "ReconcileResponse{}",
			}).
			Error(err)

		return ReconcileResponse{}, err
	}

	return resp, nil
}
