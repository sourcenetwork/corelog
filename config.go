package corelog

import (
	"flag"
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

// FlagSet is the set of flags used to configure the logging package.
var FlagSet = flag.NewFlagSet("corelog", flag.ContinueOnError)

func init() {
	FlagSet.String("log-level", "", "Specifies the logging level.")
	FlagSet.String("log-format", "", "Specifies the output format of the logger")
	FlagSet.Bool("log-stacktrace", false, "Enables logging error stacktraces.")
	FlagSet.Bool("log-source", false, "Enables logging the source location.")
	FlagSet.String("log-output", "", "Specifies the output path for the logger.")
	FlagSet.String("log-overrides", "", "Specifies logger specific overrides.")
}

// LoadConfig returns a config with values set from environment variables and cli flags.
func LoadConfig() Config {
	// first load the environment variables
	level := os.Getenv("LOG_LEVEL")
	output := os.Getenv("LOG_OUTPUT")
	format := os.Getenv("LOG_FORMAT")
	enableStackTrace := os.Getenv("LOG_STACKTRACE")
	enableSource := os.Getenv("LOG_SOURCE")
	overrides := os.Getenv("LOG_OVERRIDES")

	if !FlagSet.Parsed() {
		FlagSet.Parse(os.Args[1:])
	}

	// override environment variables with cli flags
	FlagSet.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "log-level":
			level = f.Value.String()
		case "log-format":
			format = f.Value.String()
		case "log-stacktrace":
			enableStackTrace = f.Value.String()
		case "log-source":
			enableSource = f.Value.String()
		case "log-output":
			output = f.Value.String()
		case "log-overrides":
			overrides = f.Value.String()
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
