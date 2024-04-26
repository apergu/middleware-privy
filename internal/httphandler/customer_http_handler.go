package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"middleware/internal/constants"
	"middleware/internal/helper"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rhelper"
	"gitlab.com/rteja-library3/rresponser"
)

type CustomerHttpHandler struct {
	Command  usecase.CustomerCommandUsecase
	Query    usecase.CustomerQueryUsecase
	Decorder rdecoder.Decoder
}

func NewCustomerHttpHandler(prop HTTPHandlerProperty) http.Handler {
	ucProp := usecase.CustomerUsecaseProperty{
		CustomerRepo:  repository.NewCustomerRepositoryPostgre(prop.DBPool),
		CustomerPrivy: prop.DefaultCredential,
	}

	handler := CustomerHttpHandler{
		Command:  usecase.NewCustomerCommandUsecaseGeneral(ucProp),
		Query:    usecase.NewCustomerQueryUsecaseGeneral(ucProp),
		Decorder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Post("/", handler.Create)
	r.Post("/lead", handler.CreateLead)
	r.Put("/lead/{id}", handler.UpdateLead)
	r.Put("/id/{id}", handler.Update)
	r.Delete("/id/{id}", handler.Delete)

	r.Get("/", handler.Find)
	r.Get("/id/{id}", handler.FindById)

	return r
}

func (h CustomerHttpHandler) Create(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Customer
	//var payloadLead model.Lead

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "CustomerHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"CustomerHttpHandler.Create",
			nil,
		)

		// response = rresponser.NewResponserError(err)
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	// get user from context
	user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = user

	errors := payload.Validate()
	if payload.EntityStatus == "6" || payload.EntityStatus == "" {
		if payload.SubIndustry == "" {
			errors = append(errors, map[string]interface{}{
				"field":   "SubIndustry",
				"message": "Sub Industry is required",
			})
		}
	}

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

	// if payload.CRMLeadID == "" {

	// } else {
	// if payload.CRMLeadID != "" {
	if payload.EntityStatus == "13" {
		log.Println("payload masuk 13", payload)
		roleId, meta, err := h.Command.Create(ctx, payload)
		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, "Create Customer Failed", map[string]interface{}{
				"roleId": roleId,
				"meta":   meta,
			})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		response, _ := helper.GenerateJSONResponse(http.StatusCreated, true, "Customer successfully created", map[string]interface{}{
			"roleId": roleId,
			"meta":   meta,
		})

		// response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", roleId, meta)
		helper.WriteJSONResponse(w, response, http.StatusCreated)
	}

	if payload.EntityStatus == "6" || payload.EntityStatus == "" {
		log.Println("payload masuk 6", payload)

		roleId, meta, err := h.Command.CreateLeadZD(ctx, payload)
		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		url := os.Getenv("ACZD_BASE") + "api/v1/privy/zendesk/lead"

		fmt.Println("url", url)

		// Replace the following map with your actual data
		data := map[string]interface{}{
			"zd_lead_id":          payload.CRMLeadID,
			"first_name":          payload.FirstName,
			"last_name":           payload.LastName,
			"enterprise_privy_id": payload.EnterprisePrivyID,
			"enterprise_name":     payload.CustomerName,
			"address":             payload.Address,
			"email":               payload.Email,
			"zip":                 payload.ZipCode,
			"state":               payload.State,
			"country":             "Indonesia",
			"city":                payload.City,
			"npwp":                payload.NPWP,
			"sub_industry":        payload.SubIndustry,
			"phone":               payload.PhoneNo,
		}

		// Convert data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		// Make the HTTP POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD"))

		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println("response", err)

		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		// if resp.StatusCode != 200 {
		// 	fmt.Println("response status code", resp.StatusCode)
		// 	response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
		// 	// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		// 	helper.WriteJSONResponse(w, response, http.StatusBadRequest)
		// 	return
		// }
		defer resp.Body.Close()

		response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{
			"roleId": roleId,
			"meta":   meta,
		})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusCreated)

		// response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", resp.StatusCode, resp.Body)
	}

	if payload.EntityStatus == "7" {
		roleId, meta, err := h.Command.CreateLead2(ctx, payload)
		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{
			"roleId": roleId,
			"meta":   meta,
		})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusCreated)
		// response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", roleId, meta)
	}
	// } else {
	// 	log.Println("CRM LEAD ID KOSONG", payload)
	// 	roleId, meta, err := h.Command.CreateLead2(ctx, payload)
	// 	if err != nil {
	// 		response = rresponser.NewResponserError(err)
	// 		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	// 		return
	// 	}
	// 	response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", roleId, meta)
	// }

	// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h CustomerHttpHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Lead

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "CustomerHttpHandler.Create",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		// err = rapperror.ErrBadRequest(
		// 	rapperror.AppErrorCodeBadRequest,
		// 	"Invalid body",
		// 	"CustomerHttpHandler.Create",
		// 	nil,
		// )
		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

		return
	}

	// get user from context
	user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = user

	errors := payload.ValidateLead()
	if len(errors) > 0 {
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

	roleId, meta, err := h.Command.CreateLead(ctx, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{
		"roleId": roleId,
		"meta":   meta,
	})
	// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	helper.WriteJSONResponse(w, response, http.StatusCreated)

}

func (h CustomerHttpHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Lead

	id := chi.URLParam(r, "id")
	//id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	//if id < 1 {
	//	err = rapperror.ErrBadRequest(
	//		rapperror.AppErrorCodeBadRequest,
	//		"Invalid id",
	//		"CustomerHttpHandler.Update",
	//		nil,
	//	)
	//
	//	response = rresponser.NewResponserError(err)
	//	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	//	return
	//}

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "CustomerHttpHandler.Update",
				"src":    "rdecoder.DecodeRest",
			}).
			Error(err)

		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}
	}

	// get user from context
	user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = user

	errors := payload.ValidateLead()
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

	roleId, meta, err := h.Command.UpdateLead(ctx, id, payload)

	if err != nil {
		response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{
		"roleId": roleId,
		"meta":   meta,
	})

	helper.WriteJSONResponse(w, response, http.StatusCreated)

}

func (h CustomerHttpHandler) Update(w http.ResponseWriter, r *http.Request) {
	// var response rresponser.Responser
	var err error
	ctx := r.Context()

	var payload model.Customer

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"CustomerHttpHandler.Update",
			nil,
		)

		if err != nil {
			response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, http.StatusBadRequest)
			return
		}

	}

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)

	if err != nil {
		response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	// if err != nil {
	// 	logrus.
	// 		WithFields(logrus.Fields{
	// 			"action": "try to decode data",
	// 			"at":     "CustomerHttpHandler.Update",
	// 			"src":    "rdecoder.DecodeRest",
	// 		}).
	// 		Error(err)

	// 	err = rapperror.ErrBadRequest(
	// 		rapperror.AppErrorCodeBadRequest,
	// 		"Invalid body",
	// 		"CustomerHttpHandler.Update",
	// 		nil,
	// 	)

	// 	response = rresponser.NewResponserError(err)
	// 	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
	// 	return
	// }

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

	roleId, meta, err := h.Command.Update(ctx, id, payload)
	if err != nil {
		response, _ := helper.GenerateJSONResponse(http.StatusBadRequest, false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusBadRequest)
		return
	}

	response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{
		"roleId": roleId,
		"meta":   meta,
	})

	helper.WriteJSONResponse(w, response, http.StatusCreated)

}

func (h CustomerHttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"CustomerHttpHandler.Delete",
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

	response = rresponser.NewResponserSuccessOK("", "Customer successfully deleted", roleId, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h CustomerHttpHandler) Find(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	limit := rhelper.QueryStringToInt64(r, "limit", 0)
	if limit < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid limit",
			"CustomerHttpHandler.Find",
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
			"CustomerHttpHandler.Find",
			nil,
		)

		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	filter := repository.CustomerFilter{
		Sort: rhelper.QueryString(r, "sort"),
	}

	roles, meta, err := h.Query.Find(ctx, filter, limit, skip)
	if err != nil {
		response = rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		return
	}

	response = rresponser.NewResponserSuccessOK("", "Customer successfully retrieved", roles, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}

func (h CustomerHttpHandler) FindById(w http.ResponseWriter, r *http.Request) {
	var response rresponser.Responser
	var err error
	ctx := r.Context()

	id := rhelper.ToInt64(chi.URLParam(r, "id"), 0)
	if id < 1 {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"CustomerHttpHandler.FindById",
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

	response = rresponser.NewResponserSuccessOK("", "Customer successfully retrieved", role, meta)
	rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
}
