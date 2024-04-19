package corelog

import (
	"context"
	"fmt"
	"log/slog"
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

// Info logs a message at info log level.
func (l *Logger) Info(msg string, args ...slog.Attr) {
	l.log(context.Background(), slog.LevelInfo, nil, msg, args)
}

// InfoContext logs a message at info log level.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelInfo, nil, msg, args)
}

// Error logs a message at error log level.
func (l *Logger) Error(msg string, args ...slog.Attr) {
	l.log(context.Background(), slog.LevelError, nil, msg, args)
}

// ErrorE logs a message at error log level with an error stacktrace.
func (l *Logger) ErrorE(msg string, err error, args ...slog.Attr) {
	l.log(context.Background(), slog.LevelError, err, msg, args)
}

// ErrorContext logs a message at error log level.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...slog.Attr) {
	l.log(ctx, slog.LevelError, nil, msg, args)
}

// ErrorContextE logs a message at error log level with an error stacktrace.
func (l *Logger) ErrorContextE(ctx context.Context, msg string, err error, args ...slog.Attr) {
	l.log(ctx, slog.LevelError, err, msg, args)
}

// log wraps calls to the underlying logger so that the caller source can be corrected and
// an optional stacktrace can be included.
func (l *Logger) log(ctx context.Context, level slog.Level, err error, msg string, args []slog.Attr) {
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
	// add logger name
	r.Add("name", l.name)
	// add remaining attributes
	r.AddAttrs(args...)

	// add stack trace if enabled
	if err != nil && config.EnableStackTrace {
		r.Add("stack", fmt.Sprintf("%+v", err))
	}

	_ = l.handler.Handle(ctx, r)
}
