package logs

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var (
	logger Logger
	once   sync.Once
)

func Get(logLever string) *Logger {
	once.Do(func() {
		logrusLogger := logrus.New()
		switch logLever {
		case "Debug":
			logrusLogger.SetLevel(logrus.DebugLevel)
		}
		logger = Logger{logrusLogger}
	})
	return &logger
}
