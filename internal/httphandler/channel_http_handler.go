package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rhelper"
	"gitlab.com/rteja-library3/rresponser"

	"middleware/internal/helper"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/internal/usecase"
)

type ChannelHttpHandler struct {
	Command  usecase.ChannelCommandUsecase
	Query    usecase.ChannelQueryUsecase
	Decorder rdecoder.Decoder
}

func NewChannelHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.ChannelUsecaseProperty{
		ChannelRepo:  repository.NewChannelRepositoryPostgre(prop.DBPool),
		ChannelPrivy: prop.DefaultCredential,
		MerchantRepo: repository.NewMerchantRepositoryPostgre(prop.DBPool),
	}

	handler := ChannelHttpHandler{
		Command:  usecase.NewChannelCommandUsecaseGeneral(ucProp),
		Query:    usecase.NewChannelQueryUsecaseGeneral(ucProp),
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

func (h ChannelHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Channel

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)

	if err != nil {

		msg := err.Error()
		re := regexp.MustCompile(`Channel\.(\w+)`)
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
	// 			"at":     "ChannelHttpHandler.Create",
	// 			"src":    "rdecoder.DecodeRest",
	// 		}).
	// 		Error(err)

	// 	err = rapperror.ErrBadRequest(
	// 		rapperror.AppErrorCodeBadRequest,
	// 		"Invalid body",
	// 		"ChannelHttpHandler.Create",
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

		helper.LoggerValidateStructfunc(w, r, "ChannelHttpHandler.Create", "channel", message, "", payload)

		logrus.
			WithFields(logrus.Fields{
				"at":     "ChannelHttpHandler.Create",
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
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
		return
	}
	response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Channel successfully created", map[string]interface{}{
		"roleId": roleId,
		"meta":   meta,
	})

	helper.WriteJSONResponse(w, response, http.StatusCreated)
	helper.LoggerSuccessStructfunc(w, r, "ChannelHttpHandler.Create", "channel", "Channel successfully created", "")
}

func (h ChannelHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Channel

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ChannelHttpHandler.Update",
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
				"at":     "ChannelHttpHandler.Update",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ChannelHttpHandler.Update",
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

		helper.LoggerValidateStructfunc(w, r, "ChannelHttpHandler.Update", "channel", message, "", payload)

		logrus.
			WithFields(logrus.Fields{
				"at":     "ChannelHttpHandler.Update",
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

	response = rresponser.NewResponserSuccessOK("", "Channel successfully updated", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "ChannelHttpHandler.Update", "channel", "Channel successfully updated", "")
}

func (h ChannelHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ChannelHttpHandler.Delete",
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

	response = rresponser.NewResponserSuccessOK("", "Channel successfully deleted", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "ChannelHttpHandler.Delete", "channel", "Channel successfully deleted", "")
}

func (h ChannelHttpHandler) Find(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	limit := rhelper.QueryStringToInt64(r, "limit", 0)
	if limit < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid limit",
			"ChannelHttpHandler.Find",
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
			"ChannelHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	filter := repository.ChannelFilter{
		Sort: rhelper.QueryString(r, "sort"),
	}

	roles, meta, err := h.Query.Find(ctx, filter, limit, skip)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Channel successfully retrieved", roles, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "ChannelHttpHandler.Find", "channel", "Channel successfully retrieved", "")
}

func (h ChannelHttpHandler) FindById(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ChannelHttpHandler.FindById",
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

	response = rresponser.NewResponserSuccessOK("", "Channel successfully retrieved", role, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.LoggerSuccessStructfunc(w, r, "ChannelHttpHandler.FindById", "channel", "Channel successfully retrieved", "")
}
