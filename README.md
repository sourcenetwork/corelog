# Corelog

Logging library used across Source Network projects.

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
