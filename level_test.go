package corelog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceLevelLabel(t *testing.T) {
	for level, label := range levelLabels {
		replace := replaceLoggerLevel(level)
		attr := replace(nil, slog.Any(slog.LevelKey, level))
		assert.Equal(t, slog.String(slog.LevelKey, label), attr)
	}
}
