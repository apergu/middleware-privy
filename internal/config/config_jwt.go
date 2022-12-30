package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rhelper"
)

type JWT struct {
	SignatureKey string
	VerifyKey    string
	Expiration   time.Duration
	Audience     string
	Issuers      string
}

func ReadJWTConfig(env Getenv) JWT {
	logrus.Info("[config] Start read jwt config")
	defer logrus.Info("[config] End read jwt config")

	conf := JWT{
		SignatureKey: env("JWT_SIGNATURE"),
		VerifyKey:    env("JWT_VERIFYKEY"),
		Expiration:   time.Duration(rhelper.ToInt(env("JWT_EXPIRATION"), 5)) * time.Second,
		Audience:     env("JWT_AUD"),
		Issuers:      env("JWT_ISS"),
	}

	logrus.Infof("[config] jwt config : %+v", conf)

	return conf
}
