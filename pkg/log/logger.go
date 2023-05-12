package logger

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogInterface interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Logger struct {
	logger *logrus.Logger
}

func New(level string, output io.Writer) *Logger {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "debug":
		l = logrus.DebugLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	case "info":
		l = logrus.InfoLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetOutput(output)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}
