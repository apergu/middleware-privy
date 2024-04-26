package credential

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

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
