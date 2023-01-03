package adapters

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LogLevel int
type LogFormat int

const (
	DEBUG LogLevel = iota // Head = 0
	INFO                  // Shoulder = 1
	ERROR                 // Knee = 2
)

const (
	TEXT LogFormat = iota // Head = 0
	JSON                  // Shoulder = 1
)

func (ll LogLevel) String() string {
	return []string{"debug", "info", "error"}[ll]
}

func (lf LogFormat) String() string {
	return []string{"text", "json"}[lf]
}

func NewLogger() *logrus.Logger {
	var formatter utcFormatter
	formatter.f = &logrus.JSONFormatter{}

	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &utcFormatter{
			f: &logrus.JSONFormatter{},
		},
		Level: logrus.InfoLevel,
	}
}

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func SetLogger(logger *logrus.Logger, level string, format string) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case DEBUG.String():
		logLevel = logrus.DebugLevel
	case INFO.String():
		logLevel = logrus.InfoLevel
	case ERROR.String():
		logLevel = logrus.ErrorLevel
	default:
		logLevel = logrus.InfoLevel
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case TEXT.String():
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case JSON.String():
		formatter.f = &logrus.JSONFormatter{}
	default:
		formatter.f = &logrus.JSONFormatter{}
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&formatter)
}
