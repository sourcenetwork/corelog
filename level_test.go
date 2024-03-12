package corelog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedLeveler(t *testing.T) {
	leveler := namedLeveler("core")
	assert.Equal(t, levelInfo, leveler.Level())

	attr := leveler.ReplaceAttr(nil, slog.Any(slog.LevelKey, levelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "INFO"), attr)

	SetConfigOverride("core", Config{Level: LevelInfo})
	assert.Equal(t, levelInfo, leveler.Level())

	attr = leveler.ReplaceAttr(nil, slog.Any(slog.LevelKey, levelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "INFO"), attr)

	SetConfigOverride("core", Config{Level: LevelDebug})
	assert.Equal(t, levelDebug, leveler.Level())

	attr = leveler.ReplaceAttr(nil, slog.Any(slog.LevelKey, levelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "DEBUG"), attr)

	SetConfigOverride("core", Config{Level: LevelError})
	assert.Equal(t, levelError, leveler.Level())

	attr = leveler.ReplaceAttr(nil, slog.Any(slog.LevelKey, levelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "ERROR"), attr)

	SetConfigOverride("core", Config{Level: LevelFatal})
	assert.Equal(t, levelFatal, leveler.Level())

	attr = leveler.ReplaceAttr(nil, slog.Any(slog.LevelKey, levelInfo))
	assert.Equal(t, slog.String(slog.LevelKey, "FATAL"), attr)
}
