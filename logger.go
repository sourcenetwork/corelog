package corelog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger is a logger that wraps the slog package.
type Logger struct {
	name    string
	handler slog.Handler
}

// NewLogger returns a new named logger.
func NewLogger(name string) *Logger {
	return &Logger{
		name:    name,
		handler: &namedHandler{name: name},
	}
}

// WithAttrs returns a new Logger whose attributes consist of
// both the receiver's attributes and the arguments.
func (l *Logger) WithAttrs(attrs ...slog.Attr) *Logger {
	return &Logger{
		name:    l.name,
		handler: l.handler.WithAttrs(attrs),
	}
}

// WithGroup returns a new Logger with the given group appended to
// the receiver's existing groups.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		name:    l.name,
		handler: l.handler.WithGroup(name),
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
	os.Exit(1)
}

// FatalE logs a message at fatal log level with an error stacktrace.
func (l *Logger) FatalE(msg string, err error, args ...any) {
	l.log(context.Background(), levelFatal, err, msg, args)
	os.Exit(1)
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
	os.Exit(1)
}

// FatalContextE logs a message at fatal log level with an error stacktrace and calls os.Exit(1).
func (l *Logger) FatalContextE(ctx context.Context, msg string, err error, args ...any) {
	l.log(ctx, levelFatal, err, msg, args)
	os.Exit(1)
}

// log wraps calls to the underlying logger so that the caller source can be corrected and
// an optional stacktrace can be included.
func (l *Logger) log(ctx context.Context, level slog.Level, err error, msg string, args []any) {
	// check if logger is enabled
	if !l.handler.Enabled(ctx, level) {
		return
	}

	// use latest config values
	config := GetConfig(l.name)

	var pcs [1]uintptr
	// add caller source if enabled
	if config.EnableSource {
		runtime.Callers(3, pcs[:]) // skip [Callers, log, Info]
	}

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(slog.Any("name", l.name))
	r.Add(args...)

	// add stack trace if enabled
	if err != nil && config.EnableStackTrace {
		r.Add("stacktrace", fmt.Sprintf("%+v", err))
	}

	_ = l.handler.Handle(ctx, r)
}
