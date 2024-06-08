package logrus

import (
	// "middleware/config"
	// "middleware/helper/utils"
	"middleware/infrastructure/logger"
	"middleware/internal/config"

	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// type Logrus struct {
// 	cfg *config.Config
// }

// func NewLoggerLogrus(cfg *config.Config) *Logrus {
// 	return &Logrus{cfg: cfg}
// }

type Logrus struct {
	cfg *config.Config
}

func NewLoggerLogrus(cfg *config.Config) *Logrus {
	return &Logrus{cfg: cfg}
}

func (l *Logrus) CreateLog(data *logger.Log, types, logFileName string) error {

	// cg := l.cfg

	log := logrus.New()
	// You could set this to any `io.Writer` such as a file
	baseDir := ""
	_, fileName, line, _ := runtime.Caller(1)

	projectDir, _ := filepath.Abs(".")

	relativePath := strings.TrimPrefix(fileName, projectDir+"/")
	data.FileName = relativePath
	data.Line = line
	if os.Getenv("ENVIRONMENT") == "preproduction" || os.Getenv("ENVIRONMENT") == "production" {
		baseDir = "/"
	}

	file, err := os.OpenFile(baseDir+"logs/"+logFileName+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	logField := logrus.Fields{
		"event": data.Event,
	}

	if logFileName == logger.DefaultLogFileName {
		logField = logrus.Fields{
			"event":       data.Event,
			"status_code": data.StatusCode,
			"filename":    data.FileName,
			"line":        data.Line,
			"response":    data.Response,
			"url":         data.URL,
			"service":     data.Service,
			"request":     data.Request,
		}
	}

	if types == logger.LogWarn {
		log.WithFields(logField).Warn(data.Message)
	}

	if types == logger.LogInfo {
		log.WithFields(logField).Info(data.Message)
	}

	if types == logger.LogError {
		log.WithFields(logField).Error(data.Message)
	}

	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout

	return nil
}
