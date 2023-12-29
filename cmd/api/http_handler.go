package api

import (
	"fmt"
	"net/http"
	"strings"

	"middleware/internal/config"
	"middleware/internal/httphandler"
	"middleware/pkg/appmiddleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
	"gitlab.com/rteja-library3/rresponser"
)

func InitHttpHandler(pool *pgxpool.Pool, corsOpt cors.Options, prop httphandler.HTTPHandlerProperty, jwtAuth *jwtauth.JWTAuth, basicAuth config.BasicAuth) http.Handler {
	logrus.Info("[HttpHandler] Start")
	defer logrus.Info("[HttpHandler] End")
	defer logrus.Info("[HttpHandler] Ready")

	r := chi.NewRouter()

	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger:  logrus.StandardLogger(),
			NoColor: false,
		},
	)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(corsOpt).Handler)
	r.Use(middleware.Compress(5, "gzip"))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := rapperror.ErrNotFound(
			"",
			fmt.Sprintf("Route %s not found ", r.RequestURI),
			"",
			nil,
		)

		response := rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, rdecoder.NewJSONEncoder(), response)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := rapperror.ErrBadRequest(
			"",
			fmt.Sprintf("Route %s is not allowed for %s ", r.RequestURI, r.Method),
			"",
			nil,
		)

		response := rresponser.NewResponserError(err)
		rdecoder.EncodeRestWithResponser(w, rdecoder.NewJSONEncoder(), response)
	})

	r.Route("/api", func(r chi.Router) {
		// Health Checking...
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("SUCCESS"))
		})

		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// r.Use(jwtauth.Verifier(jwtAuth))
				// r.Use(appmiddleware.JWTAuthenticatorMiddleware(prop.DefaultToken, prop.DefaultDecoder))
				// r.Use(appmiddleware.Session(prop.DefaultToken, prop.DefaultCache, prop.DefaultDecoder))

				r.Use(appmiddleware.BasicAuth(basicAuth.Username, basicAuth.Password, prop.DefaultDecoder))

				r.Get("/healthcheck/logged", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("SUCCESS Logged"))
				})

				r.Mount("/customer", httphandler.NewCustomerHttpHandler(prop))
				r.Mount("/customer-usage", httphandler.NewCustomerUsageHttpHandler(prop))
				r.Mount("/transfer-balance", httphandler.NewTransferBalanceHttpHandler(prop))
				r.Mount("/sales-order", httphandler.NewSalesOrderHttpHandler(prop))
				r.Mount("/merchant", httphandler.NewMerchantHttpHandler(prop))
				r.Mount("/channel", httphandler.NewChannelHttpHandler(prop))
				r.Mount("/divission", httphandler.NewDivissionHttpHandler(prop))
				r.Mount("/top-up-data", httphandler.NewTopUpDataHttpHandler(prop))
			})

			r.Group(func(r chi.Router) {
			})
		})
	})

	chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		logrus.Printf("[HttpHandler] [%s] %s", method[0:3], route)
		return nil
	})

	return r
}
