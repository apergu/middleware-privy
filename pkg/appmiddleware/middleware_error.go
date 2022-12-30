package appmiddleware

import (
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rcache"
)

func sessionCoverError(err error, src string) *rapperror.AppError {
	if err.Error() == "token is expired" {
		return rapperror.ErrUnauthorized(
			"",
			tokenExpired,
			src,
			nil,
		)
	}

	if err == rcache.ErrNotFound {
		return rapperror.ErrUnauthorized(
			"",
			tokenInvalid,
			src,
			nil,
		)
	}

	if err == rcache.ErrNoContent {
		return rapperror.ErrUnauthorized(
			"",
			tokenNoContent,
			src,
			nil,
		)
	}

	return rapperror.ErrInternalServerError(
		"",
		err.Error(),
		src,
		nil,
	)
}
