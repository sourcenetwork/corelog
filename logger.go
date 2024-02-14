package corelog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

const (
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
	levelError = slog.LevelError
	levelFatal = slog.Level(12)
)

// levelLabels is a mapping of log levels to custom labels.
var levelLabels = map[slog.Level]string{
	levelDebug: "DEBUG",
	levelInfo:  "INFO",
	levelError: "ERROR",
	levelFatal: "FATAL",
}

// Logger is a logger that wraps the slog package.
type Logger struct {
	handler          slog.Handler
	enableStackTrace bool
	enableSource     bool
	// skipExit is used for tests that don't want to call os.Exit when logging fatal.
	skipExit bool
}

// NewLogger returns a new logger configured with the default config.
func NewLogger(name string) *Logger {
	// use overrides if specified
	if val, ok := defaultConfig.Overrides[name]; ok {
		return NewLoggerWithConfig(name, val)
	}
	return NewLoggerWithConfig(name, defaultConfig)
}

// NewLoggerWithConfig returns a new logger configured with the given config.
func NewLoggerWithConfig(name string, config Config) *Logger {
	var level slog.Leveler
	switch config.Level {
	case LevelDebug:
		level = levelDebug
	case LevelInfo:
		level = levelInfo
	case LevelError:
		level = levelError
	case LevelFatal:
		level = levelFatal
	default:
		level = levelDebug
	}

	var output io.Writer
	switch config.Output {
	case OutputStdOut:
		output = os.Stdout
	case OutputStdErr:
		output = os.Stderr
	default:
		output = os.Stderr
	}

	opts := &slog.HandlerOptions{
		AddSource: config.EnableSource,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// replace level labels with custom labels
			if a.Key == slog.LevelKey {
				a.Value = slog.StringValue(levelLabels[level.Level()])
			}
			return a
		},
	}

	var handler slog.Handler
	switch config.Format {
	case FormatText:
		handler = slog.NewTextHandler(output, opts)
	case FormatJSON:
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewJSONHandler(output, opts)
	}

	return &Logger{
		handler:          handler.WithGroup(name),
		enableStackTrace: config.EnableStackTrace,
		enableSource:     config.EnableSource,
	}
}

// Debug logs a message at debug log level.
func (l *Logger) Debug(msg string, args ...any) {
	l.DebugContext(context.Background(), msg, args...)
}

// Info logs a message at info log level.
func (l *Logger) Info(msg string, args ...any) {
	l.InfoContext(context.Background(), msg, args...)
}

// Error logs a message at error log level.
func (l *Logger) Error(msg string, args ...any) {
	l.ErrorContext(context.Background(), msg, args...)
}

// ErrorE logs a message at error log level with an error stacktrace.
func (l *Logger) ErrorE(msg string, err error, args ...any) {
	l.ErrorContextE(context.Background(), msg, err, args...)
}

// Fatal logs a message at fatal log level.
func (l *Logger) Fatal(msg string, args ...any) {
	l.FatalContext(context.Background(), msg, args...)
}

// FatalE logs a message at fatal log level with an error stacktrace.
func (l *Logger) FatalE(msg string, err error, args ...any) {
	l.FatalContextE(context.Background(), msg, err, args...)
}

// DebugContext logs a message at debug log level.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, levelDebug, nil, msg, args)
}

// InfoContext logs a message at info log level.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, levelInfo, nil, msg, args)
}

// ErrorContext logs a message at error log level.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, levelError, nil, msg, args)
}

// ErrorContextE logs a message at error log level with an error stacktrace.
func (l *Logger) ErrorContextE(ctx context.Context, msg string, err error, args ...any) {
	l.log(ctx, levelError, err, msg, args)
}

// FatalContext logs a message at fatal log level and calls os.Exit(1).
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, levelFatal, nil, msg, args)
	if !l.skipExit {
		os.Exit(1)
	}
}

// FatalContextE logs a message at fatal log level with an error stacktrace and calls os.Exit(1).
func (l *Logger) FatalContextE(ctx context.Context, msg string, err error, args ...any) {
	l.log(ctx, levelFatal, err, msg, args)
	if !l.skipExit {
		os.Exit(1)
	}
}

// log wraps calls to the underlying logger so that the caller source can be corrected and
// an optional stacktrace can be included.
func (l *Logger) log(ctx context.Context, level slog.Level, err error, msg string, args []any) {
	if !l.handler.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// add caller source if enabled
	if l.enableSource {
		runtime.Callers(3, pcs[:]) // skip [Callers, log, Info]
	}

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	// add stack trace if enabled
	if err != nil && l.enableStackTrace {
		r.Add("stacktrace", fmt.Sprintf("%+v", err))
	}

	_ = l.handler.Handle(ctx, r)
}
