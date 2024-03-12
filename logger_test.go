package corelog

import (
	"context"
	"errors"
	"fmt"
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

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test")
	assert.Equal(t, "test", logger.name)

	_, ok := logger.stderrHandler.(*slog.TextHandler)
	assert.True(t, ok)

	_, ok = logger.stdoutHandler.(*slog.TextHandler)
	assert.True(t, ok)
}

func TestNewLoggerWithOutputStderr(t *testing.T) {
	SetConfig(Config{Output: OutputStderr})

	stderrHandler := &TestHandler{level: levelDebug}
	stdoutHandler := &TestHandler{level: levelDebug}

	logger := NewLogger("test")
	logger.stderrHandler = stderrHandler
	logger.stdoutHandler = stdoutHandler

	logger.Debug("test")
	assert.Len(t, stdoutHandler.records, 0)
	assert.Len(t, stderrHandler.records, 1)
}

func TestNewLoggerWithOutputStdout(t *testing.T) {
	SetConfig(Config{Output: OutputStdout})

	stderrHandler := &TestHandler{level: levelDebug}
	stdoutHandler := &TestHandler{level: levelDebug}

	logger := NewLogger("test")
	logger.stderrHandler = stderrHandler
	logger.stdoutHandler = stdoutHandler

	logger.Debug("test")
	assert.Len(t, stderrHandler.records, 0)
	assert.Len(t, stdoutHandler.records, 1)
}

func TestNewLoggerWithFormatText(t *testing.T) {
	SetConfig(Config{Format: FormatText})
	logger := NewLogger("test")

	_, ok := logger.stderrHandler.(*slog.TextHandler)
	assert.True(t, ok)

	_, ok = logger.stdoutHandler.(*slog.TextHandler)
	assert.True(t, ok)
}

func TestNewLoggerWithFormatJSON(t *testing.T) {
	SetConfig(Config{Format: FormatJSON})
	logger := NewLogger("test")

	_, ok := logger.stderrHandler.(*slog.JSONHandler)
	assert.True(t, ok)

	_, ok = logger.stdoutHandler.(*slog.JSONHandler)
	assert.True(t, ok)
}

func TestNewLoggerWithLevelDebug(t *testing.T) {
	SetConfig(Config{Level: LevelDebug})
	logger := NewLogger("test")

	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelDebug))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelInfo))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelError))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelFatal))
}

func TestNewLoggerWithLevelInfo(t *testing.T) {
	SetConfig(Config{Level: LevelInfo})
	logger := NewLogger("test")

	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelDebug))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelInfo))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelError))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelFatal))
}

func TestNewLoggerWithLevelError(t *testing.T) {
	SetConfig(Config{Level: LevelError})
	logger := NewLogger("test")

	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelDebug))
	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelInfo))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelError))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelFatal))
}

func TestNewLoggerWithLevelFatal(t *testing.T) {
	SetConfig(Config{Level: LevelFatal})
	logger := NewLogger("test")

	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelDebug))
	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelInfo))
	assert.False(t, logger.stderrHandler.Enabled(context.Background(), levelError))
	assert.True(t, logger.stderrHandler.Enabled(context.Background(), levelFatal))
}

func TestLoggerLogWithConfigOverride(t *testing.T) {
	SetConfig(Config{
		Level:            LevelFatal,
		Format:           FormatJSON,
		Output:           OutputStderr,
		EnableStackTrace: false,
		EnableSource:     false,
	})

	SetConfigOverride("test", Config{
		Level:            LevelInfo,
		Format:           FormatText,
		Output:           OutputStdout,
		EnableStackTrace: true,
		EnableSource:     true,
	})

	stderrHandler := &TestHandler{}
	stdoutHandler := &TestHandler{}

	logger := NewLogger("test")
	logger.stderrHandler = stderrHandler
	logger.stdoutHandler = stdoutHandler

	ctx := context.Background()
	err := errors.New("test error")

	logger.log(ctx, levelInfo, err, "test", []any{slog.Any("arg1", "val1")})
	assert.Len(t, stderrHandler.records, 0)
	require.Len(t, stdoutHandler.records, 1)

	assert.NotEqual(t, uintptr(0x00), stdoutHandler.records[0].PC)
	assert.Equal(t, levelInfo, stdoutHandler.records[0].Level)
	assert.Equal(t, "test", stdoutHandler.records[0].Message)

	attrs := []slog.Attr{
		slog.Any("name", "test"),
		slog.Any("arg1", "val1"),
		slog.Any("stacktrace", fmt.Sprintf("%+v", err)),
	}
	assertRecordAttrs(t, stdoutHandler.records[0], attrs...)
}

func TestLoggerInfoWithEnableSource(t *testing.T) {
	SetConfig(Config{EnableSource: true})

	handler := &TestHandler{}
	logger := &Logger{
		name:          "test",
		stderrHandler: handler,
		stdoutHandler: handler,
	}

	logger.Info("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assert.NotEqual(t, uintptr(0x00), handler.records[0].PC)
	assertRecordAttrs(t, handler.records[0], slog.Any("name", "test"), slog.Any("arg1", "val1"))
}

func TestLoggerWithAttrs(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		stderrHandler: handler,
		stdoutHandler: handler,
	}

	attrs := []slog.Attr{slog.Any("extra", "value")}
	other := logger.WithAttrs(attrs...)

	stderrHandler, ok := other.stderrHandler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, attrs, stderrHandler.attrs)

	stdoutHandler, ok := other.stdoutHandler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, attrs, stdoutHandler.attrs)
}

func TestLoggerWithGroup(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		stderrHandler: handler,
		stdoutHandler: handler,
	}

	other := logger.WithGroup("group")

	stderrHandler, ok := other.stderrHandler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, "group", stderrHandler.group)

	stdoutHandler, ok := other.stdoutHandler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, "group", stdoutHandler.group)
}

// assertRecordAttrs asserts that the record has matching attributes.
func assertRecordAttrs(
	t *testing.T,
	record slog.Record,
	expected ...slog.Attr,
) {
	var actual []slog.Attr
	record.Attrs(func(a slog.Attr) bool {
		actual = append(actual, a)
		return true
	})
	assert.Equal(t, expected, actual)
}
