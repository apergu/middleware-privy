package config

import "github.com/sirupsen/logrus"

type ERPPrivy struct {
	Host           string
	Username       string
	Password       string
	ApplicationKey string
	RequestId      string
}

func ReadERPPrivyConfig(env Getenv) ERPPrivy {
	logrus.Info("[config] Start read ERP Privy config")
	defer logrus.Info("[config] End read ERP Privy config")

	conf := ERPPrivy{
		Host:           env("ERP_PRIVY_HOST"),
		Username:       env("ERP_PRIVY_USERNAME"),
		Password:       env("ERP_PRIVY_PASSWORD"),
		ApplicationKey: env("ERP_PRIVY_APPLICATION_KEY"),
		RequestId:      env("ERP_PRIVY_REQUEST_ID"),
	}

	logrus.Infof("[config] Privy config : %+v", conf)

	return conf
}
