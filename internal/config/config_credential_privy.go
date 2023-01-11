package config

import "github.com/sirupsen/logrus"

type CredentialPrivy struct {
	Host     string
	Username string
	Password string
}

func ReadCredentialPrivyConfig(env Getenv) CredentialPrivy {
	logrus.Info("[config] Start read CredentialPrivy config")
	defer logrus.Info("[config] End read CredentialPrivy config")

	conf := CredentialPrivy{
		Host:     env("PRIVY_HOST"),
		Username: env("PRIVY_USERNAME"),
		Password: env("PRIVY_PASSWORD"),
	}

	logrus.Infof("[config] CredentialPrivy config : %+v", conf)

	return conf
}
