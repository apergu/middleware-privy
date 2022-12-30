package config

import (
	"net/http"
	"strings"

	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rhelper"
)

func ReadCorsConfig(env Getenv) cors.Options {
	logrus.Info("[config] Start read cors config")
	defer logrus.Info("[config] End read cors config")

	var headers []string = []string{
		"Accept",
		"Authorization",
		"Content-Type",
		"X-CSRF-Token",
	}

	var origins []string = []string{"*"}
	var methods []string = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}

	corsHeadersEnv := strings.TrimSpace(env("CORS_HEADERS"))
	corsOriginsEnv := strings.TrimSpace(env("CORS_ORIGINS"))
	corsMethodsEnv := strings.TrimSpace(env("CORS_METHODS"))

	if corsOriginsEnv != "" {
		origins = strings.Split(corsOriginsEnv, ",")
	}

	if corsHeadersEnv != "" {
		headers = strings.Split(corsHeadersEnv, "")
	}

	if corsMethodsEnv != "" {
		methods = strings.Split(corsMethodsEnv, "")
	}

	maxAgeInSeconds := rhelper.ToInt(env("CORS_MAX_AGE"), 60)

	conf := cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
		MaxAge:         maxAgeInSeconds,
	}

	logrus.Infof("[config] cors config : %+v", conf)

	return conf
}
