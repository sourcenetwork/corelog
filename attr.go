package corelog

import (
	"log/slog"
	"time"
)

// Any returns an slog.Attr for the supplied value.
func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

// Bool returns an slog.Attr for a bool.
func Bool(key string, value bool) slog.Attr {
	return slog.Bool(key, value)
}

// Duration returns an slog.Attr for a time.Duration.
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Float64 returns an slog.Attr for a floating-point number.
func Float64(key string, value float64) slog.Attr {
	return slog.Float64(key, value)
}

// Group returns an slog.Attr for a Group Value.
func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

// Int converts an int to an int64 and returns
// an slog.Attr with that value.
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Int64 returns an slog.Attr for an int64.
func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

// String returns an slog.Attr for a string value.
func String(key string, value string) slog.Attr {
	return slog.String(key, value)
}

// Time returns an slog.Attr for a time.Time.
func Time(key string, value time.Time) slog.Attr {
	return slog.Time(key, value)
}

// Uint64 returns an slog.Attr for a uint64.
func Uint64(key string, value uint64) slog.Attr {
	return slog.Uint64(key, value)
}
