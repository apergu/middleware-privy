package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func CreateResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := map[string]interface{}{
		"status":  http.StatusText(statusCode),
		"code":    statusCode,
		"message": message,
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

type JSONResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func HandleResponse(e echo.Context, statusCode int, message string, data interface{}) error {
	log.Println("-----------------------------------------------")
	log.Println("status code : ", statusCode)
	log.Println("message : ", message)
	log.Println("data : ", data)

	dataResponse := &JSONResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	return e.JSON(statusCode, dataResponse)
}

func WriteJSONResponse(w http.ResponseWriter, response map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle encoding error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GenerateJSONResponse(code int, status bool, message string, data interface{}) (map[string]interface{}, error) {
	response := map[string]interface{}{
		"code":    code,
		"success": status,
		"message": message,
		"data":    data,
	}

	return response, nil
}

func GetErrorStatusCode(err error) int {
	newErr := strings.Split(err.Error(), " ")

	switch newErr[0] {
	case "":
		return http.StatusOK
	case "[err_duplicate]":
		return http.StatusConflict
	case "[err_invalid_payload]":
		return http.StatusUnprocessableEntity
	case "[err_bad_request]":
		return http.StatusBadRequest
	case "[err_internal_server]":
		return http.StatusInternalServerError
	case "[err_not_found]":
		return http.StatusNotFound
	case "[err_unauthorized]":
		return http.StatusUnauthorized
	case "[err_forbidden]":
		return http.StatusForbidden
	case "[err_unprocessable_entity]":
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
