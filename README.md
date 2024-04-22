# Corelog

Logging library used across Source Network projects.

## Usage

```go
package main

import "github.com/sourcenetwork/corelog"

var log = corelog.NewLogger("main")

func main() {
    // with attributes
    log.Info("message", corelog.String("key", "val"))

    // with context
    log.InfoContext(ctx, "message", corelog.Int("key", 10))

    // with error stacktrace
    log.ErrorE("message", err, corelog.Bool("key", true))

    // with common attributes
    attrs := log.WithAttrs(corelog.Float64("key", float64(1.234)))
    attrs.Info("message")

    // with common group
    group := log.WithGroup("group")
    group.Info("message", corelog.Any("key", struct{}{}))
}
```

## Configuration

Default config values can be set via environment variables.

| Env              | Description               | Values                              |
| ---------------- | ------------------------- | ----------------------------------- |
| `LOG_LEVEL`      | sets logging level        | `info` `error`                      |
| `LOG_FORMAT`     | sets logging format       | `json` `text`                       |
| `LOG_STACKTRACE` | enables stacktraces       | `true` `false`                      |
| `LOG_SOURCE`     | enables source location   | `true` `false`                      |
| `LOG_OUTPUT`     | sets the output path      | `stderr` `stdout`                   |
| `LOG_OVERRIDES`  | logger specific overrides | `net,level=info;core,output=stdout` |
| `LOG_NO_COLOR`   | disable color text output | `true` `false`                      |
