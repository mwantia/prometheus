package log

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

var def hclog.Logger

func GetDefaultLogger() hclog.Logger {
	return def
}

func Setup(level string) error {
	def = hclog.New(&hclog.LoggerOptions{
		Name:        "queueverse",
		Level:       parseLevel(level),
		Output:      os.Stdout,
		JSONFormat:  false,
		Color:       hclog.AutoColor,
		TimeFormat:  "02.01.2006 15:04:05",
		DisableTime: false,
	})

	log.SetOutput(io.Discard)
	hclog.SetDefault(def)

	return nil
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
