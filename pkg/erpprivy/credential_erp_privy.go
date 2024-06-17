package erpprivy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"golang.org/x/exp/slices"
)

func defaultSuccessCode() []int {
	return []int{200, 208, 201}
}

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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy TopUpBalance",
			"CredentialERPPrivy.TopUpBalance",
			err.Error(),
		)

		return err.Error(), errs
	}
	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.TopUpBalance",
					"src": "TopUpBalanceFailedResponse{}",
				}).
				Error(err)

			var resp TopUpBalanceFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeTopUpBalanceFailedResponse401",
						"at":     "ERPPrivy.TopUpBalance",
						"src":    "TopUpBalanceBadRequestResponse{}",
					}).
					Error(err)
				return TopUpBalanceResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.TopUpBalance",
				"Unauthorized",
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy TopUpBalance validation failed",
				"CredentialERPPrivy.TopUpBalance",
				strErr,
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy top up balance unknown code",
				"CredentialERPPrivy.TopUpBalance",
				resp,
			)

			return resp, err
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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy checktopupstatus",
			"CredentialERPPrivy.CheckTopUpStatus",
			err.Error(),
		)

		return err.Error(), errs
	}

	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.CheckTopUpStatus",
					"src": "CheckTopUpStatusFailedResponse{}",
				}).
				Error(err)

			var resp CheckTopUpStatusFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeCheckTopUpStatusFailedResponse401",
						"at":     "ERPPrivy.CheckTopUpStatus",
						"src":    "CheckTopUpStatusBadRequestResponse{}",
					}).
					Error(err)
				return CheckTopUpStatusResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.CheckTopUpStatus",
				"Unauthorized",
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy CheckTopUpStatus validation failed",
				"CredentialERPPrivy.CheckTopUpStatus",
				strErr,
			)

			return resp, err
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy void balance unknown code",
				"CredentialERPPrivy.VoidBalance",
				resp,
			)

			return resp, err
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

	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
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
				"CredentialERPPrivy.VoidBalance",
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
				"request erp privy VoidBalance validation failed",
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy void balance unknown code",
				"CredentialERPPrivy.VoidBalance",
				resp,
			)

			return resp, err
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

	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
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

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy adendum unknown code",
				"CredentialERPPrivy.Adendum",
				resp,
			)

			return resp, err
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

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy reconcile",
			"CredentialERPPrivy.Reconcile",
			err.Error(),
		)

		return err.Error(), errs
	}

	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
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
				"request erp privy reconcile unknown code",
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

func (c *CredentialERPPrivy) TransferBalanceERP(ctx context.Context, param TransferBalanceERPParam) (interface{}, error) {
	TransferBalanceURL := c.host + EndpointTransferBalance

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "ERPPrivy.TransferBalance",
			"src":  "TransferBalance{}.beforeDo",
			"host": TransferBalanceURL,
		}).
		Info(body.String())

	req, _ := http.NewRequest(http.MethodPost, TransferBalanceURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lang", "en")
	req.Header.Set("X-Request-Id", c.requestid)
	req.Header.Set("Application-Key", c.applicationkey)
	req.SetBasicAuth(c.username, c.password)

	resp := TransferBalanceERPResponse{}
	http := &http.Client{}
	res, err := http.Do(req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.TransferBalanceERP",
				"src": "TransferBalanceERP{}",
			}).
			Error(err)

		errs := rapperror.ErrInternalServerError(
			"",
			"failed to request erp privy TransferBalanceERP",
			"CredentialERPPrivy.TransferBalanceERP",
			err.Error(),
		)

		return err.Error(), errs
	}

	if !slices.Contains(defaultSuccessCode(), res.StatusCode) {
		var strErr string
		switch res.StatusCode {
		case 401:
			logrus.
				WithFields(logrus.Fields{
					"at":  "ERPPrivy.TransferBalanceERP",
					"src": "TransferBalanceERPFailedResponse{}",
				}).
				Error(err)
			var resp TransferBalanceERPFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeTransferBalanceERPFailedResponse",
						"at":     "ERPPrivy.TransferBalanceERP",
						"src":    "TransferBalanceERPBadRequestResponse{}",
					}).
					Error(err)
				return TransferBalanceERPResponse{}, err
			}

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy unauthorized",
				"CredentialERPPrivy.TransferBalanceERP",
				"Unauthorized",
			)

			return resp, err
		case 422:
			var resp TransferBalanceERPBadRequestResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"at":  "ERPPrivy.TransferBalanceERP",
						"src": "TransferBalanceERPBadRequestResponse{}",
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
				"CredentialERPPrivy.TransferBalanceERP",
				strErr,
			)

			return resp, err
		default:
			var resp TransferBalanceERPFailedResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			if err != nil {
				logrus.
					WithFields(logrus.Fields{
						"action": "DecodeTransferBalanceERPFailedResponse",
						"at":     "ERPPrivy.TransferBalanceERP",
						"src":    "TransferBalanceERPBadRequestResponse{}",
					}).
					Error(err)
				return TransferBalanceERPResponse{}, err
			}

			logrus.
				WithFields(logrus.Fields{
					"action": "GetResponseNot200Privy",
					"at":     "ERPPrivy.TransferBalanceERP",
					"src":    "TransferBalanceERPBadRequestResponse{}",
				}).
				Error(resp)

			err = rapperror.ErrInternalServerError(
				"",
				"request erp privy TransferBalance unknown code",
				"CredentialERPPrivy.TransferBalanceERP",
				resp,
			)

			return resp, err
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ERPPrivy.TransferBalanceERP",
				"src": "TransferBalanceERPResponse{}",
			}).
			Error(err)

		return TransferBalanceERPResponse{}, err
	}

	return resp, nil
}
