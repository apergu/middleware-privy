package httphandler

import (
	"net/http"

	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	"middleware/internal/helper"
	service "middleware/internal/services/privy"
	usecase "middleware/internal/usecase/top_up_payment"

	"github.com/go-chi/chi/v5"
	"gitlab.com/rteja-library3/rdecoder"
)

type TopUpPaymentUsecaseHttpHandler struct {
	Command  usecase.TopUpPaymentUsecase
	Decorder rdecoder.Decoder
}

func NewTopUpPaymentUsecaseHttpHandler(prop HTTPHandlerProperty, inf *infrastructure.Infrastructure) http.Handler {

	ucProp := service.NewToNetsuitService(inf)
	handler := TopUpPaymentUsecaseHttpHandler{
		Command:  usecase.NewTopUpPaymentUsecaseGeneral(ucProp, inf),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()
	r.Post("/", handler.TopUpPayment)

	return r
}

func (h TopUpPaymentUsecaseHttpHandler) TopUpPayment(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload request.CustomerDetails

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {

		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	_, err = h.Command.TopUpPayment(payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	res, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "TopUpPayment successfully created", map[string]interface{}{})

	helper.WriteJSONResponse(w, res, http.StatusCreated)
}
