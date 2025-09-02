//go:build silent

package corelog

import (
	"context"
	"log/slog"
	"os"
)

// This dummies out all of the logging functionality, so that code using the
// logger will be silent, if the build tag silent is used.

const (
	nameKey   = "$name"
	stackKey  = "$stack"
	errorKey  = "$err"
	msgKey    = "$msg"
	timeKey   = "$time"
	levelKey  = "$level"
	sourceKey = "$source"
)

type namedHandler struct {
	name  string
	attrs []slog.Attr
	group string
}

func (h namedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}
func (h namedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}
func (h namedHandler) WithGroup(name string) slog.Handler {
	return h
}
func (h namedHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func newTintHandler(_ Config, _ string, _ *os.File) slog.Handler {
	return namedHandler{}
}
func newJSONHandler(_ Config, _ string, _ *os.File) *slog.JSONHandler {
	return slog.NewJSONHandler(nil, nil)
}
