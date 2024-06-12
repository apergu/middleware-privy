package httphandler

import (
	"encoding/json"
	"fmt"
	"middleware/internal/model"
	"middleware/internal/usecase"
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
				"at":     "ErpPrivyHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ErpPrivyHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	errors := payload.Validate()
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.Create",
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
				"at":     "ErpPrivyHttpHandler.Create",
				"src":    "h.Command.TopUpBalance",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"ErpPrivyHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "CheckTopUpStatus successfully created", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h ErpPrivyHttpHandler) CheckTopUpStatus(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error

	var ctx = r.Context()

	var payload model.CheckTopUpStatus

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ErpPrivyHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	errors := payload.Validate()
	if len(errors) > 0 {
		logrus.
			WithFields(logrus.Fields{
				"at":     "ErpPrivyHttpHandler.Create",
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
		logrus.
			WithFields(logrus.Fields{
				"action": "try to check top up status",
				"at":     "ErpPrivyHttpHandler.Create",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"ErpPrivyHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "CheckTopUpStatus successfully", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h ErpPrivyHttpHandler) VoidBalance(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error

	var ctx = r.Context()

	var payload model.VoidBalance

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ErpPrivyHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	errors := payload.Validate()
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

	res, err := h.Command.VoidBalance(ctx, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to check top up status",
				"at":     "ErpPrivyHttpHandler.VoidBalance",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"ErpPrivyHttpHandler.VoidBalance",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "VoidBalance successfully", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h ErpPrivyHttpHandler) Adendum(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
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

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ErpPrivyHttpHandler.Adendum",
			err.Error(),
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	errors := payload.Validate()
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

	res, err := h.Command.Adendum(ctx, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to check adendum",
				"at":     "ErpPrivyHttpHandler.Adendum",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"ErpPrivyHttpHandler.Adendum",
			err.Error(),
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Adendum successfully", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h ErpPrivyHttpHandler) Reconcile(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
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

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ErpPrivyHttpHandler.Reconcile",
			err.Error(),
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	errors := payload.Validate()
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

	res, err := h.Command.Reconcile(ctx, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to check reconcile",
				"at":     "ErpPrivyHttpHandler.Reconcile",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"ErpPrivyHttpHandler.Reconcile",
			err.Error(),
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Reconcile successfully", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}
