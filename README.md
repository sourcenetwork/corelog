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

## Configuration

Default config values can be set via environment variables.

| Env              | Description               | Values                              |
| ---------------- | ------------------------- | ----------------------------------- |
| `LOG_LEVEL`      | sets logging level        | `info` `debug` `error` `fatal`      |
| `LOG_FORMAT`     | sets logging format       | `json` `text`                       |
| `LOG_STACKTRACE` | enables stacktraces       | `true` `false`                      |
| `LOG_SOURCE`     | enables source location   | `true` `false`                      |
| `LOG_OUTPUT`     | sets the output path      | `stderr` `stdout`                   |
| `LOG_OVERRIDES`  | logger specific overrides | `net,level=info;core,output=stdout` |
