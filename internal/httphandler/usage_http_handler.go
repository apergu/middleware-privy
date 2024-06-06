package httphandler

import (
	"net/http"

	"middleware/internal/helper"
	"middleware/internal/model"
	"middleware/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
)

type UsageUsecaseHttpHandler struct {
	Command  usecase.UsageUsecases
	Decorder rdecoder.Decoder
}

func NewUsageUsecaseHttpHandler(prop HTTPHandlerProperty) http.Handler {

	ucProp := usecase.CustomerUsageUsecaseProperty{
		CustomerPrivy: prop.DefaultCredential,
	}

	handler := UsageUsecaseHttpHandler{
		Command:  usecase.NewUsageUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()
	r.Post("/", handler.Create)

	return r
}

func (h UsageUsecaseHttpHandler) Create(w http.ResponseWriter, r *http.Request) {

	var err error
	ctx := r.Context()
	var payload model.UsageModel

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "UsageUsecaseHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"UsageUsecaseHttpHandler.Create",
			nil,
		)

		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	_, err = h.Command.Create(ctx, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	res, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Usage successfully created", map[string]interface{}{})

	helper.WriteJSONResponse(w, res, http.StatusCreated)
}
