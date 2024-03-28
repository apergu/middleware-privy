package appmiddleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"

	"middleware/internal/constants"
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
				"Invalid basic auth",
				"",
				nil,
			)

			logrus.
				WithFields(logrus.Fields{
					"at":  "basicAuthHandler",
					"src": "r.BasicAuth()",
				}).
				Error(err)

			response := rresponser.NewResponserError(err)
			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		fmt.Println(username)
		fmt.Println(password)
		fmt.Println(basicUsername)
		fmt.Println(basicPassword)
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

			response := rresponser.NewResponserError(err)
			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		ctxWthValue := context.WithValue(ctx, constants.SessionUserId, int64(0))

		next.ServeHTTP(w, r.WithContext(ctxWthValue))
	})
}
