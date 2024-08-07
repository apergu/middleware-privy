package httphandler

import (
	"encoding/json"
	"fmt"

	// "middleware/infrastructure/logger/logrus"
	"middleware/internal/helper"

	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/internal/usecase"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rhelper"
	"gitlab.com/rteja-library3/rresponser"
)

type ApplicationHttpHandler struct {
	Command  usecase.ApplicationCommandUsecase
	Query    usecase.ApplicationQueryUsecase
	Decorder rdecoder.Decoder
}

func NewApplicationHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.ApplicationUsecaseProperty{
		ApplicationPrivy: prop.DefaultCredential,
		CustomerRepo:     repository.NewCustomerRepositoryPostgre(prop.DBPool),
		ApplicationRepo:  repository.NewApplicationRepositoryPostgre(prop.DBPool),
	}

	handler := ApplicationHttpHandler{
		Command:  usecase.NewApplicationCommandUsecaseGeneral(ucProp),
		Query:    usecase.NewApplicationQueryUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/", handler.Create)
	r.Put("/id/{id}", handler.Update)
	r.Delete("/id/{id}", handler.Delete)

	r.Get("/", handler.Find)
	r.Get("/id/{id}", handler.FindById)

	return r
}

func (h ApplicationHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Application

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)

	if err != nil {
		msg := err.Error()
		re := regexp.MustCompile(`Application\.(\w+)`)
		custm := re.FindStringSubmatch(msg)
		re = regexp.MustCompile(`([a-z])([A-Z])`)
		spaced := re.ReplaceAllString(custm[1], `$1 $2`)
		re = regexp.MustCompile(`type ([^\]]+)`)
		format := re.FindStringSubmatch(msg)
		message := fmt.Sprintf("Unprocessable entity - %s value must in %s format", spaced, format[1])

		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "MerchantHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrUnprocessableEntity(
			"",
			message,
			"CustomerHttpHandler.Create",
			nil,
		)

		response, _ := helper.GenerateJSONResponse(422, false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}

	// if err != nil {
	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"action": "try to decode data",
	// 			"at":     "ApplicationHttpHandler.Create",
	// 			"src":    "rdecoder.DecodeRest",
	// 		}).
	// 		Error(err)

	// 	err = rapperror.ErrBadRequest(
	// 		rapperror.AppErrorCodeBadRequest,
	// 		"Invalid body",
	// 		"ApplicationHttpHandler.Create",
	// 		nil,
	// 	)

	// 	response = rresponser.NewResponserError(err)
	// 	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	// 	return
	// }

	// get user from context
	// user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = 0

	errors := payload.Validate()
	if len(errors) > 0 {
		var message string
		for _, v := range errors {
			if message == "" {
				message = v["Description"].(string)
			} else {
				message = message + "; " + v["Description"].(string)
			}
		}

		helper.LoggerValidateStructfunc(w, r, "CustomerHttpHandler.Create", "application", message, "", payload)

		logrus.
			WithFields(logrus.Fields{
				"at":     "CustomerHttpHandler.Create",
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
	println("payload =>", payload.Address)
	roleId, meta, err := h.Command.Create(ctx, payload)

	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}
	response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Application fail created", map[string]interface{}{
		"roleId": roleId,
		"meta":   meta,
	})

	helper.WriteJSONResponse(w, response, http.StatusCreated)
	helper.LoggerSuccessStructfunc(w, r, "CustomerUsageHttpHandler.Create", "Application", "Application successfully created", "")
}

func (h ApplicationHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Application

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ApplicationHttpHandler.Update",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ApplicationHttpHandler.Update",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ApplicationHttpHandler.Update",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	// get user from context
	// user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = 0

	errors := payload.Validate()
	if len(errors) > 0 {
		var message string
		for _, v := range errors {
			if message == "" {
				message = v["Description"].(string)
			} else {
				message = message + "; " + v["Description"].(string)
			}
		}

		helper.LoggerValidateStructfunc(w, r, "CustomerHttpHandler.Update", "channel", message, "", payload)

		logrus.
			WithFields(logrus.Fields{
				"at":     "CustomerHttpHandler.Update",
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

	roleId, meta, err := h.Command.Update(ctx, id, payload)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Application successfully updated", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "CustomerHttpHandler.Update", "channel", "Application successfully updated", "")
}

func (h ApplicationHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ApplicationHttpHandler.Delete",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	roleId, meta, err := h.Command.Delete(ctx, id)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Application successfully deleted", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "CustomerHttpHandler.Delete", "channel", "Application successfully deleted", "")
}

func (h ApplicationHttpHandler) Find(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	limit := rhelper.QueryStringToInt64(r, "limit", 0)
	if limit < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid limit",
			"ApplicationHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	skip := rhelper.QueryStringToInt64(r, "skip", 0)
	if skip < 0 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid skip",
			"ApplicationHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	filter := repository.ApplicationFilter{
		Sort: rhelper.QueryString(r, "sort"),
	}

	roles, meta, err := h.Query.Find(ctx, filter, limit, skip)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Application successfully retrieved", roles, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "CustomerHttpHandler.Find", "channel", "Application successfully retrieved", "")
}

func (h ApplicationHttpHandler) FindById(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ApplicationHttpHandler.FindById",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	role, meta, err := h.Query.FindById(ctx, id)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Application successfully retrieved", role, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "CustomerHttpHandler.FindById", "channel", "Application successfully retrieved", "")
}
