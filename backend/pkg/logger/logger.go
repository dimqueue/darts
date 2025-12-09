package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const RequestIDCtxKey ContextKey = "requestId"

type Config struct {
	Level        string
	Format       string
	Output       string
	ReportCaller bool
}

type ConsoleFormatter struct{}

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[97m"
	colorGray    = "\033[90m"
)

func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor, levelBadge string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = colorCyan
		levelBadge = "[DEBUG]"
	case logrus.InfoLevel:
		levelColor = colorGreen
		levelBadge = "[INFO] "
	case logrus.WarnLevel:
		levelColor = colorYellow
		levelBadge = "[WARN] "
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = colorRed
		levelBadge = "[ERROR]"
	default:
		levelColor = colorReset
		levelBadge = "[LOG]  "
	}

	timestamp := entry.Time.Format("2006/01/02 - 15:04:05")

	msg := fmt.Sprintf("%s%s%s %s%s%s %s|%s %s",
		levelColor, levelBadge, colorReset,
		colorGray, timestamp, colorReset,
		colorGray, colorReset,
		entry.Message)

	if rid, ok := entry.Data["request_id"]; ok {
		shortRid := fmt.Sprintf("%v", rid)
		if len(shortRid) > 8 {
			shortRid = shortRid[:8]
		}
		msg += fmt.Sprintf(" %s|%s %s", colorGray, colorReset, shortRid)
	}

	var fields []string

	if val, ok := entry.Data["game_id"]; ok {
		fields = append(fields, fmt.Sprintf("%sgame_id:%s%s%v%s", colorGray, colorBlue, colorWhite, val, colorReset))
	}
	if val, ok := entry.Data["user_id"]; ok {
		fields = append(fields, fmt.Sprintf("%suser_id:%s%v%s", colorGray, colorMagenta, val, colorReset))
	}
	if val, ok := entry.Data["word"]; ok {
		fields = append(fields, fmt.Sprintf("%sword:%s%s\"%v\"%s", colorGray, colorReset, colorYellow, val, colorReset))
	}
	if val, ok := entry.Data["target_word"]; ok {
		fields = append(fields, fmt.Sprintf("%starget:%s%s\"%v\"%s", colorGray, colorReset, colorYellow, val, colorReset))
	}
	if val, ok := entry.Data["guess"]; ok {
		fields = append(fields, fmt.Sprintf("%sguess:%s%s\"%v\"%s", colorGray, colorReset, colorCyan, val, colorReset))
	}
	if val, ok := entry.Data["distance"]; ok {
		fields = append(fields, fmt.Sprintf("%sdistance:%s%s%v%s", colorGray, colorReset, colorGreen, val, colorReset))
	}
	if val, ok := entry.Data["error"]; ok {
		fields = append(fields, fmt.Sprintf("%serror:%s%s%v%s", colorGray, colorReset, colorRed, val, colorReset))
	}
	if val, ok := entry.Data["username"]; ok {
		fields = append(fields, fmt.Sprintf("%susername:%s%s%v%s", colorGray, colorReset, colorCyan, val, colorReset))
	}
	if val, ok := entry.Data["language"]; ok {
		fields = append(fields, fmt.Sprintf("%slang:%s%v%s", colorGray, colorWhite, val, colorReset))
	}
	if val, ok := entry.Data["calculation_time"]; ok {
		fields = append(fields, fmt.Sprintf("%stime:%s%.3fs%s", colorGray, colorGreen, val, colorReset))
	}

	if len(fields) > 0 {
		msg += " " + strings.Join(fields, " ")
	}

	return []byte(msg + "\n"), nil
}

func Init(config Config) error {
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warnf("Invalid log level %s, defaulting to info", config.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&ConsoleFormatter{})

	switch strings.ToLower(config.Output) {
	case "stdout", "":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	default:
		dir := filepath.Dir(config.Output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logrus.Warnf("Failed to create log directory %s: %v", dir, err)
		}

		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Warnf("Failed to open log file %s: %v, using stdout only", config.Output, err)
			logrus.SetOutput(os.Stdout)
		} else {
			multiWriter := io.MultiWriter(os.Stdout, file)
			logrus.SetOutput(multiWriter)

			jsonFile, err := os.OpenFile(config.Output+".json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				logrus.AddHook(&JSONFileHook{file: jsonFile})
			}
		}
	}

	logrus.SetReportCaller(config.ReportCaller)
	logrus.Infof("Logger initialized - Level: %s, Output: %s", config.Level, config.Output)

	return nil
}

type JSONFileHook struct {
	file      *os.File
	formatter *logrus.JSONFormatter
}

func (hook *JSONFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *JSONFileHook) Fire(entry *logrus.Entry) error {
	if hook.formatter == nil {
		hook.formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	}
	line, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.file.Write(line)
	return err
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithContext(ctx context.Context) *logrus.Entry {
	entry := logrus.NewEntry(logrus.StandardLogger())
	if ctx == nil {
		return entry
	}
	if requestID := getRequestIDFromContext(ctx); requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}
	return entry
}

func getRequestIDFromContext(ctx context.Context) string {
	if id := ctx.Value(RequestIDCtxKey); id != nil {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}

func ErrorCtx(ctx context.Context, msg string) {
	WithContext(ctx).Error(msg)
}

func WarnCtx(ctx context.Context, msg string) {
	WithContext(ctx).Warn(msg)
}

func InfoCtx(ctx context.Context, msg string) {
	WithContext(ctx).Info(msg)
}

func DebugCtx(ctx context.Context, msg string) {
	WithContext(ctx).Debug(msg)
}
