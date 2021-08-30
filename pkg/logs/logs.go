package logs

import (
	"sync"

	"github.com/sirupsen/logrus"
)

//Logger based logger struct
type Logger struct {
	*logrus.Logger
}

var (
	logger Logger
	once   sync.Once
)

//Get returns the logger instance with specified parameters
func Get(logLever interface{}) *Logger {
	once.Do(func() {
		ll := logrus.New()
		ll.SetLevel(getLogLevel(logLever.(string)))
		ll.SetReportCaller(true)
		logger = Logger{ll}
	})
	return &logger
}

//getLogLevel returns the log level
func getLogLevel(ll string) logrus.Level {
	switch ll {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.ErrorLevel
	}
}
