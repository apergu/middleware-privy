package controllers

import (
	"encoding/json"
	"gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/services"
	"gitlab.com/mohamadikbal/project-privy/system"
	"log"
	"net/http"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting process create skill")

	response := new(system.JSONResponseReturn)

	var req model.Customer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		/* Log Handle error */
		system.HandleError("Error on binding request", err.Error())

		response.StatusCode = http.StatusBadRequest
		response.Message = "Bad Request"
		response.Data = nil
		w.WriteHeader(response.StatusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	dataPost := map[string]interface{}{
		"Data": req,
	}

	status := services.CreateCustomer(dataPost)

	if status == "" {
		response.StatusCode = 400
		response.Message = "ERROR"
		response.Data = nil
		w.WriteHeader(response.StatusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.StatusCode = 200
	response.Message = "SUKSES"
	response.Data = status
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
