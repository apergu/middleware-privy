package httphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/mohamadikbal/project-privy/internal/constants"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/repository"
	"gitlab.com/mohamadikbal/project-privy/internal/usecase"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rhelper"
	"gitlab.com/rteja-library3/rresponser"
)

type MerchantHttpHandler struct {
	Command  usecase.MerchantCommandUsecase
	Query    usecase.MerchantQueryUsecase
	Decorder rdecoder.Decoder
}

func NewMerchantHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.MerchantUsecaseProperty{
		MerchantRepo:  repository.NewMerchantRepositoryPostgre(prop.DBPool),
		MerchantPrivy: prop.DefaultCredential,
	}

	handler := MerchantHttpHandler{
		Command:  usecase.NewMerchantCommandUsecaseGeneral(ucProp),
		Query:    usecase.NewMerchantQueryUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/", handler.Create)
	// r.Put("/id/{id}", handler.Update)
	// r.Delete("/id/{id}", handler.Delete)

	r.Get("/", handler.Find)
	r.Get("/id/{id}", handler.FindById)

	return r
}

func (h MerchantHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Merchant

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "MerchantHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"MerchantHttpHandler.Create",
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

	err = payload.Validate()
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":     "MerchantHttpHandler.Create",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	roleId, meta, err := h.Command.Create(ctx, payload)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessCreated("", "Merchant successfully created", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h MerchantHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Merchant

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"MerchantHttpHandler.Update",
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
				"at":     "MerchantHttpHandler.Update",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"MerchantHttpHandler.Update",
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

	err = payload.Validate()
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":     "MerchantHttpHandler.Update",
				"src":    "payload.Validate",
				"params": payload,
			}).
			Error(err)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	roleId, meta, err := h.Command.Update(ctx, id, payload)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Merchant successfully updated", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h MerchantHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"MerchantHttpHandler.Delete",
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

	response = rresponser.NewResponserSuccessOK("", "Merchant successfully deleted", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h MerchantHttpHandler) Find(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	limit := rhelper.QueryStringToInt64(r, "limit", 0)
	if limit < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid limit",
			"MerchantHttpHandler.Find",
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
			"MerchantHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	filter := repository.MerchantFilter{
		Sort: rhelper.QueryString(r, "sort"),
	}

	roles, meta, err := h.Query.Find(ctx, filter, limit, skip)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Merchant successfully retrieved", roles, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h MerchantHttpHandler) FindById(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"MerchantHttpHandler.FindById",
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

	response = rresponser.NewResponserSuccessOK("", "Merchant successfully retrieved", role, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}
