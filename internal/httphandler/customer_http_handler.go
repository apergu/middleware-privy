package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rhelper"
	"gitlab.com/rteja-library3/rresponser"

	"middleware/internal/constants"
	"middleware/internal/helper"
	"middleware/internal/model"
	"middleware/internal/repository"
	"middleware/internal/usecase"
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
	fmt.Println("err =>", err)
	if err != nil {
		msg := err.Error()
		re := regexp.MustCompile(`Customer\.(\w+)`)
		custm := re.FindStringSubmatch(msg)
		re = regexp.MustCompile(`([a-z])([A-Z])`)
		spaced := re.ReplaceAllString(custm[1], `$1 $2`)
		re = regexp.MustCompile(`type ([^\]]+)`)
		format := re.FindStringSubmatch(msg)
		message := fmt.Sprintf("Unprocessable entity - %s value must in %s format", spaced, format[1])
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "CustomerHttpHandler.Create",
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

	// get user from context
	user := ctx.Value(constants.SessionUserId).(int64)

	// set created by value
	payload.CreatedBy = user

	// fmt.Println("BEFORE ERRROR ")

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
			// fmt.Println("error", err.Error())
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{
				"roleId": roleId,
				"meta":   meta,
			})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

		if payload.CRMLeadID != "" {
			// GET DETAIL DATA

			urlDetailData := "https://api.getbase.com/v2/leads/"
			reqDetailData, err := http.NewRequest("GET", urlDetailData+payload.CRMLeadID, nil)

			reqDetailData.Header.Add("Content-Type", "application/json")
			reqDetailData.Header.Add("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")

			clientDetailData := &http.Client{}
			respDetailData, err := clientDetailData.Do(reqDetailData)
			fmt.Println("response", respDetailData.Body)

			defer respDetailData.Body.Close()

			bodyDetailData, err := ioutil.ReadAll(respDetailData.Body)

			var responsDetailData struct {
				Data interface{} `json:"data"`
			}

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			err = json.Unmarshal(bodyDetailData, &responsDetailData)
			fmt.Println("response Body", responsDetailData)

			var responseDetail struct {
				FirstName    string                 `json:"first_name"`
				LastName     string                 `json:"last_name"`
				Email        string                 `json:"email"`
				PhoneNumber  string                 `json:"phone"`
				CustomFields map[string]interface{} `json:"custom_fields"`
			}

			customFieldsData := responsDetailData.Data.(map[string]interface{})["custom_fields"].(map[string]interface{})

			newResp := responsDetailData.Data.(map[string]interface{})
			responseDetail.FirstName = newResp["first_name"].(string)
			responseDetail.LastName = newResp["last_name"].(string)
			responseDetail.Email = newResp["email"].(string)
			if newResp["phone"] != nil {
				responseDetail.PhoneNumber = responsDetailData.Data.(map[string]interface{})["phone"].(string)
			}

			if newResp["mobile"] != nil {
				responseDetail.PhoneNumber = responsDetailData.Data.(map[string]interface{})["mobile"].(string)
			}
			// err = json.Unmarshal([]byte(newResp), &responseDetail)
			fmt.Println("response Body Detail", responseDetail)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// END DETAIL DATA

			// START GET DATA =================

			urlGetData := "https://api.getbase.com/v2/leads/"

			// Make the HTTP POST request
			payloadEdit := map[string]interface{}{
				"organization_name": payload.CustomerName,
				"address": map[string]interface{}{
					"line1":       payload.Address,
					"city":        payload.City,
					"postal_code": payload.ZipCode,
					"state":       payload.State,
					"country":     "ID",
				},
				"custom_fields": map[string]interface{}{
					"Sub Industry":  payload.SubIndustry,
					"NPWP":          payload.NPWP,
					"Enterprise ID": payload.EnterprisePrivyID,
				},
			}

			if responseDetail.FirstName != payload.FirstName {
				payloadEdit["custom_fields"].(map[string]interface{})["First Name - Adonara"] = payload.FirstName
			}

			if responseDetail.LastName != payload.LastName {
				payloadEdit["custom_fields"].(map[string]interface{})["Last Name - Adonara"] = payload.LastName
			}

			if responseDetail.Email != payload.Email {
				payloadEdit["custom_fields"].(map[string]interface{})["Email - Adonara"] = payload.Email
			}

			if responseDetail.PhoneNumber != payload.PhoneNo {
				payloadEdit["custom_fields"].(map[string]interface{})["Phone Number - Adonara"] = payload.PhoneNo
			}

			dataPayloadEdit := map[string]interface{}{
				"data": payloadEdit,
			}

			// Convert data to JSON
			jsonDataEdit, err := json.Marshal(dataPayloadEdit)

			reqGetData, err := http.NewRequest("PUT", urlGetData+payload.CRMLeadID, bytes.NewBuffer(jsonDataEdit))

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			}

			reqGetData.Header.Add("Content-Type", "application/json")
			reqGetData.Header.Add("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")

			clientGetData := &http.Client{}
			respGetData, err := clientGetData.Do(reqGetData)
			fmt.Println("response", respGetData.Body)

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			// fmt.Println("respGetData", customFieldsData["ActiveCampaign Contact ID"].(string))

			defer respGetData.Body.Close()

			// END GET DATA =================

			// LEADS CONVERSION =============
			url := "https://api.getbase.com/v2/lead_conversions"

			// fmt.Println("url", url)

			// // Replace the following map with your actual data
			leadID, _ := strconv.Atoi(payload.CRMLeadID)
			data := map[string]interface{}{
				"lead_id": leadID,
			}

			payloadData := map[string]interface{}{
				"data": data,
			}

			// // Convert data to JSON
			jsonData, err := json.Marshal(payloadData)
			if err != nil {
				panic(err)
			}

			// // Make the HTTP POST request
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")
			// req.SetBasicAuth(os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD"))

			client := &http.Client{}
			resp, err := client.Do(req)
			fmt.Println("response", err)

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			defer resp.Body.Close()

			type Data struct {
				DealId int `json:"deal_id"`
			}

			body, err := ioutil.ReadAll(resp.Body)
			// jsonBody := json.RawMessage(body)
			fmt.Println("TEST", string(body))
			if err != nil {
				fmt.Println("Error reading body:", err)
				return
			}

			var respons struct {
				Data Data `json:"data"`
			}

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			err = json.Unmarshal(body, &respons)
			fmt.Println("response Body", respons)

			// ACTIVE CAMPAIGN UPDATE

			urlAC := "https://privy1706071639.api-us1.com/api/3/contacts/" + customFieldsData["ActiveCampaign Contact ID"].(string)

			fmt.Println("url", urlAC)

			payloadAc := map[string]interface{}{
				"contact": map[string]interface{}{
					"lastName": payload.LastName,
					"email":    responseDetail.Email,
					"phone":    responseDetail.PhoneNumber,
					"fieldValues": []map[string]interface{}{
						{
							"field": 1,
							"value": payload.CustomerName,
						}, {
							"field": 2,
							"value": payload.SubIndustry,
						}, {
							"field": 3,
							"value": "New Client - Inbound",
						},
						{
							"field": 4,
							"value": "Won - Contract Signed / Award Letter Issued",
						}, {
							"field": 5,
							"value": payload.EnterprisePrivyID,
						}, {
							"field": 7,
							"value": strconv.Itoa(respons.Data.DealId),
						},
					},
				},
			}

			// Convert data to JSON
			jsonDataAc, err := json.Marshal(payloadAc)
			if err != nil {
				panic(err)
			}

			// Make the HTTP POST request
			reqAc, err := http.NewRequest("PUT", urlAC, bytes.NewBuffer(jsonDataAc))

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			reqAc.Header.Add("Content-Type", "application/json")
			reqAc.Header.Add("Api-Token", "83098f1b9181f163ee582823ba5bdcde7a02db14d75b8fc3dc2eea91738a49a47e100e68")

			clientAc := &http.Client{}
			respAc, err := clientAc.Do(reqAc)
			fmt.Println("response", err)

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			defer respAc.Body.Close()

			// DEALS UPDATE WON ============
			urlWon := "https://api.getbase.com/v2/deals/"
			fmt.Println("urlWon", urlWon+strconv.Itoa(respons.Data.DealId))

			// fmt.Println("url", url)

			// // Replace the following map with your actual data
			dataWon := map[string]interface{}{
				"stage_id": 34108700,
			}

			payloadDataWon := map[string]interface{}{
				"data": dataWon,
			}

			// // Convert data to JSON
			jsonDataWon, err := json.Marshal(payloadDataWon)
			if err != nil {
				panic(err)
			}

			// Make the HTTP POST request
			reqWon, err := http.NewRequest("PUT", urlWon+strconv.Itoa(respons.Data.DealId), bytes.NewBuffer(jsonDataWon))

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			reqWon.Header.Add("Content-Type", "application/json")
			reqWon.Header.Add("Authorization", "Bearer 26bed09778079a78eb96acb73feb1cb2d9b36267e992caa12b0d960c8f760e2c")
			// req.SetBasicAuth(os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD"))

			clientWon := &http.Client{}
			respWon, err := clientWon.Do(reqWon)
			fmt.Println("response", err)

			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}

			defer respWon.Body.Close()

			bodyWon, err := ioutil.ReadAll(respWon.Body)
			jsonBodyWon := json.RawMessage(bodyWon)
			response, _ := helper.GenerateJSONResponse(http.StatusCreated, true, "Customer successfully created", map[string]interface{}{
				"roleId": 1,
				"meta":   jsonBodyWon,
			})
			helper.WriteJSONResponse(w, response, http.StatusCreated)
		}

		response, _ := helper.GenerateJSONResponse(http.StatusCreated, true, "Customer successfully created", map[string]interface{}{
			"roleId": 1,
			"meta":   nil,
		})
		helper.WriteJSONResponse(w, response, http.StatusCreated)

		// // response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", roleId, meta)
	}

	if payload.EntityStatus == "6" || payload.EntityStatus == "" {
		log.Println("payload masuk 6", payload)

		if payload.CRMLeadID == "" {
			_, _, err := h.Command.CreateLeadZD(ctx, payload)
			if err != nil {
				response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
				// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
				helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
				return
			}
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
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.SetBasicAuth(os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD"))

		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println("response", err)

		if err != nil {
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

		defer resp.Body.Close()

		response, _ := helper.GenerateJSONResponse(http.StatusCreated, false, "Customer successfully created", map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, http.StatusCreated)

		// response = rresponser.NewResponserSuccessCreated("", "Customer successfully created", resp.StatusCode, resp.Body)
	}

	if payload.EntityStatus == "7" {
		roleId, meta, err := h.Command.CreateLead2(ctx, payload)
		if err != nil {
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

	}

	err = rdecoder.DecodeRest(r, h.Decorder, &payload)

	if err != nil {
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
		response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
		// rdecoder.EncodeRestWithResponser(w, h.Decorder, response)
		helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
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
