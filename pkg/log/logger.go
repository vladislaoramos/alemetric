// Package logger provides a mechanism for logging incident.
package logger

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogInterface defines the interface of interaction between a client and the tool.
type LogInterface interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// Logger is the logger implementing LogInterface.
type Logger struct {
	logger *logrus.Logger
}

// New creates an object of Logger.
// The logger allows to specify the level and the output mode.
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

// Debug is the debug method for Logger.
func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Info is the info mode for Logger.
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Warn is the warn mode for Logger.
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Error is the debug mode for Logger.
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Fatal is the fatal mode for Logger.
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}
