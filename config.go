package corelog

import (
	"os"
	"strconv"
	"strings"
)

const (
	// LevelDebug specifies debug log level.
	LevelDebug = "debug"
	// LevelDebug specifies info log level.
	LevelInfo = "info"
	// LevelDebug specifies error log level.
	LevelError = "error"
	// LevelDebug specifies fatal log level.
	LevelFatal = "fatal"
	// FormatText specifies text output for a logger.
	FormatText = "text"
	// FormatText specifies json output for a logger.
	FormatJSON = "json"
	// OutputStdOut specifies stdout output for a logger.
	OutputStdOut = "stdout"
	// OutputStdErr specifies stderr output for a logger.
	OutputStdErr = "stderr"
)

// Config contains general settings for a logger.
type Config struct {
	// Level specifies the logging level.
	//
	// This value is set by the "CLOG_LEVEL" environment variable.
	Level string
	// Format specifies the output format of the logger.
	//
	// This value is set by the "CLOG_FORMAT" environment variable.
	Format string
	// EnableStackTrace enables logging error stack traces.
	//
	// This value is set by the "CLOG_STACKTRACE" environment variable.
	EnableStackTrace bool
	// EnableSource enables logging the source location.
	//
	// This value is set by the "CLOG_SOURCE" environment variable.
	EnableSource bool
	// Output specifies the output path for the logger.
	//
	// This value is set by the "CLOG_OUTPUT" environment variable.
	Output string
	// Overrides is a mapping of logger names to override configs.
	//
	// This value is set by the "CLOG_OVERRIDES" environment variable.
	Overrides map[string]Config
}

// defaultConfig is the package level config that all loggers use by default.
var defaultConfig = Config{
	Level:            LevelInfo,
	Output:           OutputStdErr,
	Format:           FormatJSON,
	EnableStackTrace: false,
	EnableSource:     false,
	Overrides:        make(map[string]Config),
}

// load the configuration from env before
// any loggers are created
func init() {
	values := make(map[string]string)
	values["level"] = os.Getenv("CLOG_LEVEL")
	values["output"] = os.Getenv("CLOG_OUTPUT")
	values["format"] = os.Getenv("CLOG_FORMAT")
	values["stacktrace"] = os.Getenv("CLOG_STACKTRACE")
	values["source"] = os.Getenv("CLOG_SOURCE")

	// parse and set default config
	defaultConfig = parseConfigMap(values)
	defaultConfig.Overrides = parseConfigOverrides(os.Getenv("CLOG_OVERRIDES"))
}

// parseConfigMap parses a map of strings into a config.
func parseConfigMap(values map[string]string) Config {
	enableSource, _ := strconv.ParseBool(values["source"])
	enableStacktrace, _ := strconv.ParseBool(values["stacktrace"])

	return Config{
		Level:            strings.ToLower(values["level"]),
		Output:           strings.ToLower(values["output"]),
		Format:           strings.ToLower(values["format"]),
		EnableStackTrace: enableStacktrace,
		EnableSource:     enableSource,
		Overrides:        make(map[string]Config),
	}
}

// parseConfigOverrides parses a mapping of config overrides from the given text.
//
// Overrides are separated by ";", and override values are comma separated,
// where the first value is the name, and the remaining values are key value
// pairs separated by "=".
func parseConfigOverrides(text string) map[string]Config {
	overrides := make(map[string]Config)
	// overrides are separated by ";"
	for _, part := range strings.Split(text, ";") {
		configMap := make(map[string]string)
		// values are separated by ","
		values := strings.Split(part, ",")
		for _, kv := range values[1:] {
			// key values are separated by "="
			values := strings.SplitN(kv, "=", 2)
			if len(values) != 2 {
				continue // invalid key value
			}
			configMap[values[0]] = values[1]
		}
		// first value is the override name
		overrideName := values[0]
		overrides[overrideName] = parseConfigMap(configMap)
	}
	return overrides
}
