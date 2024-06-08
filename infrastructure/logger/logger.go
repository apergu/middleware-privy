package logger

const (
	DefaultLogFileName = "service"
	LogWarn            = "warn"
	LogInfo            = "info"
	LogError           = "error"
)

type Log struct {
	Event      string
	StatusCode int
	Method     string
	Request    interface{}
	URL        string
	Message    string
	Response   interface{}
	Service    string
	FileName   string
	Line       int
}
