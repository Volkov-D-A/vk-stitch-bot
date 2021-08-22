package logs

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func Get(logLever string) *Logger {
		ll := logrus.New()
		ll.SetLevel(getLogLevel(logLever))
		ll.SetReportCaller(true)
		return &Logger{ll}
}

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