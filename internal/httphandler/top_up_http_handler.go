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

type TopUpHttpHandler struct {
	Command  *usecase.ErpPrivyCommandUsecaseGeneral
	Decorder rdecoder.Decoder
}

func NewTopUpHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.ErpPrivyUsecaseProperty{
		ErpPrivyDataPrivy: prop.DefaultPrivy,
		ErpPrivyPrivy:     prop.DefaultCredential,
	}

	handler := TopUpHttpHandler{
		Command:  usecase.NewErpPrivyCommandUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/check-status", handler.CheckTopUpStatus)
	r.Post("/void-balance", handler.VoidBalance)

	return r
}

func (h TopUpHttpHandler) CheckTopUpStatus(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error

	var ctx = r.Context()

	var payload model.CheckTopUpStatus

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "TopUpHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"TopUpHttpHandler.Create",
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
				"at":     "CustomerUsageHttpHandler.Create",
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
				"at":     "TopUpHttpHandler.Create",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"TopUpHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "CheckTopUpStatus successfully created", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h TopUpHttpHandler) VoidBalance(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error

	var ctx = r.Context()

	var payload model.VoidBalance

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "TopUpHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"TopUpHttpHandler.Create",
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
				"at":     "CustomerUsageHttpHandler.Create",
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
				"at":     "TopUpHttpHandler.Create",
				"src":    "h.Command.CheckTopUpStatus",
			}).
			Error(err)

		err = rapperror.ErrInternalServerError(
			rapperror.AppErrorCodeInternalServerError,
			"Internal server error",
			"TopUpHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "CheckTopUpStatus successfully created", nil, res)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}