package logger

import (
	"io"
	"os"

	log "github.com/rs/zerolog"

	"github.com/Ropho/avito-bootcamp-assignment/config"
)

type Logger struct {
	logger *log.Logger
}

func (logger Logger) Fatal(msg string, err error) {
	logger.logger.Fatal().Msgf(msg, err)
}

func (logger Logger) Infof(msg string, params ...any) {
	logger.logger.Info().Msgf(msg, params...)
}

func (logger Logger) Errorf(err error, msg string, params ...any) {
	logger.logger.Error().Err(err).Msgf(msg, params...)
}

func NewLogger(cfg *config.LoggerConfig) Logger {

	logger := log.New(io.MultiWriter(os.Stderr, cfg.File)).
		Level(cfg.LogLevel).
		With().
		Timestamp().
		Logger()

	return Logger{
		logger: &logger,
	}
}
