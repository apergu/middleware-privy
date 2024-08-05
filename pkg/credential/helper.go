package credential

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func HttpRequest(method, url string, data []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
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
