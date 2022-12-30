package config

import (
	"os"

	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Getenv func(string) string

type Config struct {
	Application Application
	Database    Database
	Cors        cors.Options
	Jwt         JWT
	RefreshJWT  RefreshJWT
	BasicAuth   BasicAuth
}

func Init() *Config {
	appConfig := new(Config)
	ReadLoggerConfig()

	logrus.Info("[config] Start init config")
	defer logrus.Info("[config] End init config")

	appConfig.Application = ReadApplicationConfig(os.Getenv)
	appConfig.Database = ReadDatabaseConfig(os.Getenv)
	appConfig.Cors = ReadCorsConfig(os.Getenv)
	appConfig.Jwt = ReadJWTConfig(os.Getenv)
	appConfig.RefreshJWT = ReadRefreshJWTConfig(os.Getenv)
	appConfig.BasicAuth = ReadBasicAuthConfig(os.Getenv)

	return appConfig
}
