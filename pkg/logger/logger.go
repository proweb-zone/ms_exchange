package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger() *zerolog.Logger {
	var logger zerolog.Logger

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	return &logger
}

type Logger struct {
	Logger *zerolog.Logger
}

func NewLogger(logger *zerolog.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) GetLogger() *zerolog.Logger {
	return l.Logger
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.Logger.Info().Fields(fields).Msg(msg)
}

func (l *Logger) Error(err error, msg string, fields map[string]interface{}) {
	l.Logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Fatal(err error, msg string, fields map[string]interface{}) {
	l.Logger.Fatal().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Warn(msg string, fields map[string]interface{}, err error) {
	l.Logger.Warn().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.Logger.Debug().Fields(fields).Msg(msg)
}
