package appmiddleware

import (
	"context"
	"net/http"

	"middleware/internal/constants"
	"middleware/internal/model"

	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
)

type RoleAccess func(access int8) bool

func RoleMenuFilter(menu string, decorder rdecoder.Decoder) func(http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			menus := ctx.Value(constants.SessionUserMenu).(model.SessionMenu)

			access := menus.FindAccess(menu)
			if access <= 0 {
				// return error
				err := rapperror.ErrForbidden(
					"",
					"User don't have access",
					"role-menu-filter-middleware",
					nil,
				)

				response := rresponser.NewResponserError(err)
				rdecoder.EncodeRestWithResponser(w, decorder, response)
				return
			}

			ctx = context.WithValue(ctx, constants.SessionUserMenuAccess, access)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return fn
}

func RoleAccessFilter(menu string, fnAccess RoleAccess, decorder rdecoder.Decoder) func(http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			menus := ctx.Value(constants.SessionUserMenu).(model.SessionMenu)

			access := menus.FindAccess(menu)
			if access <= 0 {
				// return error
				err := rapperror.ErrForbidden(
					"",
					"User don't have access",
					"role-menu-filter-middleware",
					nil,
				)

				response := rresponser.NewResponserError(err)
				rdecoder.EncodeRestWithResponser(w, decorder, response)
				return
			}

			if !fnAccess(access) {
				// return error
				err := rapperror.ErrForbidden(
					"",
					"User don't have access",
					"role-menu-filter-middleware",
					nil,
				)

				response := rresponser.NewResponserError(err)
				rdecoder.EncodeRestWithResponser(w, decorder, response)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return fn
}
