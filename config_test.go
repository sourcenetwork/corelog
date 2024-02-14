package corelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigOverrides(t *testing.T) {
	value := "net,level=debug,source=true,invalid;core,output=stdout,stacktrace=true"
	overrides := parseConfigOverrides(value)

	net, ok := overrides["net"]
	require.True(t, ok)

	assert.Equal(t, LevelDebug, net.Level)
	assert.Equal(t, "", net.Output)
	assert.Equal(t, "", net.Format)
	assert.Equal(t, false, net.EnableStackTrace)
	assert.Equal(t, true, net.EnableSource)

	core, ok := overrides["core"]
	require.True(t, ok)

	assert.Equal(t, "", core.Level)
	assert.Equal(t, OutputStdOut, core.Output)
	assert.Equal(t, "", core.Format)
	assert.Equal(t, true, core.EnableStackTrace)
	assert.Equal(t, false, core.EnableSource)
}
