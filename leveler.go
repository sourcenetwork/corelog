package corelog

import "log/slog"

// namedLeveler is an slog.Leveler that gets its value from a named config.
type namedLeveler string

func (n namedLeveler) Level() slog.Level {
	switch cfg := GetConfig(string(n)); cfg.Level {
	case LevelInfo:
		return slog.LevelInfo
	case LevelError:
		return slog.LevelError
	default:
		// default to info if no value is set
		// or the set value is invalid
		return slog.LevelInfo
	}
}
