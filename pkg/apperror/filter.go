package apperror

import (
	"net/http"

	"gitlab.com/rteja-library3/rapperror"
)

func IsConflictError(err error) (b bool) {
	if e, ok := err.(*rapperror.AppError); ok {
		return e.Status == http.StatusConflict
	}

	return
}
