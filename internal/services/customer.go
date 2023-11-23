package services

// import (
// 	model2 "middleware/internal/model"
// 	"middleware/internal/models"
// 	"middleware/system"
// 	"net/http"
// 	"time"
// )

// func CreateCustomer(data interface{}) string {
// 	id := system.GenerateRandInt()

// 	/* Bind Request data to Struct */
// 	var result model2.Customer
// 	var newData model2.Customer
// 	modelx := models.ModelCustomer{}
// 	system.MarshalUnmarshal(data, &newData)
// 	modelx.Customer = newData
// 	result, err := modelx.CreateCustomer()

// 	if err != nil {
// 		return system.HandleJSONResponse(id, 1, http.StatusBadRequest, "Failed", "", time.Now(), nil)
// 	}

// 	return system.HandleJSONResponse(id, 1, http.StatusOK, "Success", "", time.Now(), result)
// }
