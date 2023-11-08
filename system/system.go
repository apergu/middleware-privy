package system

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
)

func MarshalUnmarshal(param interface{}, result interface{}) error {
	paramByte, err := json.Marshal(param)
	if err != nil {
		log.Println("Error marshal", err.Error())
		return err
	}

	err = json.Unmarshal(paramByte, &result)
	if err != nil {
		log.Println("Error unmarshal", err.Error())
		return err
	}

	return nil
}

func GenerateRandInt() int64 {
	y1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := y1.Int63n(10000)
	return id
}

type JSONResponse struct {
	Id            int64       `json:"Id"`
	RequstId      int64       `json:"ReqId"`
	Status        int         `json:"Status"`
	StatusMessage string      `json:"StatusMessage"`
	ErrorCode     string      `json:"ErrorCode"`
	Time          time.Time   `json:"Time"`
	Signature     string      `json:"Signature`
	Data          interface{} `json:"Data"`
}

func (response *JSONResponse) Error(id, reqId int64, status int, statusMessage, errCode string, time time.Time, data interface{}) {
	response.Id = id
	response.RequstId = reqId
	response.Status = status
	response.StatusMessage = statusMessage
	response.ErrorCode = errCode
	response.Time = time
	response.Data = data
}

func (response *JSONResponse) Success(id, reqId int64, status int, statusMessage string, time time.Time, data interface{}) {
	response.Id = id
	response.RequstId = reqId
	response.Status = status
	response.StatusMessage = statusMessage
	response.ErrorCode = ""
	response.Time = time
	response.Data = data
}

func HandleError(message string, err interface{}) {
	log.Println()
	log.Println("========== Start Error Message ==========")
	log.Println("Message => " + message + ".")

	if err != nil {
		log.Println("Error => ", err)
	}

	log.Println("========== End Of Error Message ==========")
	log.Println()
}

// JSONEncode is a function to encode data to JSON
func JSONEncode(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}

func HandleJSONResponse(id, reqId int64, status int, statusMessage, errCode string, time time.Time, data interface{}) string {
	var responseStruct = new(JSONResponse)

	if statusMessage == "Success" {
		responseStruct.Success(id, reqId, status, statusMessage, time, data)
	} else {
		HandleError(statusMessage, data)

		if data == nil {
			responseStruct.Error(id, reqId, status, statusMessage, errCode, time, nil)
		} else if fmt.Sprintf("%v", reflect.TypeOf(data).Kind()) == "ptr" {
			responseStruct.Error(id, reqId, status, statusMessage, errCode, time, fmt.Sprintf("%v", data))
		} else {
			responseStruct.Error(id, reqId, status, statusMessage, errCode, time, data)
		}
	}

	log.Println("Closing")

	return JSONEncode(responseStruct)
}

type Request struct {
	ID         int64       `json:"Id"`
	Command    string      `json:"Command"`
	Time       string      `json:"Time"`
	ModuleId   string      `json:"ModuleId"`
	Properties interface{} `json:"Properties"`
	Signature  string      `json:"Signature"`
	Data       interface{} `json:"Data"`
}

type JSONResponseReturn struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
