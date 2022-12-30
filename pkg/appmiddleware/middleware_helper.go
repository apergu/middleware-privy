package appmiddleware

import (
	"net/http"

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
