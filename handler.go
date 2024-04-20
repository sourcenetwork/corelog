package corelog

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const (
	// nameKey is the key for the logger name attribute
	nameKey = "$name"
	// stackKey is the key for the logger stack attribute
	stackKey = "$stack"
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
		handler = slog.NewJSONHandler(output, opts)
	default:
		// default to tint.Handler if no value is set
		// or the set value is invalid
		handler = tint.NewHandler(output, &tint.Options{
			AddSource: opts.AddSource,
			Level:     opts.Level,
			// disable color if not a tty or config requested no color
			NoColor: !isatty.IsTerminal(output.Fd()) || config.DisableColor,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if attr.Key == nameKey {
					return slog.Attr{} // ignore name as it is prended to message
				}
				return attr
			},
		})
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
