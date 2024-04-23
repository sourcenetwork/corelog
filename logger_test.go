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

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test")
	assert.Equal(t, "test", logger.name)
}

func TestLoggerLogWithConfigOverride(t *testing.T) {
	SetConfig(Config{
		Level:            LevelError,
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

	handler := &TestHandler{}

	logger := NewLogger("test")
	logger.handler = handler

	ctx := context.Background()
	err := errors.New("test error")

	logger.log(ctx, slog.LevelInfo, err, "test", []slog.Attr{slog.Any("arg1", "val1")})
	require.Len(t, handler.records, 1)

	assert.NotEqual(t, uintptr(0x00), handler.records[0].PC)
	assert.Equal(t, slog.LevelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)

	attrs := []slog.Attr{
		slog.Any(nameKey, "test"),
		slog.Any("arg1", "val1"),
		slog.Any(stackKey, fmt.Sprintf("%+v", err)),
	}
	assertRecordAttrs(t, handler.records[0], attrs...)
}

func TestLoggerInfoWithEnableSource(t *testing.T) {
	SetConfig(Config{EnableSource: true})

	handler := &TestHandler{}
	logger := &Logger{
		name:    "test",
		handler: handler,
	}

	logger.Info("test", String("arg1", "val1"))
	require.Len(t, handler.records, 1)

	assert.Equal(t, slog.LevelInfo, handler.records[0].Level)
	assert.Equal(t, "test", handler.records[0].Message)
	assert.NotEqual(t, uintptr(0x00), handler.records[0].PC)

	attrs := []slog.Attr{
		slog.Any(nameKey, "test"),
		slog.Any("arg1", "val1"),
	}
	assertRecordAttrs(t, handler.records[0], attrs...)
}

func TestLoggerWithAttrs(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler: handler,
	}

	attrs := []slog.Attr{slog.Any("extra", "value")}
	other := logger.WithAttrs(attrs...)

	otherHandler, ok := other.handler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, attrs, otherHandler.attrs)
}

func TestLoggerWithGroup(t *testing.T) {
	handler := &TestHandler{}
	logger := &Logger{
		handler: handler,
	}

	other := logger.WithGroup("group")

	otherHandler, ok := other.handler.(*TestHandler)
	require.True(t, ok)
	assert.Equal(t, "group", otherHandler.group)
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
