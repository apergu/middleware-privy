package exceptions

import (
	"errors"
	"net/http"
)

var (
	ErrPayload             = errors.New("invalid request payload")
	ErrInternalServerError = errors.New("error - 500 internal server error - an error occurred in the system. please try again later")
	ErrBadRequest          = errors.New("error - 400 bad Request - please try again later")
	ErrUnAuthorized        = errors.New("401 unauthorized")
	ErrNotFound            = errors.New("not found")
	ErrMappingNotfound     = errors.New("mapping not found")
	ErrCb                  = errors.New("service is on maintenance, please try again later")
)

func MappingStatusCode(err error) int {
	var statusCode int
	switch err {
	case ErrNotFound:
		statusCode = http.StatusNotFound
	case ErrUnAuthorized:
		statusCode = http.StatusUnauthorized
	case ErrMappingNotfound:
		statusCode = http.StatusBadRequest
	case ErrBadRequest:
		statusCode = http.StatusBadRequest
	case ErrPayload:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	return statusCode
}
