package corelog

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	// LevelDebug specifies info log level.
	LevelInfo = "info"
	// LevelDebug specifies error log level.
	LevelError = "error"
	// FormatText specifies text output for a logger.
	FormatText = "text"
	// FormatJSON specifies json output for a logger.
	FormatJSON = "json"
	// OutputStdout specifies stdout output for a logger.
	OutputStdout = "stdout"
	// OutputStderr specifies stderr output for a logger.
	OutputStderr = "stderr"
)

var (
	configMutex     sync.RWMutex
	configValue     Config
	configOverrides = make(map[string]Config)
)

func init() {
	SetConfig(DefaultConfig())
	SetConfigOverrides(os.Getenv("LOG_OVERRIDES"))
}

// Config contains general settings for a logger.
type Config struct {
	// Level specifies the logging level.
	Level string
	// Format specifies the output format of the logger.
	Format string
	// EnableStackTrace enables logging error stack traces.
	EnableStackTrace bool
	// EnableSource enables logging the source location.
	EnableSource bool
	// Output specifies the output path for the logger.
	Output string
	// DisableColor specifies if colored output is disabled.
	DisableColor bool
}

// DefaultConfig returns a config with default values.
//
// The default values are derived from environment variables.
func DefaultConfig() Config {
	enableSource, _ := strconv.ParseBool(os.Getenv("LOG_SOURCE"))
	enableStacktrace, _ := strconv.ParseBool(os.Getenv("LOG_STACKTRACE"))
	disableColor, _ := strconv.ParseBool(os.Getenv("LOG_NO_COLOR"))

	return Config{
		Level:            strings.ToLower(os.Getenv("LOG_LEVEL")),
		Output:           strings.ToLower(os.Getenv("LOG_OUTPUT")),
		Format:           strings.ToLower(os.Getenv("LOG_FORMAT")),
		EnableSource:     enableSource,
		EnableStackTrace: enableStacktrace,
		DisableColor:     disableColor,
	}
}

// GetConfig returns the config for a named logger.
func GetConfig(name string) Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if val, ok := configOverrides[name]; ok {
		return val
	}
	return configValue
}

// SetConfig sets the config values for all loggers.
func SetConfig(cfg Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configValue = cfg
}

// SetConfigOverride sets the config override for the given named logger.
func SetConfigOverride(name string, cfg Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configOverrides[name] = cfg
}

// SetConfigOverrides parses and sets config overrides from the given text.
//
// Overrides are separated by ";", and override values are comma separated,
// where the first value is the name, and the remaining values are key value
// pairs separated by "=".
func SetConfigOverrides(text string) {
	// overrides are separated by ";"
	for _, override := range strings.Split(text, ";") {
		// override parts are separated by ","
		parts := strings.Split(override, ",")
		// first part is the override name
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue // empty logger name
		}
		config := DefaultConfig()
		// remaining parts are key value pairs
		for _, pair := range parts[1:] {
			// key value pairs are separated by "="
			values := strings.SplitN(pair, "=", 2)
			if len(values) != 2 {
				continue // invalid key value
			}
			key := strings.TrimSpace(values[0])
			val := strings.TrimSpace(values[1])
			switch strings.ToLower(key) {
			case "level":
				config.Level = strings.ToLower(val)
			case "format":
				config.Format = strings.ToLower(val)
			case "output":
				config.Output = strings.ToLower(val)
			case "stacktrace":
				config.EnableStackTrace, _ = strconv.ParseBool(val)
			case "source":
				config.EnableSource, _ = strconv.ParseBool(val)
			case "no-color":
				config.DisableColor, _ = strconv.ParseBool(val)
			}
		}
		SetConfigOverride(name, config)
	}
}
