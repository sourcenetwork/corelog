package corelog

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/term"
)

const (
	// nameKey is the key for the logger name attribute
	nameKey = "$name"
	// stackKey is the key for the logger stack attribute
	stackKey = "$stack"
	// errorKey is the key for the logger error attribute
	errorKey = "$err"
	// msgKey is the key for the logger message attribute
	msgKey = "$msg"
	// timeKey is the key for the logger time attribute
	timeKey = "$time"
	// levelKey is the key for the logger level attribute
	levelKey = "$level"
	// sourceKey is the key for the logger source attribute
	sourceKey = "$source"
)

type namedHandler struct {
	name  string
	attrs []slog.Attr
	group string
}

func (h namedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= namedLeveler(h.name).Level()
}

func (h namedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &namedHandler{
		name:  h.name,
		group: h.group,
		attrs: attrs,
	}
}

func (h namedHandler) WithGroup(name string) slog.Handler {
	return &namedHandler{
		name:  h.name,
		attrs: h.attrs,
		group: name,
	}
}

func (h namedHandler) Handle(ctx context.Context, record slog.Record) error {
	config := GetConfig(h.name)

	var output *os.File
	switch config.Output {
	case OutputStdout:
		output = os.Stdout
	default:
		// default to os.Stderr if no value is set
		// or the set value is invalid
		output = os.Stderr
	}

	var handler slog.Handler
	switch config.Format {
	case FormatJSON:
		handler = newJSONHandler(config, h.name, output)
	default:
		// default to tint.Handler if no value is set
		// or the set value is invalid
		handler = newTintHandler(config, h.name, output)
		// prepend logger name to message
		record.Message = h.name + " " + record.Message
	}

	if len(h.attrs) > 0 {
		handler = handler.WithAttrs(h.attrs)
	}
	if len(h.group) > 0 {
		handler = handler.WithGroup(h.group)
	}
	return handler.Handle(ctx, record)
}

func newTintHandler(config Config, name string, output *os.File) slog.Handler {
	isTerminal := term.IsTerminal(int(output.Fd()))
	return tint.NewHandler(output, &tint.Options{
		AddSource: config.EnableSource,
		Level:     namedLeveler(name),
		NoColor:   !isTerminal || config.DisableColor,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			// ignore name as it is prended to message
			if attr.Key == nameKey {
				return slog.Attr{}
			}
			return attr
		},
	})
}

func newJSONHandler(config Config, name string, output *os.File) *slog.JSONHandler {
	return slog.NewJSONHandler(output, &slog.HandlerOptions{
		AddSource: config.EnableSource,
		Level:     namedLeveler(name),
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			switch attr.Key {
			case slog.TimeKey:
				attr.Key = timeKey
			case slog.LevelKey:
				attr.Key = levelKey
			case slog.MessageKey:
				attr.Key = msgKey
			case slog.SourceKey:
				attr.Key = sourceKey
			}
			return attr
		},
	})
}
