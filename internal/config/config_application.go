package config

import "github.com/sirupsen/logrus"

type Application struct {
	Mode    string
	Port    string
	Name    string
	IsRedis bool
	IsS3    bool
}

func ReadApplicationConfig(env Getenv) Application {
	logrus.Info("[config] Start read application config")
	defer logrus.Info("[config] End read application config")

	conf := Application{
		Mode:    env("APP_MODE"),
		Port:    env("APP_PORT"),
		Name:    env("APP_NAME"),
		IsRedis: env("IS_REDIS") == "1",
		IsS3:    env("IS_S3") == "1",
	}

	logrus.Infof("[config] application config : %+v", conf)

	return conf
}
