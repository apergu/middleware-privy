package appmiddleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"gitlab.com/rteja-library3/rcache"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
	"gitlab.com/rteja-library3/rtoken"
)

const (
	unexpectedJWT  string = "unexpected error"
	tokenExpired   string = "token expired"
	tokenInvalid   string = "token invalid"
	tokenNoContent string = "token no content"
)

// JWTAuthenticatorMiddleware definition
func JWTAuthenticatorMiddleware(tokener rtoken.TokenParseAble, decorder rdecoder.Decoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtAuthenticator(next, tokener, decorder)
	}
}

// JWTAuthenticatorMiddleware definition
func jwtAuthenticator(next http.Handler, tokener rtoken.TokenParseAble, decorder rdecoder.Decoder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			response := rresponser.NewResponserError(
				sessionCoverError(err, "JWT authenticator"),
			)

			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		if token == nil {
			response := rresponser.NewResponserError(
				sessionCoverError(rcache.ErrNotFound, "JWT authenticator"),
			)

			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		if token.Expiration().Before(time.Now()) {
			response := rresponser.NewResponserError(
				sessionCoverError(errors.New("token is expired"), "JWT authenticator"),
			)

			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		next.ServeHTTP(w, r)
	})
}
