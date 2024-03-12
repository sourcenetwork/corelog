package corelog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigWithEnv(t *testing.T) {
	os.Setenv("LOG_LEVEL", LevelError)
	os.Setenv("LOG_OUTPUT", OutputStdout)
	os.Setenv("LOG_FORMAT", FormatJSON)
	os.Setenv("LOG_SOURCE", "true")
	os.Setenv("LOG_STACKTRACE", "true")
	t.Cleanup(os.Clearenv)

	cfg := DefaultConfig()
	assert.Equal(t, LevelError, cfg.Level)
	assert.Equal(t, OutputStdout, cfg.Output)
	assert.Equal(t, FormatJSON, cfg.Format)
	assert.Equal(t, true, cfg.EnableStackTrace)
	assert.Equal(t, true, cfg.EnableSource)
}

func TestSetConfigOverrides(t *testing.T) {
	SetConfigOverrides("net,level=debug,source=true,format=json,invalid;core,output=stdout,stacktrace=true")

	cfg := GetConfig("")
	assert.Equal(t, "", cfg.Level)
	assert.Equal(t, "", cfg.Output)
	assert.Equal(t, "", cfg.Format)
	assert.Equal(t, false, cfg.EnableStackTrace)
	assert.Equal(t, false, cfg.EnableSource)

	net := GetConfig("net")
	assert.Equal(t, LevelDebug, net.Level)
	assert.Equal(t, "", net.Output)
	assert.Equal(t, FormatJSON, net.Format)
	assert.Equal(t, false, net.EnableStackTrace)
	assert.Equal(t, true, net.EnableSource)

	core := GetConfig("core")
	assert.Equal(t, "", core.Level)
	assert.Equal(t, OutputStdout, core.Output)
	assert.Equal(t, "", core.Format)
	assert.Equal(t, true, core.EnableStackTrace)
	assert.Equal(t, false, core.EnableSource)
}
