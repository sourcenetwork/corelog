package corelog

import (
	"context"
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
}

var _ (slog.Handler) = (*TestHandler)(nil)

func (h *TestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
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

	_, ok := logger.handler.(*slog.TextHandler)
	assert.True(t, ok)

	assert.False(t, logger.enableSource)
	assert.False(t, logger.enableStackTrace)
}

func TestNewLoggerWithConfigOverrides(t *testing.T) {
	config := Config{
		EnableStackTrace: false,
		EnableSource:     false,
		Format:           FormatJSON,
		Overrides:        make(map[string]Config),
	}
	config.Overrides["test"] = Config{
		EnableStackTrace: true,
		EnableSource:     true,
		Format:           FormatText,
	}

	logger := NewLoggerWithConfig("test", config)

	_, ok := logger.handler.(*slog.TextHandler)
	assert.True(t, ok)

	assert.True(t, logger.enableSource)
	assert.True(t, logger.enableStackTrace)
}

func TestLoggerInfo(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.Info("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerDebug(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.Debug("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelDebug, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerError(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.Error("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelError, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerErrorE(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:          handler,
		enableStackTrace: true,
	}

	err := fmt.Errorf("this is a test error")
	logger.ErrorE("test", err, "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelError, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"), slog.Any("stacktrace", fmt.Sprintf("%+v", err)))
}

func TestLoggerFatal(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:  handler,
		skipExit: true,
	}

	logger.Fatal("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelFatal, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerFatalE(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:          handler,
		enableStackTrace: true,
		skipExit:         true,
	}

	err := fmt.Errorf("this is a test error")
	logger.FatalE("test", err, "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelFatal, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"), slog.Any("stacktrace", fmt.Sprintf("%+v", err)))
}

func TestLoggerInfoContext(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.InfoContext(context.Background(), "test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerDebugContext(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.DebugContext(context.Background(), "test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelDebug, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerErrorContext(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{handler: handler}

	logger.ErrorContext(context.Background(), "test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelError, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerErrorContextE(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:          handler,
		enableStackTrace: true,
	}

	err := fmt.Errorf("this is a test error")
	logger.ErrorContextE(context.Background(), "test", err, "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelError, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"), slog.Any("stacktrace", fmt.Sprintf("%+v", err)))
}

func TestLoggerFatalContext(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:  handler,
		skipExit: true,
	}

	logger.FatalContext(context.Background(), "test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelFatal, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
}

func TestLoggerFatalContextE(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:          handler,
		enableStackTrace: true,
		skipExit:         true,
	}

	err := fmt.Errorf("this is a test error")
	logger.FatalContextE(context.Background(), "test", err, "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelFatal, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"), slog.Any("stacktrace", fmt.Sprintf("%+v", err)))
}

func TestLoggerInfoWithEnableSource(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler:      handler,
		enableSource: true,
	}

	logger.Info("test", "arg1", "val1")
	require.Len(t, handler.records, 1)

	assert.Equal(t, levelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assert.NotEqual(t, uintptr(0x00), handler.records[0].PC)
	assertRecordAttrs(t, handler.records[0], slog.Any("arg1", "val1"))
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
