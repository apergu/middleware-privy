package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	"middleware/internal/helper"
	"middleware/internal/model"
	service "middleware/internal/services/privy"
	usecaseErpPrivy "middleware/internal/usecase"
	usecase "middleware/internal/usecase/top_up_payment"
	"middleware/pkg/pkgvalidator"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
)

type TopUpPaymentGateWayHttpHandler struct {
	Command      usecase.TopUpPaymentGateWayUsecase
	CommandTopUp *usecaseErpPrivy.ErpPrivyCommandUsecaseGeneral
	Decorder     rdecoder.Decoder
}

func NewTopUpPaymentGateWayHttpHandler(prop HTTPHandlerProperty, inf *infrastructure.Infrastructure) http.Handler {

	ucProp := service.NewToNetsuitService(inf)

	ucPropErp := usecaseErpPrivy.ErpPrivyUsecaseProperty{
		ErpPrivy: prop.DefaultERPPrivy,
	}

	handler := TopUpPaymentGateWayHttpHandler{
		Command:      usecase.NewTopUpPaymentGateWayGeneral(ucProp, inf),
		CommandTopUp: usecaseErpPrivy.NewErpPrivyCommandUsecaseGeneral(ucPropErp),
		Decorder:     prop.DefaultDecoder,
	}

	r := chi.NewRouter()
	r.Post("/", handler.TopUpPayment)

	return r
}

func (h TopUpPaymentGateWayHttpHandler) TopUpPayment(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload request.PaymentGateway

	xRequestId := r.Header.Get("X-Request-Id")
	if xRequestId == "" {
		logrus.
			WithFields(logrus.Fields{
				"at":  "ErpPrivyHttpHandler.TopUpPayment",
				"src": "payload.X-Request-Id",
			}).Error(errors.New("please provide X-Request-Id in header"))

		response, _ := helper.GenerateJSONResponse(422, false, "please provide X-Request-Id in header", map[string]interface{}{})
		helper.WriteJSONResponse(w, response, 422)
		return
	}

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ErpPrivyHttpHandler.TopUpPayment",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)
		switch jsonerr := err.(type) {
		case *json.UnmarshalTypeError:
			if jsonerr.Field == "" {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					"Invalid body",
					"ErpPrivyHttpHandler.TopUpPayment",
					nil,
				)
			} else {
				err = rapperror.ErrUnprocessableEntity(
					rapperror.AppErrorCodeUnprocessableEntity,
					fmt.Sprintf(jsonerr.Field+" must be a "+jsonerr.Type.String()),
					"ErpPrivyHttpHandler.TopUpPayment",
					nil,
				)
			}
		default:
			err = rapperror.ErrUnprocessableEntity(
				rapperror.AppErrorCodeUnprocessableEntity,
				"invalid body",
				"ErpPrivyHttpHandler.TopUpPayment",
				nil,
			)

		}

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	errors := pkgvalidator.Validate(payload)
	if len(errors) > 0 {
		var message string
		for _, v := range errors {
			if message == "" {
				message = v["description"].(string)
			} else {
				message = message + "; " + v["description"].(string)
			}
		}

		helper.LoggerValidateStructfunc(w, r, "PaymentGatewayHttpHandler.TopUpPaymentGateway", "ERPTopUpPaymentGateway", message, "", payload)

		logrus.
			WithFields(logrus.Fields{
				"at":     "PaymentGatewayHttpHandler.TopUpPaymentGateway",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(message)

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

	resTopUp, err := h.Command.TopUpPaymentGateWay(payload)
	if err != nil {
		helper.LoggerErrorStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", err.Error(), "", payload, nil)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	var resTopUpmodel request.ResPaymentGateway

	resByte, err := json.Marshal(resTopUp)
	if err != nil {
		helper.LoggerErrorStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", err.Error(), "", payload, nil)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	err = json.Unmarshal(resByte, &resTopUpmodel)
	if err != nil {
		helper.LoggerErrorStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", err.Error(), "", payload, nil)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	reqTopUpBalance := model.TopUpBalance{
		TopUPID:         resTopUpmodel.Data.TopupID,
		EnterpriseId:    resTopUpmodel.Data.EnterpriseID,
		MerchantId:      resTopUpmodel.Data.MerchantID,
		ChannelId:       resTopUpmodel.Data.ChannelID,
		ServiceId:       resTopUpmodel.Data.ServiceID,
		PostPaid:        resTopUpmodel.Data.PostPaid,
		Qty:             resTopUpmodel.Data.Qty,
		UnitPrice:       resTopUpmodel.Data.UnitPrice,
		StartPeriodDate: resTopUpmodel.Data.StartPeriodDate.Format("02/01/2006"),
		EndPeriodDate:   resTopUpmodel.Data.EndPeriodDate.Format("02/01/2006"),
		TransactionDate: resTopUpmodel.Data.TransactionDate.Format("02/01/2006"),
	}

	_, _, err = h.CommandTopUp.TopUpBalance(context.Background(), reqTopUpBalance, xRequestId)
	if err != nil {
		helper.LoggerErrorStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWayWithTopUpBalance", err.Error(), "", reqTopUpBalance, nil)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	res, _ := helper.GenerateJSONResponse(http.StatusCreated, true, "TopUpPayment successfully created", map[string]interface{}{
		"trx_id": resTopUpmodel.Data.TopupID,
	})
	helper.LoggerSuccessStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", "TopUpPayment successfully created", "")
	helper.WriteJSONResponse(w, res, http.StatusCreated)
}
