package appmiddleware

import (
	"context"
	"middleware/internal/constants"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

func findToken(r *http.Request, fns ...func(r *http.Request) string) string {
	for _, fn := range fns {
		tokenString := fn(r)
		if tokenString != "" {
			return tokenString
		}
	}

	return ""
}

func defaultFindToken(r *http.Request) string {
	return findToken(r, jwtauth.TokenFromCookie, jwtauth.TokenFromHeader)
}

func HandlerSetContextValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = setContextValue(r)
		next.ServeHTTP(w, r)
	})
}

func setContextValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, constants.Timestamp, time.Now())
	return r.WithContext(ctx)
}
