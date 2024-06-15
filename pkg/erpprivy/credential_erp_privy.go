package erpprivy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
)

func (c *CredentialERPPrivy) TopUpBalance(ctx context.Context, param TopUpBalanceParam) (interface{}, error) {
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

	if res.StatusCode != 200 && res.StatusCode != 208 {
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
						"action": "DecodeTopUpBalanceFailedResponse",
						"at":     "ERPPrivy.TopUpBalance",
						"src":    "TopUpBalanceBadRequestResponse{}",
					}).
					Error(err)
				return TopUpBalanceResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.TopUpBalance",
					"src":    "TopUpBalanceBadRequestResponse{}",
				}).
				Error(resp)

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

func (c *CredentialERPPrivy) CheckTopUpStatus(ctx context.Context, param CheckTopUpStatusParam) (interface{}, error) {
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

	if res.StatusCode != 200 && res.StatusCode != 208 {
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
						"action": "DecodeCheckTopUpStatusFailedResponse",
						"at":     "ERPPrivy.CheckTopUpStatus",
						"src":    "CheckTopUpStatusBadRequestResponse{}",
					}).
					Error(err)
				return CheckTopUpStatusResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.CheckTopUpStatus",
					"src":    "CheckTopUpStatusBadRequestResponse{}",
				}).
				Error(resp)

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

func (c *CredentialERPPrivy) VoidBalance(ctx context.Context, param VoidBalanceParam) (interface{}, error) {
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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy void balance",
			"CredentialERPPrivy.VoidBalance",
			err.Error(),
		)

		return err.Error(), errs
	}

	if res.StatusCode != 200 && res.StatusCode != 208 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.VoidBalance",
					"src": "VoidBalanceFailedResponse{}",
				}).
				Error(err)

			var resp VoidBalanceFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeVoidBalanceFailedResponse401",
						"at":     "ERPPrivy.VoidBalance",
						"src":    "VoidBalanceBadRequestResponse{}",
					}).
					Error(err)
				return VoidBalanceResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.Adendum",
				"Unauthorized",
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy validation failed",
				"CredentialERPPrivy.VoidBalance",
				strErr,
			)

			return resp, err
		default:
			var resp VoidBalanceFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeVoidBalanceFailedResponse",
						"at":     "ERPPrivy.VoidBalance",
						"src":    "VoidBalanceBadRequestResponse{}",
					}).
					Error(err)
				return VoidBalanceResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.VoidBalance",
					"src":    "VoidBalanceBadRequestResponse{}",
				}).
				Error(resp)

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

func (c *CredentialERPPrivy) Adendum(ctx context.Context, param AdendumParam) (interface{}, error) {
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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy adendum",
			"CredentialERPPrivy.Adendum",
			err.Error(),
		)

		return err.Error(), errs
	}

	if res.StatusCode != 200 && res.StatusCode != 208 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.Adendum",
					"src": "AdendumFailedResponse{}",
				}).
				Error(err)

			var resp AdendumFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeAdendumFailedResponse401",
						"at":     "ERPPrivy.Adendum",
						"src":    "AdendumBadRequestResponse{}",
					}).
					Error(err)
				return AdendumResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.Adendum",
				"Unauthorized",
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy validation failed",
				"CredentialERPPrivy.Adendum",
				strErr,
			)

			return resp, err
		default:
			var resp AdendumFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeAdendumFailedResponse",
						"at":     "ERPPrivy.Adendum",
						"src":    "AdendumBadRequestResponse{}",
					}).
					Error(err)
				return AdendumResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.Adendum",
					"src":    "AdendumBadRequestResponse{}",
				}).
				Error(resp)

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

func (c *CredentialERPPrivy) Reconcile(ctx context.Context, param ReconcileParam) (interface{}, error) {
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
	// req.Header.Set("X-Request-Id", c.requestid)
	// req.Header.Set("Application-Key", c.applicationkey)
	// req.SetBasicAuth(c.username, c.password)

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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy reconcile",
			"CredentialERPPrivy.Reconcile",
			err.Error(),
		)

		return err.Error(), errs
	}

	if res.StatusCode != 200 && res.StatusCode != 208 {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.Reconcile",
					"src": "ReconcileFailedResponse{}",
				}).
				Error(err)
			var resp ReconcileFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeReconcileFailedResponse",
						"at":     "ERPPrivy.Reconcile",
						"src":    "ReconcileBadRequestResponse{}",
					}).
					Error(err)
				return ReconcileResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.Reconcile",
				"Unauthorized",
			)

			return resp, err
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
				return resp, errors.New(resp.Message)
			}

			for _, v := range resp.Errors {
				strErr += v.Field + " " + v.Description + " "
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy validation failed",
				"CredentialERPPrivy.Reconcile",
				strErr,
			)

			return resp, err
		default:
			var resp ReconcileFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeReconcileFailedResponse",
						"at":     "ERPPrivy.Reconcile",
						"src":    "ReconcileBadRequestResponse{}",
					}).
					Error(err)
				return ReconcileResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.Reconcile",
					"src":    "ReconcileBadRequestResponse{}",
				}).
				Error(resp)

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.Reconcile",
				resp,
			)

			return resp, err
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
