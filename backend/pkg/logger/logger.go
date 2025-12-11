package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ContextKey string

const RequestIDCtxKey ContextKey = "requestId"

const (
	FieldRequestID       = "request_id"
	FieldOp              = "op"
	FieldUserID          = "user_id"
	FieldGameID          = "game_id"
	FieldLanguage        = "language"
	FieldWord            = "word"
	FieldGuess           = "guess"
	FieldDistance        = "distance"
	FieldUsername        = "username"
	FieldErrorCode       = "error_code"
	FieldDurationMs      = "duration_ms"
	FieldCalculationTime = "calculation_time"
	FieldError           = "error"
)

var defaultLogger *slog.Logger

type Config struct {
	Level        string
	Format       string
	Output       string
	ReportCaller bool
}

func Init(cfg Config) error {
	level := parseLevel(cfg.Level)
	output, cleanup := getOutput(cfg.Output)

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = NewColoredTextHandler(output, &slog.HandlerOptions{
			Level: level,
		})
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)

	if cleanup != nil {
		logFileCleanup = cleanup
	}

	defaultLogger.Info("logger initialized", "level", cfg.Level, "format", cfg.Format, "output", cfg.Output)
	return nil
}

var logFileCleanup func()

func CloseLogFiles() {
	if logFileCleanup != nil {
		logFileCleanup()
		logFileCleanup = nil
	}
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getOutput(output string) (io.Writer, func()) {
	switch strings.ToLower(output) {
	case "stdout", "":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Warn("failed to create log directory, using stdout", "dir", dir, "error", err)
			return os.Stdout, nil
		}

		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Warn("failed to open log file, using stdout", "path", output, "error", err)
			return os.Stdout, nil
		}

		multiWriter := io.MultiWriter(os.Stdout, file)
		cleanup := func() { file.Close() }
		return multiWriter, cleanup
	}
}

func Default() *slog.Logger {
	if defaultLogger == nil {
		return slog.Default()
	}
	return defaultLogger
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := Default()
	if ctx == nil {
		return logger
	}

	if requestID := GetRequestIDFromContext(ctx); requestID != "" {
		return logger.With(FieldRequestID, requestID)
	}
	return logger
}

func Op(ctx context.Context, operation string) *slog.Logger {
	return FromContext(ctx).With(FieldOp, operation)
}

func GetRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id := ctx.Value(RequestIDCtxKey); id != nil {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

type ColoredTextHandler struct {
	opts   *slog.HandlerOptions
	output io.Writer
	mu     *sync.Mutex
	attrs  []slog.Attr
	groups []string
}

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

func NewColoredTextHandler(w io.Writer, opts *slog.HandlerOptions) *ColoredTextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredTextHandler{
		opts:   opts,
		output: w,
		mu:     &sync.Mutex{},
	}
}

func (h *ColoredTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *ColoredTextHandler) Handle(ctx context.Context, r slog.Record) error {
	var levelColor, levelBadge string
	switch r.Level {
	case slog.LevelDebug:
		levelColor, levelBadge = colorCyan, "[DEBUG]"
	case slog.LevelInfo:
		levelColor, levelBadge = colorGreen, "[INFO] "
	case slog.LevelWarn:
		levelColor, levelBadge = colorYellow, "[WARN] "
	case slog.LevelError:
		levelColor, levelBadge = colorRed, "[ERROR]"
	default:
		levelColor, levelBadge = colorReset, "[LOG]  "
	}

	timestamp := r.Time.Format("2006/01/02 - 15:04:05")

	var buf []byte
	buf = fmt.Appendf(buf, "%s%s%s %s%s%s %s|%s %s",
		levelColor, levelBadge, colorReset,
		colorGray, timestamp, colorReset,
		colorGray, colorReset,
		r.Message)

	for _, a := range h.attrs {
		buf = h.appendAttr(buf, a)
	}

	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a)
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.output.Write(buf)
	return err
}

func (h *ColoredTextHandler) appendAttr(buf []byte, a slog.Attr) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}

	if a.Key == FieldRequestID {
		rid := a.Value.String()
		if len(rid) > 8 {
			rid = rid[:8]
		}
		return fmt.Appendf(buf, " %s|%s %s", colorGray, colorReset, rid)
	}

	var color string
	switch a.Key {
	case FieldUserID:
		color = colorMagenta
	case FieldGameID:
		color = colorBlue
	case FieldError:
		color = colorRed
	case FieldWord, "target_word", FieldGuess:
		color = colorYellow
	case FieldDistance, FieldCalculationTime:
		color = colorGreen
	case FieldUsername:
		color = colorCyan
	case FieldLanguage:
		color = colorWhite
	default:
		color = colorGray
	}

	val := a.Value.Any()
	switch v := val.(type) {
	case string:
		if a.Key == FieldWord || a.Key == "target_word" || a.Key == FieldGuess {
			return fmt.Appendf(buf, " %s%s:%s\"%s\"%s", colorGray, a.Key, color, v, colorReset)
		}
		return fmt.Appendf(buf, " %s%s:%s%s%s", colorGray, a.Key, color, v, colorReset)
	case error:
		return fmt.Appendf(buf, " %s%s:%s%v%s", colorGray, a.Key, color, v, colorReset)
	case float64:
		if a.Key == FieldCalculationTime {
			return fmt.Appendf(buf, " %stime:%s%.3fs%s", colorGray, color, v, colorReset)
		}
		return fmt.Appendf(buf, " %s%s:%s%v%s", colorGray, a.Key, color, v, colorReset)
	default:
		return fmt.Appendf(buf, " %s%s:%s%v%s", colorGray, a.Key, color, val, colorReset)
	}
}

func (h *ColoredTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)

	return &ColoredTextHandler{
		opts:   h.opts,
		output: h.output,
		mu:     h.mu,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *ColoredTextHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups), len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups = append(newGroups, name)

	return &ColoredTextHandler{
		opts:   h.opts,
		output: h.output,
		mu:     h.mu,
		attrs:  h.attrs,
		groups: newGroups,
	}
}
