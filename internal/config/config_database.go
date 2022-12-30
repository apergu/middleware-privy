package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rhelper"
)

type Database struct {
	Driver          string
	Dsn             string
	DBName          string
	MaxConnLifetime time.Duration
	MaxOpenConn     int32
	MaxConnIdleTime time.Duration
}

func ReadDatabaseConfig(env Getenv) Database {
	logrus.Info("[config] Start read database config")
	defer logrus.Info("[config] End read database config")

	maxopencon := rhelper.ToInt(env("DB_MAX_OPEN_CONN"), 4)
	maxidlecon := rhelper.ToInt(env("DB_MAX_IDLE_CONN"), 20)
	maxlifecon := rhelper.ToInt(env("DB_MAX_LIFE_CONN"), 3)

	conf := Database{
		Driver:          "postgre",
		Dsn:             env("DB_DSN"),
		DBName:          env("DB_NAME"),
		MaxConnLifetime: time.Minute * time.Duration(maxlifecon),
		MaxOpenConn:     int32(maxopencon),
		MaxConnIdleTime: time.Second * time.Duration(maxidlecon),
	}

	logrus.Infof("[config] database config : %+v", conf)

	return conf
}
