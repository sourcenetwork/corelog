package corelog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedLeveler(t *testing.T) {
	leveler := namedLeveler("core")
	assert.Equal(t, slog.LevelInfo, leveler.Level())

	SetConfigOverride("core", Config{Level: LevelInfo})
	assert.Equal(t, slog.LevelInfo, leveler.Level())

	SetConfigOverride("core", Config{Level: LevelError})
	assert.Equal(t, slog.LevelError, leveler.Level())
}
