// Package log provides a logger interface for logger libraries
// so that bee does not depend on any of them directly.
// It also includes a default implementation using Logrus (used by bee previously).
package log

import (
	adapters "github.com/luminosita/honeycomb/pkg/log/adapters"
	"github.com/sirupsen/logrus"
	"sync"
)

// Logger serves as an adapter interface for logger libraries
// so that bee does not depend on any of them directly.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var (
	once   sync.Once
	logger Logger
)

func Log() Logger {
	once.Do(func() { // <-- atomic, does not allow repeating
		logger = adapters.NewLogger()
	})

	return logger
}

func SetLogger(level string, format string) {
	adapters.SetLogger(Log().(*logrus.Logger), level, format)
}
