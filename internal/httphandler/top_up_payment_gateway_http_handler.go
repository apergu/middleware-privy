package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"middleware/infrastructure"
	request "middleware/infrastructure/http/request"
	"middleware/internal/helper"
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

	// payloadLines := []request.LineItem{}

	// for _, v := range payload.Lines {
	// 	payloadLines = append(payloadLines, request.LineItem{
	// 		Item:                         v.Item,
	// 		CustColPrivyMerchant:         v.CustColPrivyMerchant,
	// 		CustColPrivyChannel:          v.CustColPrivyChannel,
	// 		TaxCode:                      v.TaxCode,
	// 		CustColPrivyMainProduct:      v.CustColPrivyMainProduct,
	// 		CustColPrivySubProduct:       v.CustColPrivySubProduct,
	// 		Description:                  v.Description,
	// 		Quantity:                     v.Quantity,
	// 		CustColPrivyStartDateLayanan: v.CustColPrivyStartDateLayanan,
	// 		CustColPrivyDateLayanan:      v.CustColPrivyDateLayanan,
	// 		CustColPrivyTrxID:            v.CustColPrivyTrxID,
	// 		CustColPaymentGatewayFee:     v.CustColPaymentGatewayFee,
	// 		Amount:                       v.Amount,
	// 		CustColAmountBeforeDisc:      v.CustColAmountBeforeDisc,
	// 	})
	// }

	// payloadReq := request.PaymentGateway{
	// 	RecordType:                       "salesorder",
	// 	CustomForm:                       "144",
	// 	CustBodyPrivySoCustID:            payload.CustBodyPrivySoCustID,
	// 	Entity:                           payload.Entity,
	// 	StartDate:                        payload.StartDate,
	// 	EndDate:                          payload.EndDate,
	// 	CustBodyPrivyTermOfPayment:       payload.CustBodyPrivyTermOfPayment,
	// 	OtherRefNum:                      payload.OtherRefNum,
	// 	CustBodyPrivyBilling:             payload.CustBodyPrivyBilling,
	// 	CustBodyPrivyIntegrasi:           payload.CustBodyPrivyIntegrasi,
	// 	Memo:                             payload.Memo,
	// 	CustBodyPrivyBDA:                 payload.CustBodyPrivyBDA,
	// 	CustBodyPrivyBDM:                 payload.CustBodyPrivyBDM,
	// 	CustBodyPrivySalesSupport:        payload.CustBodyPrivySalesSupport,
	// 	CustBodyPrivySalesSupportManager: payload.CustBodyPrivySalesSupportManager,
	// 	CustBody10:                       payload.CustBody10,
	// 	CustBody9:                        payload.CustBody9,
	// 	CustBody7:                        payload.CustBody7,
	// 	Lines:                            payloadLines,
	// }

	// resTopUp, err := h.Command.TopUpPaymentGateWay(payloadReq)
	_, err = h.Command.TopUpPaymentGateWay(payload)
	if err != nil {
		helper.LoggerErrorStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", err.Error(), "", payload, nil)
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	res, _ := helper.GenerateJSONResponse(http.StatusCreated, true, "TopUpPayment successfully created", map[string]interface{}{})
	helper.LoggerSuccessStructfunc(w, r, "TOP_UP_PAYMENT_GATEWAY", "TopUpPaymentGateWay", "TopUpPayment successfully created", "")
	helper.WriteJSONResponse(w, res, http.StatusCreated)
}
