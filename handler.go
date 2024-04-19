package corelog

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
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
	leveler := namedLeveler(h.name)
	opts := &slog.HandlerOptions{
		AddSource: config.EnableSource,
		Level:     leveler,
	}

	var output io.Writer
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
		handler = slog.NewJSONHandler(output, opts)
	default:
		// default to tint.Handler if no value is set
		// or the set value is invalid
		handler = tint.NewHandler(output, &tint.Options{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: opts.ReplaceAttr,
		})
	}

	if len(h.attrs) > 0 {
		handler = handler.WithAttrs(h.attrs)
	}
	if len(h.group) > 0 {
		handler = handler.WithGroup(h.group)
	}
	return handler.Handle(ctx, record)
}
