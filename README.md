# Corelog

Logging library used across Source Network projects.

## Configuration

Loggers are configured via a common set of environment variables:

- `CLOG_LEVEL` sets logging level for all loggers
    -  Can be one of `info` `debug` `error` `fatal`
- `CLOG_FORMAT` sets logging format for all loggers
    - Can be one of `json` `text`
- `CLOG_STACKTRACE` enables stacktraces for all loggers
    - Can be one of `true` `false`
- `CLOG_SOURCE` enables source location for all loggers
    - Can be one of `true` `false`
- `CLOG_OUTPUT` sets the output path for all loggers
    - Can be one of `stderr` `stdout`
- `CLOG_OVERRIDES` sets logger specific overrides
    - Format `net,level=info;core,output=stdout`