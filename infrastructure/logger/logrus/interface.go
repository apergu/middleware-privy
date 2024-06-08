package logrus

import "middleware/infrastructure/logger"

type LogrusInterface interface {
	CreateLog(data *logger.Log, types, logFileName string) error
}
