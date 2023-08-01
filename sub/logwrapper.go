package sub

import (
	"github.com/rs/zerolog"
	"gitlab.com/etke.cc/go/validator"
)

type validatorLoggerWrapper struct {
	log *zerolog.Logger
}

// NewLogWrapper creates a wrapper around zerolog.Logger to implement validator.Logger interface
func NewLogWrapper(log *zerolog.Logger) validator.Logger {
	return &validatorLoggerWrapper{log}
}

func (l validatorLoggerWrapper) Info(msg string, args ...interface{}) {
	l.log.Info().Msgf(msg, args...)
}

func (l validatorLoggerWrapper) Error(msg string, args ...interface{}) {
	l.log.Warn().Msgf(msg, args...)
}
