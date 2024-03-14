# Corelog

Logging library used across Source Network projects.

## Usage

```go
package main

import "github.com/sourcenetwork/corelog"

var log = corelog.NewLogger("main")

func main() {
    // with alternating key value pairs
    log.Debug("message", "key", "val")
    
    // with explicit attributes
    log.Debug("message", corelog.String("key", "val"))

    // with common attributes
    attrs := log.WithAttrs(corelog.String("key", "val"))
    attrs.Info("message")

    // with common group
    group := log.WithGroup("group")
    group.Info("message", "key", "val")
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
