package appmiddleware

import (
	"context"
	"encoding/json"
	"net/http"

	"middleware/internal/constants"
	"middleware/internal/model"

	"gitlab.com/rteja-library3/rcache"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
	"gitlab.com/rteja-library3/rtoken"
)

func Session(token rtoken.TokenParseAble, cache rcache.CacheGetter, decorder rdecoder.Decoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sessionHandler(next, token, cache, decorder)
	}
}

func sessionHandler(next http.Handler, tokenParse rtoken.TokenParseAble, cache rcache.CacheGetter, decorder rdecoder.Decoder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := defaultFindToken(r)

		data, err := tokenParse.Parse(ctx, token)
		if err != nil {
			response := rresponser.NewResponserError(
				sessionCoverError(err, "session-handler"),
			)

			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		sessionId := data[string(constants.SessionSessionId)].(string)

		user, err := cache.Get(ctx, sessionId+"-user")
		if err != nil {
			response := rresponser.NewResponserError(
				sessionCoverError(err, "session-handler"),
			)

			rdecoder.EncodeRestWithResponser(w, decorder, response)
			return
		}

		loggedUser := model.SessionUser{}
		_ = json.Unmarshal(user, &loggedUser)

		ctxWthValue := context.WithValue(ctx, constants.SessionToken, token)
		ctxWthValue = context.WithValue(ctxWthValue, constants.SessionSessionId, sessionId)
		ctxWthValue = context.WithValue(ctxWthValue, constants.SessionUser, loggedUser)
		ctxWthValue = context.WithValue(ctxWthValue, constants.SessionUserId, loggedUser.ID)
		ctxWthValue = context.WithValue(ctxWthValue, constants.SessionUserMenu, loggedUser.Menus)

		next.ServeHTTP(w, r.WithContext(ctxWthValue))
	})
}
