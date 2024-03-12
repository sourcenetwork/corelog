package corelog

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type namedHandler struct {
	name  string
	attrs []slog.Attr
	group string
}

func (h namedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	leveler := namedLeveler(h.name)
	return level >= leveler.Level()
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
		AddSource:   config.EnableSource,
		Level:       leveler,
		ReplaceAttr: leveler.ReplaceAttr,
	}

	var output io.Writer
	switch config.Output {
	case OutputStderr:
		output = os.Stderr

	case OutputStdout:
		output = os.Stdout

	default:
		output = os.Stderr
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

	if len(h.attrs) > 0 {
		handler = handler.WithAttrs(h.attrs)
	}
	if len(h.group) > 0 {
		handler = handler.WithGroup(h.group)
	}
	return handler.Handle(ctx, record)
}
