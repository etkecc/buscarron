package utils

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var loglevel zerolog.Level

// SetLogLevel sets the log level
func SetLogLevel(level string) {
	var err error
	loglevel, err = zerolog.ParseLevel(level)
	if err != nil {
		loglevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(loglevel)
}

// NewContext creates a new context with a logger
func NewContext(parent ...context.Context) context.Context {
	ctx := context.Background()
	if len(parent) > 0 {
		ctx = parent[0]
	}
	return newLogger().WithContext(ctx)
}

func newLogger() zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}}
	return zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
}
