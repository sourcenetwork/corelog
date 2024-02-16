# Corelog

Logging library used across Source Network projects.

## Usage

```go
package main

import (
    "log/slog"

    "github.com/sourcenetwork/corelog"
)

var log = corelog.NewLogger("main")

func main() {
    log.Debug("doing stuff", "key", "val")

    attrs := log.WithAttrs(slog.Any("key", "val"))
    attrs.Info("doing stuff with attrs")

    group := log.WithGroup("group")
    group.Info("doing stuff with group", "key", "val")
}
```

### Usage with Cobra

Add config flags to Cobra.

```go
package main

import (
    "github.com/sourcenetwork/corelog"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{...}
    cmd.PersistentFlags().AddGoFlagSet(corelog.FlagSet)
}
```

## Configuration

Loggers are configured via environment variables and/or command line flags.

| Env              | Flag             | Description               | Values                              |
| ---------------- | ---------------- | ------------------------- | ----------------------------------- |
| `LOG_LEVEL`      | `log-level`      | sets logging level        | `info` `debug` `error` `fatal`      |
| `LOG_FORMAT`     | `log-format`     | sets logging format       | `json` `text`                       |
| `LOG_STACKTRACE` | `log-stacktrace` | enables stacktraces       | `true` `false`                      |
| `LOG_SOURCE`     | `log-source`     | enables source location   | `true` `false`                      |
| `LOG_OUTPUT`     | `log-output`     | sets the output path      | `stderr` `stdout`                   |
| `LOG_OVERRIDES`  | `log-overrides`  | logger specific overrides | `net,level=info;core,output=stdout` |
