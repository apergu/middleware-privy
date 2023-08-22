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

type TopUpDataHttpHandler struct {
	Command  usecase.TopUpDataCommandUsecase
	Query    usecase.TopUpDataQueryUsecase
	Decorder rdecoder.Decoder
}

func NewTopUpDataHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.TopUpDataUsecaseProperty{
		TopUpDataRepo:  repository.NewTopUpDataRepositoryPostgre(prop.DBPool),
		TopUpDataPrivy: prop.DefaultCredential,
	}

	handler := TopUpDataHttpHandler{
		Command:  usecase.NewTopUpDataCommandUsecaseGeneral(ucProp),
		Query:    usecase.NewTopUpDataQueryUsecaseGeneral(ucProp),
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

func (h TopUpDataHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.TopUpData

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "TopUpDataHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"TopUpDataHttpHandler.Create",
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
				"at":     "TopUpDataHttpHandler.Create",
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

	response = rresponser.NewResponserSuccessCreated("", "TopUpData successfully created", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h TopUpDataHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.TopUpData

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"TopUpDataHttpHandler.Update",
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
				"at":     "TopUpDataHttpHandler.Update",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"TopUpDataHttpHandler.Update",
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
				"at":     "TopUpDataHttpHandler.Update",
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

	response = rresponser.NewResponserSuccessOK("", "TopUpData successfully updated", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h TopUpDataHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"TopUpDataHttpHandler.Delete",
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

	response = rresponser.NewResponserSuccessOK("", "TopUpData successfully deleted", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h TopUpDataHttpHandler) Find(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	limit := rhelper.QueryStringToInt64(r, "limit", 0)
	if limit < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid limit",
			"TopUpDataHttpHandler.Find",
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
			"TopUpDataHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	filter := repository.TopUpDataFilter{
		Sort: rhelper.QueryString(r, "sort"),
	}

	roles, meta, err := h.Query.Find(ctx, filter, limit, skip)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "TopUpData successfully retrieved", roles, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h TopUpDataHttpHandler) FindById(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"TopUpDataHttpHandler.FindById",
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

	response = rresponser.NewResponserSuccessOK("", "TopUpData successfully retrieved", role, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}
