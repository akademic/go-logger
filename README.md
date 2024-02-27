## FEATURES

- Only stdlib dependencies
- Error, Info, Debug levels (Errors print always, Info and Debug only if configured)
- Logging with prefixes. Log output with particular prefix can be controlled with config with same levels (Error, Info, Debug)

## INSTALL

go get github.com/akademic/go-logger

## USAGE

### Full example

```go

package main

import (
    "github.com/akademic/go-logger"
)

func main() {
    l := logger.New("", logger.Config{
        Level: logger.LogDebug,
        ComponentLevel: map[string]logger.LogLevel{
            "db": logger.LogError,
            "api": logger.LogInfo,
        },
    })

    l.Error("number: %d", 1) // [err] number: 1
    l.Info("number: %d", 2) // [inf] number: 2
    l.Debug("number: %d", 3) // [dbg] number: 3

    dbL := l.WithPrefix("db")

    dbL.Error("number: %d", 1) // [err] [db] number: 1
    dbL.Info("number: %d", 2) // nothing because db level is set to Error
    dbL.Debug("number: %d", 3) // nothing because db level is set to Error

    apiL := l.WithPrefix("api")

    apiL.Error("number: %d", 1) // [err] [api] number: 1
    apiL.Info("number: %d", 2) // [inf] [api] number: 2
    apiL.Debug("number: %d", 3) // nothing because api level is set to Info
}

```

### Architecture example

Prefixes can be used to control log output for different components of the application.
In this example, we have a database and an API.
We want to log only errors from the database and info from the API.

```go

package main

import (
    "github.com/akademic/go-logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

    l := logger.New("", logger.Config{
        Level: logger.LogDebug,
        ComponentLevel: map[string]logger.LogLevel{
            "db": logger.LogError,
            "api": logger.LogInfo,
        },
    })

    dbL := l.WithPrefix("db")
    repository, err := NewRepository(dbL)
    if err != nil {
        l.Error("failed to create repository: %v", err)
        return
    }

    apiL := l.WithPrefix("api")
    api, err := NewAPI(apiL, repository)
    if err != nil {
        l.Error("failed to create api: %v", err)
        return
    }

    go api.Start()

    <-ctx.Done()
}
```

### Interface

This package provides a simple interface for logging.

```go
type Logger interface {
    WithPrefix(prefix string) Logger // returns a new logger with a prefix, prefix is replaced with new
    Error(format string, args ...interface{})
    Info(format string, args ...interface{})
    Debug(format string, args ...interface{})
}
```

Use it to receive a logger instance.

```go
func NewAPI(l logger.Logger) (*API, error) {
    // ...
    l.Info("api started")
    // ...
    adminApiLogger := l.WithPrefix("admin-api")

    adminApiLogger.Info("admin api started") // [inf] [admin-api] admin api started

    customerApiLogger := l.WithPrefix("customer-api")

    customerApiLogger.Info("customer api started") // [inf] [customer-api] customer api started
}
```
