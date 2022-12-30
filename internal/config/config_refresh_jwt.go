package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rhelper"
)

type RefreshJWT struct {
	SignatureKey string
	VerifyKey    string
	Expiration   time.Duration
	Audience     string
	Issuers      string
}

func ReadRefreshJWTConfig(env Getenv) RefreshJWT {
	logrus.Info("[config] Start read refresh jwt config")
	defer logrus.Info("[config] End read refresh jwt config")

	conf := RefreshJWT{
		SignatureKey: env("REFRESH_JWT_SIGNATURE"),
		VerifyKey:    env("REFRESH_JWT_VERIFYKEY"),
		Expiration:   time.Duration(rhelper.ToInt(env("REFRESH_JWT_EXPIRATION"), 5)) * time.Second,
		Audience:     env("REFRESH_JWT_AUD"),
		Issuers:      env("REFRESH_JWT_ISS"),
	}

	logrus.Infof("[config] refresh jwt config : %+v", conf)

	return conf
}
