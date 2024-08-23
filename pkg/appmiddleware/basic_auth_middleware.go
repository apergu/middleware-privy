package appmiddleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"

	"middleware/internal/constants"
	"middleware/internal/helper"
)

func BasicAuth(basicUsername, basicPassword string, decorder rdecoder.Decoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return basicAuthHandler(next, basicUsername, basicPassword, decorder)
	}
}

func basicAuthHandler(next http.Handler, basicUsername, basicPassword string, decorder rdecoder.Decoder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		username, password, ok := r.BasicAuth()

		if !ok {
			err := rapperror.ErrUnauthorized(
				"",
				"Authorization is required",
				"",
				nil,
			)

			logrus.
				WithFields(logrus.Fields{
					"at":  "basicAuthHandler",
					"src": "r.BasicAuth()",
				}).
				Error(err)

			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

		usernameMatch := username == basicUsername
		passwordMatch := password == basicPassword

		if !usernameMatch || !passwordMatch {
			err := rapperror.ErrUnauthorized(
				"",
				"Invalid username or password",
				"",
				nil,
			)

			logrus.
				WithFields(logrus.Fields{
					"at":  "basicAuthHandler",
					"src": "!usernameMatch || !passwordMatch",
				}).
				Error(err)

			response, _ := helper.GenerateJSONResponse(helper.GetErrorStatusCode(err), false, err.Error(), map[string]interface{}{})
			helper.WriteJSONResponse(w, response, helper.GetErrorStatusCode(err))
			return
		}

		ctxWthValue := context.WithValue(ctx, constants.SessionUserId, int64(0))

		next.ServeHTTP(w, r.WithContext(ctxWthValue))
	})
}
