package corelog

import (
	"flag"
	"os"
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

func TestLoadConfigFromEnv(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_FORMAT", "json")
	os.Setenv("LOG_STACKTRACE", "true")
	os.Setenv("LOG_SOURCE", "false")
	os.Setenv("LOG_OUTPUT", "stdout")
	os.Setenv("LOG_OVERRIDES", "net,source=true,level=error")
	t.Cleanup(os.Clearenv)

	config := LoadConfig()
	assert.Equal(t, "info", config.Level)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, true, config.EnableStackTrace)
	assert.Equal(t, false, config.EnableSource)
	assert.Equal(t, "stdout", config.Output)

	net, ok := config.Overrides["net"]
	require.True(t, ok)
	assert.Equal(t, true, net.EnableSource)
	assert.Equal(t, "error", net.Level)
}

func TestSetConfigFromFlags(t *testing.T) {
	err := flag.Set("log-level", "info")
	require.NoError(t, err)

	err = flag.Set("log-format", "json")
	require.NoError(t, err)

	err = flag.Set("log-stacktrace", "true")
	require.NoError(t, err)

	err = flag.Set("log-source", "false")
	require.NoError(t, err)

	err = flag.Set("log-output", "stdout")
	require.NoError(t, err)

	err = flag.Set("log-overrides", "net,source=true,level=error")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := flag.Set("log-level", "")
		require.NoError(t, err)

		err = flag.Set("log-format", "")
		require.NoError(t, err)

		err = flag.Set("log-stacktrace", "false")
		require.NoError(t, err)

		err = flag.Set("log-source", "false")
		require.NoError(t, err)

		err = flag.Set("log-output", "")
		require.NoError(t, err)

		err = flag.Set("log-overrides", "")
		require.NoError(t, err)
	})

	config := LoadConfig()
	assert.Equal(t, "info", config.Level)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, true, config.EnableStackTrace)
	assert.Equal(t, false, config.EnableSource)
	assert.Equal(t, "stdout", config.Output)

	net, ok := config.Overrides["net"]
	require.True(t, ok)
	assert.Equal(t, true, net.EnableSource)
	assert.Equal(t, "error", net.Level)
}
