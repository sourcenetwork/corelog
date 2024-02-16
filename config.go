package corelog

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
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
	Level string
	// Format specifies the output format of the logger.
	Format string
	// EnableStackTrace enables logging error stack traces.
	EnableStackTrace bool
	// EnableSource enables logging the source location.
	EnableSource bool
	// Output specifies the output path for the logger.
	Output string
	// Overrides is a mapping of logger names to override configs.
	Overrides map[string]Config
}

var (
	// LevelFlag is a flag that sets default `Level` value.
	LevelFlag = flag.String("log-level", "", "Specifies the logging level.")
	// FormatFlag is a flag that sets the default `Format` value.
	FormatFlag = flag.String("log-format", "", "Specifies the output format of the logger")
	// EnableStackTraceFlag is a flag sets the default `EnableStackTrace` value.
	EnableStackTraceFlag = flag.Bool("log-stacktrace", false, "Enables logging error stacktraces.")
	// EnableSourceFlag is a flag that sets the default `EnableSource` value.
	EnableSourceFlag = flag.Bool("log-source", false, "Enables logging the source location.")
	// OutputFlag is a flag that sets the default `Output` value.
	OutputFlag = flag.String("log-output", "", "Specifies the output path for the logger.")
	// OverridesFlag is a flag that sets the default `Overrides` value.
	OverridesFlag = flag.String("log-overrides", "", "Specifies logger specific overrides.")
)

// parseFlagsOnce ensures that `flag.Parse` is only called once.
var parseFlagsOnce sync.Once

// LoadConfig returns a config with values set from environment variables and cli flags.
func LoadConfig() Config {
	// first load the environment variables
	level := os.Getenv("LOG_LEVEL")
	output := os.Getenv("LOG_OUTPUT")
	format := os.Getenv("LOG_FORMAT")
	enableStackTrace := os.Getenv("LOG_STACKTRACE")
	enableSource := os.Getenv("LOG_SOURCE")
	overrides := os.Getenv("LOG_OVERRIDES")

	if !flag.Parsed() {
		parseFlagsOnce.Do(flag.Parse)
	}

	// override environment variables with cli flags
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "log-level":
			level = *LevelFlag
		case "log-format":
			format = *FormatFlag
		case "log-stacktrace":
			enableStackTrace = strconv.FormatBool(*EnableStackTraceFlag)
		case "log-source":
			enableSource = strconv.FormatBool(*EnableSourceFlag)
		case "log-output":
			output = *OutputFlag
		case "log-overrides":
			overrides = *OverridesFlag
		}
	})

	values := make(map[string]string)
	values["level"] = level
	values["output"] = output
	values["format"] = format
	values["stacktrace"] = enableStackTrace
	values["source"] = enableSource

	// parse config values
	config := parseConfigMap(values)
	config.Overrides = parseConfigOverrides(overrides)
	return config
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
		if overrideName == "" {
			continue // empty logger name
		}
		overrides[overrideName] = parseConfigMap(configMap)
	}
	return overrides
}
