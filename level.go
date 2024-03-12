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

// namedLeveler is an slog.Leveler that gets its value from a named config.
type namedLeveler string

func (n namedLeveler) Level() slog.Level {
	switch cfg := GetConfig(string(n)); cfg.Level {
	case LevelDebug:
		return levelDebug
	case LevelInfo:
		return levelInfo
	case LevelError:
		return levelError
	case LevelFatal:
		return levelFatal
	default:
		return levelInfo
	}
}

func (n namedLeveler) ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		a.Value = slog.StringValue(levelLabels[n.Level()])
	}
	return a
}
