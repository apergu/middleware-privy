package services

import (
	model2 "gitlab.com/mohamadikbal/project-privy/internal/model"
	"gitlab.com/mohamadikbal/project-privy/internal/models"
	"gitlab.com/mohamadikbal/project-privy/system"
	"net/http"
	"time"
)

func CreateCustomer(data interface{}) string {
	id := system.GenerateRandInt()

	/* Bind Request data to Struct */
	var result model2.Customer
	var newData model2.Customer
	modelx := models.ModelCustomer{}
	system.MarshalUnmarshal(data, &newData)
	modelx.Customer = newData
	result, err := modelx.CreateCustomer()

	if err != nil {
		return system.HandleJSONResponse(id, 1, http.StatusBadRequest, "Failed", "", time.Now(), nil)
	}

	return system.HandleJSONResponse(id, 1, http.StatusOK, "Success", "", time.Now(), result)
}
