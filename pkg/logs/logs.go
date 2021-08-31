package logs

import (
	"sync"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

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
func Get() *Logger {
	once.Do(func() {
		cfg := config.GetConfig()
		ll := logrus.New()
		ll.SetLevel(getLogLevel(cfg.LogLevel))
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
