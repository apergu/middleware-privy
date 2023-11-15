package httphandler

import (
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/constants"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/internal/usecase"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
	"net/http"
)

type LeadHttpHandler struct {
	Command  usecase.LeadCommandUsecase
	Query    usecase.LeadQueryUsecase
	Decorder rdecoder.Decoder
}

func NewLeadHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.LeadUsecaseProperty{
		LeadRepo:  repository.NewLeadRepositoryPostgre(prop.DBPool),
		LeadPrivy: prop.DefaultCredential,
	}

	handler := LeadHttpHandler{
		Command:  usecase.NewLeadQueryUsecaseGeneral(ucProp),
		Query:    usecase.NewLeadQueryUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/", handler.Create)

	return r
}

func (h LeadHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Leads

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "LeadHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"LeadHttpHandler.Create",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	// get user from context
	user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = user

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

	roleId, meta, err := h.Command.Create(ctx, payload)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}
