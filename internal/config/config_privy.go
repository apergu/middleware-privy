package config

import "github.com/sirupsen/logrus"

type Privy struct {
	Host     string
	Username string
	Password string
}

func ReadPrivyConfig(env Getenv) Privy {
	logrus.Info("[config] Start read Privy config")
	defer logrus.Info("[config] End read Privy config")

	conf := Privy{
		Host:     env("PRIVY_GENERAL_HOST"),
		Username: env("PRIVY_GENERAL_USERNAME"),
		Password: env("PRIVY_GENERAL_PASSWORD"),
	}

	logrus.Infof("[config] Privy config : %+v", conf)

	return conf
}
