package helper

import (
	"io"
	"middleware/infrastructure/logger"
	"middleware/internal/constants"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Event struct {
	ErpTopupNo string `json:"erp_topup_no,omitempty"`
	Name       string `json:"name"`
}

type typeLog string

const (
	LogService  typeLog = "service"
	LogEvent    typeLog = "event"
	LogError    typeLog = "error"
	LogValidate typeLog = "validate"
	LogRequest  typeLog = "request"
)

type LogStruct struct {
	Env   string `json:"env"`
	Event struct {
		ErpTopupNo string `json:"erp_topup_no,omitempty"`
		Name       string `json:"name"`
	} `json:"event"`
	file            string
	Func            string  `json:"function"`
	Line            int     `json:"line"`
	Message         string  `json:"message"`
	ProcessTime     float64 `json:"process_time"`
	ProcessTimeUnit string  `json:"process_time_unit"`
	RequestID       string  `json:"request_id"`
	RequestIP       string  `json:"request_ip"`
	RequestMethod   string  `json:"request_method"`
	RequestPath     string  `json:"request_path"`
	Service         string  `json:"service"`
	StartTime       string  `json:"start_time"`
}

func LoggerValidateStructfunc(w http.ResponseWriter, r *http.Request, event, serviceName, message, id string) {
	dataLog := LogStruct{
		Env:             os.Getenv("APP_MODE"),
		Event:           Event{ErpTopupNo: id, Name: event},
		Message:         message,
		ProcessTime:     float64(time.Since(r.Context().Value(constants.Timestamp).(time.Time)).Seconds()),
		ProcessTimeUnit: "second",
		RequestID:       r.Header.Get("X-Request-Id"),
		RequestIP:       readUserIP(r),
		RequestMethod:   r.Method,
		RequestPath:     r.RequestURI,
		Service:         serviceName,
		StartTime:       r.Context().Value(constants.Timestamp).(time.Time).Format(time.RFC3339),
	}

	CreateLogERPPRivy(dataLog, "warn", serviceName, string(LogValidate))
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func CreateLogERPPRivy(data LogStruct, typeLog, typeService, logFileName string) error {
	log := logrus.New()

	counter, fileName, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(counter)
	parts := strings.Split(details.Name(), ".")
	data.Func = parts[len(parts)-1]

	data.Line = line
	projectDir, _ := filepath.Abs(".")

	relativePath := strings.TrimPrefix(fileName, projectDir+"/")
	data.file = relativePath

	logField := logrus.Fields{}
	filePathLog := "logs/" + typeService + ".log"

	ensureDir(filePathLog)

	file, _ := os.OpenFile(filePathLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	writer := io.MultiWriter(os.Stderr, file)
	log.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(writer)

	switch logFileName {
	case string(LogValidate):
		logField = logrus.Fields{
			"env":               data.Env,
			"event":             data.Event,
			"line":              data.Line,
			"file":              data.file,
			"message":           data.Message,
			"function":          data.Func,
			"process_time":      data.ProcessTime,
			"process_time_unit": data.ProcessTimeUnit,
			"request_id":        data.RequestID,
			"request_ip":        data.RequestIP,
			"request_method":    data.RequestMethod,
			"request_path":      data.RequestPath,
			"service":           data.Service,
			"start_time":        data.StartTime,
		}
	case string(LogEvent):
		logField = logrus.Fields{}
	}

	if typeLog == logger.LogWarn {
		log.WithFields(logField).Warn()
	}

	if typeLog == logger.LogInfo {
		log.WithFields(logField).Info()
	}

	if typeLog == logger.LogError {
		log.WithFields(logField).Error()
	}

	return nil
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}
