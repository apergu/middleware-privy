package config

import "github.com/sirupsen/logrus"

type BasicAuth struct {
	Username string
	Password string
}

func ReadBasicAuthConfig(env Getenv) BasicAuth {
	logrus.Info("[config] Start read basic auth config")
	defer logrus.Info("[config] End read basic auth config")

	conf := BasicAuth{
		Username: env("BASIC_AUTH_USERNAME"),
		Password: env("BASIC_AUTH_PASSWORD"),
	}

	logrus.Infof("[config] basic auth config : %+v", conf)

	return conf
}
