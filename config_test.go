package corelog

import (
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
	args := os.Args
	os.Args = []string{
		"test",
		"--log-level=info",
		"--log-format=json",
		"--log-stacktrace=true",
		"--log-source=false",
		"--log-output=stdout",
		"--log-overrides=net,source=true,level=error",
	}
	t.Cleanup(func() { os.Args = args })

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
