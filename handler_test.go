package corelog

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestHandler struct {
	attrs   []slog.Attr
	group   string
	records []slog.Record
	level   slog.Level
}

var _ (slog.Handler) = (*TestHandler)(nil)

func (h *TestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *TestHandler) Handle(ctx context.Context, record slog.Record) error {
	h.records = append(h.records, record)
	return nil
}

func (h *TestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TestHandler{attrs: attrs, group: h.group, records: h.records}
}

func (h *TestHandler) WithGroup(name string) slog.Handler {
	return &TestHandler{attrs: h.attrs, group: name, records: h.records}
}

func TestHandlerWithLevelInfo(t *testing.T) {
	SetConfig(Config{Level: LevelInfo})
	handler := namedHandler{name: "test"}

	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, handler.Enabled(context.Background(), slog.LevelError))
}

func TestHandlerWithLevelError(t *testing.T) {
	SetConfig(Config{Level: LevelError})
	handler := namedHandler{name: "test"}

	assert.False(t, handler.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, handler.Enabled(context.Background(), slog.LevelError))
}

func TestHandlerWithAttrs(t *testing.T) {
	handler := namedHandler{name: "test"}
	attrs := []slog.Attr{slog.Any("extra", "value")}
	other := handler.WithAttrs(attrs)

	otherHandler, ok := other.(*namedHandler)
	require.True(t, ok)
	assert.Equal(t, "test", otherHandler.name)
	assert.Equal(t, "", otherHandler.group)
	assert.Equal(t, attrs, otherHandler.attrs)
}

func TestHandlerWithGroup(t *testing.T) {
	handler := namedHandler{name: "test"}
	other := handler.WithGroup("group")

	otherHandler, ok := other.(*namedHandler)
	require.True(t, ok)
	assert.Equal(t, "test", otherHandler.name)
	assert.Equal(t, "group", otherHandler.group)
	assert.Equal(t, []slog.Attr(nil), otherHandler.attrs)
}
