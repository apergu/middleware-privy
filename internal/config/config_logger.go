package config

import (
	"github.com/sirupsen/logrus"
)

func ReadLoggerConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// logrus.SetLevel(cfg.Logger.Level)
	// logrus.AddHook(&apmlogrus.Hook{
	// 	LogLevels: logrus.AllLevels,
	// })
	logrus.SetReportCaller(false)
}
