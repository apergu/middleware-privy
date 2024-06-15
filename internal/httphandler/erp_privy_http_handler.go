package httphandler

import (
	"encoding/json"
	"fmt"
	"middleware/internal/helper"
	"middleware/internal/model"
	"middleware/internal/usecase"
	"middleware/pkg/pkgvalidator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
)

type ErpPrivyHttpHandler struct {
	Command  *usecase.ErpPrivyCommandUsecaseGeneral
	Decorder rdecoder.Decoder
}

func NewErpPrivyHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.ErpPrivyUsecaseProperty{
		ErpPrivy: prop.DefaultERPPrivy,
	}

	handler := ErpPrivyHttpHandler{
		Command:  usecase.NewErpPrivyCommandUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/topup-balance", handler.TopUpBalance)
	r.Post("/check-status", handler.CheckTopUpStatus)
	r.Post("/void-balance", handler.VoidBalance)
	r.Post("/topup-adendum", handler.Adendum)
	r.Post("/topup-reconcile", handler.Reconcile)

	return r
}

func (h ErpPrivyHttpHandler) TopUpBalance(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error

	var ctx = r.Context()

	var payload model.TopUpBalance

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.TopUpBalance",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.TopUpBalance",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.TopUpBalance",
					nil,
				)
			}
		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.TopUpBalance",
				nil,
			)

		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))

		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.TopUpBalance",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		errorResponse := map[string]interface{}{
			"code":    422,
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		}

		// Convert error response to JSON
		responseJSON, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			// Handle JSON marshaling error
			fmt.Println("Error encoding JSON:", marshalErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // Set the appropriate HTTP status code

		// Write the JSON response to the client
		_, writeErr := w.Write(responseJSON)
		if writeErr != nil {
			// Handle write error
			fmt.Println("Error writing response:", writeErr)
		}

		return
	}

	res, err := h.Command.TopUpBalance(ctx, payload)

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to check top up balance",
				"at":     "ErpPrivyHttpHandler.TopUpBalance",
				"src":    "h.Command.TopUpBalance",
			}).
			Error(err)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "CheckTopUpStatus successfully", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h ErpPrivyHttpHandler) CheckTopUpStatus(w http.ResponseWriter, r *http.Request) {
	var err error

	var ctx = r.Context()

	var payload model.CheckTopUpStatus

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.CheckTopUpStatus",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.CheckTopUpStatus",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.CheckTopUpStatus",
					nil,
				)
			}

		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.CheckTopUpStatus",
				nil,
			)
		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))

		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.CheckTopUpStatus",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		errorResponse := map[string]interface{}{
			"code":    422,
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		}

		// Convert error response to JSON
		responseJSON, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			// Handle JSON marshaling error
			fmt.Println("Error encoding JSON:", marshalErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // Set the appropriate HTTP status code

		// Write the JSON response to the client
		_, writeErr := w.Write(responseJSON)
		if writeErr != nil {
			// Handle write error
			fmt.Println("Error writing response:", writeErr)
		}

		return
	}

	res, err := h.Command.CheckTopUpStatus(ctx, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	responseOk, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Channel successfully created", map[string]interface{}{
		"code":    200,
		"data":    res,
		"message": "CheckTopUpStatus successfully",
		"success": true,
	})
	helper.WriteJSONResponse(w, responseOk, http.StatusCreated)
}

func (h ErpPrivyHttpHandler) VoidBalance(w http.ResponseWriter, r *http.Request) {
	var err error

	var ctx = r.Context()

	var payload model.VoidBalance

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.VoidBalance",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.VoidBalance",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.VoidBalance",
					nil,
				)
			}

		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.VoidBalance",
				nil,
			)

		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))

		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.VoidBalance",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		errorResponse := map[string]interface{}{
			"code":    422,
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		}

		// Convert error response to JSON
		responseJSON, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			// Handle JSON marshaling error
			fmt.Println("Error encoding JSON:", marshalErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // Set the appropriate HTTP status code

		// Write the JSON response to the client
		_, writeErr := w.Write(responseJSON)
		if writeErr != nil {
			// Handle write error
			fmt.Println("Error writing response:", writeErr)
		}

		return
	}

	res, resPrivy, err := h.Command.VoidBalance(ctx, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	if err != nil {
		helper.WriteJSONResponse(w, res, helper.GetErrorStatusCode(err))
		return
	}

	responseOk, _ := helper.GenerateJSONResponse(http.StatusOK, true, "VoidBalance successfully", resPrivy)
	helper.WriteJSONResponse(w, responseOk, http.StatusOK)
}

func (h ErpPrivyHttpHandler) Adendum(w http.ResponseWriter, r *http.Request) {
	var err error

	var ctx = r.Context()

	var payload model.Adendum

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.Adendum",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.Adendum",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.Adendum",
					nil,
				)
			}

		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.Adendum",
				nil,
			)

		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))

		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.Adendum",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		errorResponse := map[string]interface{}{
			"code":    422,
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		}

		// Convert error response to JSON
		responseJSON, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			// Handle JSON marshaling error
			fmt.Println("Error encoding JSON:", marshalErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // Set the appropriate HTTP status code

		// Write the JSON response to the client
		_, writeErr := w.Write(responseJSON)
		if writeErr != nil {
			// Handle write error
			fmt.Println("Error writing response:", writeErr)
		}

		return
	}

	res, resPrivy, err := h.Command.Adendum(ctx, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	if err != nil {
		helper.WriteJSONResponse(w, res, helper.GetErrorStatusCode(err))
		return
	}

	responseOk, _ := helper.GenerateJSONResponse(http.StatusOK, true, "Adendum successfully created", resPrivy)
	helper.WriteJSONResponse(w, responseOk, http.StatusOK)
}

func (h ErpPrivyHttpHandler) Reconcile(w http.ResponseWriter, r *http.Request) {
	var err error

	var ctx = r.Context()

	var payload model.Reconcile

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.Reconcile",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.Reconcile",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.Reconcile",
					nil,
				)
			}

		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.Reconcile",
				nil,
			)

		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))

		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.Reconcile",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		errorResponse := map[string]interface{}{
			"code":    422,
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		}

		// Convert error response to JSON
		responseJSON, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			// Handle JSON marshaling error
			fmt.Println("Error encoding JSON:", marshalErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // Set the appropriate HTTP status code

		// Write the JSON response to the client
		_, writeErr := w.Write(responseJSON)
		if writeErr != nil {
			// Handle write error
			fmt.Println("Error writing response:", writeErr)
		}

		return
	}

	res, respPrivy, err := h.Command.Reconcile(ctx, payload)
	if err != nil {
		helper.WriteJSONResponse(w, res, helper.GetErrorStatusCode(err))
		return
	}

	responseOk, _ := helper.GenerateJSONResponse(http.StatusOK, true, "Reconcile successfully created", respPrivy)
	helper.WriteJSONResponse(w, responseOk, http.StatusOK)
}
