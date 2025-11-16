package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Config holds logger configuration
type Config struct {
	Level        string
	Format       string
	Output       string
	ReportCaller bool
}

// Init initializes the logger with the provided configuration
func Init(config Config) error {
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warnf("Invalid log level %s, defaulting to info", config.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set formatter
	switch strings.ToLower(config.Format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	}

	// Set output
	switch strings.ToLower(config.Output) {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	case "":
		logrus.SetOutput(os.Stdout)
	default:
		// If a file path is provided, open it
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Warnf("Failed to open log file %s: %v, using stdout", config.Output, err)
			logrus.SetOutput(os.Stdout)
		} else {
			// Write to both file and stdout
			multiWriter := io.MultiWriter(os.Stdout, file)
			logrus.SetOutput(multiWriter)
		}
	}

	// Set report caller
	logrus.SetReportCaller(config.ReportCaller)

	logrus.Infof("Logger initialized - Level: %s, Format: %s, Output: %s", config.Level, config.Format, config.Output)

	return nil
}

// GetLogger returns the standard logger instance
func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

// WithFields creates a new logger entry with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

// WithField creates a new logger entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}
