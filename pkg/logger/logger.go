// Package logger provides structured logging functionality for the application.
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger represents the application logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance with the specified level
func NewLogger(level string) *Logger {
	log := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// Set formatter
	if level == "debug" {
		// Use text formatter for development
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	} else {
		// Use JSON formatter for production
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output
	log.SetOutput(os.Stdout)

	return &Logger{log}
}

// WithFields creates a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithField creates a logger with a single additional field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithError creates a logger with an error field
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithRequest creates a logger with request information
func (l *Logger) WithRequest(method, path, userAgent string) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{
		"method":     method,
		"path":       path,
		"user_agent": userAgent,
	})
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

// Default logger instance
var defaultLogger *Logger

// InitDefaultLogger initializes the default logger
func InitDefaultLogger(level string) {
	defaultLogger = NewLogger(level)
}

// GetDefaultLogger returns the default logger instance
func GetDefaultLogger() *Logger {
	if defaultLogger == nil {
		defaultLogger = NewLogger("info")
	}
	return defaultLogger
}

// Convenience functions using default logger
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}
