package log

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

var Default hclog.Logger

type HighlightFormatter struct {
	NoColor bool `bool:"no_color"`
}

type Logger struct {
	hclog.Logger
}

func Setup(level string, nc bool) error {
	Default = hclog.New(&hclog.LoggerOptions{
		Name:        "prometheus",
		Level:       parseLevel(level),
		Output:      os.Stdout,
		JSONFormat:  false,
		Color:       hclog.AutoColor,
		TimeFormat:  "02.01.2006 15:04:05",
		DisableTime: false,
	})

	log.SetOutput(io.Discard)
	hclog.SetDefault(Default)

	return nil
}

func New(name string) *Logger {
	return &Logger{
		Logger: Default.Named(name),
	}
}

func parseLevel(level string) hclog.Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return hclog.Trace
	case "DEBUG":
		return hclog.Debug
	case "INFO":
		return hclog.Info
	case "WARN":
		return hclog.Warn
	case "ERROR":
		return hclog.Error
	default:
		return hclog.Info
	}
}
