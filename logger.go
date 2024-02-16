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
	return NewLoggerWithConfig(name, LoadConfig())
}

// NewLoggerWithConfig returns a new logger configured with the given config.
func NewLoggerWithConfig(name string, config Config) *Logger {
	if config.Overrides == nil {
		config.Overrides = make(map[string]Config)
	}
	// use overrides if specified
	if val, ok := config.Overrides[name]; ok {
		return NewLoggerWithConfig(name, val)
	}

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
		level = levelInfo
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
		AddSource:   config.EnableSource,
		Level:       level,
		ReplaceAttr: replaceLoggerLevel(level.Level()),
	}

	var handler slog.Handler
	switch config.Format {
	case FormatText:
		handler = slog.NewTextHandler(output, opts)
	case FormatJSON:
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	// add logger name to all logs
	handler = handler.WithAttrs([]slog.Attr{slog.Any("logger", name)})

	return &Logger{
		handler:          handler,
		enableStackTrace: config.EnableStackTrace,
		enableSource:     config.EnableSource,
	}
}

// WithAttrs returns a new Logger whose attributes consist of
// both the receiver's attributes and the arguments.
func (l *Logger) WithAttrs(attrs ...slog.Attr) *Logger {
	return &Logger{
		handler:          l.handler.WithAttrs(attrs),
		enableStackTrace: l.enableStackTrace,
		enableSource:     l.enableSource,
		skipExit:         l.skipExit,
	}
}

// WithGroup returns a new Logger with the given group appended to
// the receiver's existing groups.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		handler:          l.handler.WithGroup(name),
		enableStackTrace: l.enableStackTrace,
		enableSource:     l.enableSource,
		skipExit:         l.skipExit,
	}
}

// Debug logs a message at debug log level.
func (l *Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), levelDebug, nil, msg, args)
}

// Info logs a message at info log level.
func (l *Logger) Info(msg string, args ...any) {
	l.log(context.Background(), levelInfo, nil, msg, args)
}

// Error logs a message at error log level.
func (l *Logger) Error(msg string, args ...any) {
	l.log(context.Background(), levelError, nil, msg, args)
}

// ErrorE logs a message at error log level with an error stacktrace.
func (l *Logger) ErrorE(msg string, err error, args ...any) {
	l.log(context.Background(), levelError, err, msg, args)
}

// Fatal logs a message at fatal log level.
func (l *Logger) Fatal(msg string, args ...any) {
	l.log(context.Background(), levelFatal, nil, msg, args)
	if !l.skipExit {
		os.Exit(1)
	}
}

// FatalE logs a message at fatal log level with an error stacktrace.
func (l *Logger) FatalE(msg string, err error, args ...any) {
	l.log(context.Background(), levelFatal, err, msg, args)
	if !l.skipExit {
		os.Exit(1)
	}
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
