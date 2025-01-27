package log

import (
	"github.com/hashicorp/go-hclog"
	"github.com/mwantia/queueverse/internal/log"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Trace(msg string, args ...interface{})

	Named(name string) Logger
}

type loggerImpl struct {
	hclog.Logger
}

func (l *loggerImpl) Named(name string) Logger {
	return &loggerImpl{
		Logger: l.Logger.Named(name),
	}
}

func New(name string) Logger {
	def := log.GetDefaultLogger()
	if def == nil {
		panic("Logger is not initialized.")
	}

	return &loggerImpl{
		Logger: def.Named(name),
	}
}
