package corelog

import "log/slog"

const (
	levelDebug = slog.LevelDebug
	levelInfo  = slog.LevelInfo
	levelError = slog.LevelError
	levelFatal = slog.Level(12)
)

// levelLabels is a mapping of log levels to custom labels.
var levelLabels = map[slog.Level]string{
	levelDebug: "DEBUG",
	levelInfo:  "INFO",
	levelError: "ERROR",
	levelFatal: "FATAL",
}

// replaceLoggerLevel returns a function that replaces log level attributes with a corrected label.
func replaceLoggerLevel(level slog.Level) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key != slog.LevelKey {
			return a
		}
		if label, ok := levelLabels[level]; ok {
			a.Value = slog.StringValue(label)
		}
		return a
	}
}
