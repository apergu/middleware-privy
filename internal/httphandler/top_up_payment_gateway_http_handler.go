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

type TopUpPaymentGateWayHttpHandler struct {
	Command  usecase.TopUpPaymentGateWayUsecase
	Decorder rdecoder.Decoder
}

func NewTopUpPaymentGateWayHttpHandler(prop HTTPHandlerProperty, inf *infrastructure.Infrastructure) http.Handler {

	ucProp := service.NewToNetsuitService(inf)
	handler := TopUpPaymentGateWayHttpHandler{
		Command:  usecase.NewTopUpPaymentGateWayGeneral(ucProp, inf),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()
	r.Post("/", handler.TopUpPayment)

	return r
}

func (h TopUpPaymentGateWayHttpHandler) TopUpPayment(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload request.PaymentGateway

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {

		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	_, err = h.Command.TopUpPaymentGateWay(payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	res, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "TopUpPayment successfully created", map[string]interface{}{})

	helper.WriteJSONResponse(w, res, http.StatusCreated)
}
