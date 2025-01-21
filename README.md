# ESL_Logger

This module provides an interface for logging messages with different severity levels: `DEBUG`, `INFO`, `WARN`, `ERROR` and `FATAL`. Messages can be logged to standard output or to a file.

## Usage

### First, create a new instance of Logger:

```go
logger := ESL_Logger.NewLogger("DevLogger", "development")
// or
logger := ESL_Logger.NewLogger("ProdLogger", "production", "app.log")
```

You can log messages with different severity levels:
```
logger.Debug("This is a debug message")
logger.Info("This is an info message")
logger.Warn("This is a warning message")
logger.Error("This is an error message")
logger.Fatal("This is a fatal message")
```
### Customizing the Logger's Format
To customize the logger's format, modify the Format field in the Config struct when creating a new logger instance.
The placeholders in the format string are:
- `%s` for the timestamp
- `%s` for the logger name
- `%s` for the severity level
- `%s `for the actual log message

```go
config := &ESL_Logger.Config{
    Name:       "CustomLogger",
    Output:     "stdout", // or "app.log" to log to a file
    Format:     "[%s] %s - %s: %s", // Custom format: [timestamp] logger_name - severity: message
    MaxSize:    10, // maximum file size in megabytes
    MaxBackups: 3,  // maximum number of backup files
    MaxAge:     30, // maximum number of days to retain backup files
}

logger := ESL_Logger.NewLogger(config)
```

### Colored Logs
Log messages are colored based on their severity level:
- `DEBUG` messages are colored in Cyan
- `INFO` messages are colored in Blue
- `WARN` messages are colored in Yellow
- `ERROR` messages are colored in Red
- `FATAL` messages are colored in Purple

### Environment-Specific Configuration
You can specify different configurations for different environments. For example, you can log messages to a file in the production environment and to standard output in the development environment.

```go
logger := ESL_Logger.NewLogger("DevLogger", "development")
// or
logger := ESL_Logger.NewLogger("ProdLogger", "production", "app.log")
```
### Configuration
The logger can be configured using the Config struct. The following fields can be set:

```go
type Config struct {
    Name       string // Name of the logger
    Output     string // Output destination: "stdout" or "app.log"
    Format     string // Format of the log message
    MaxSize    int    // Maximum file size in megabytes
    MaxBackups int    // Maximum number of backup files
    MaxAge     int    // Maximum number of days to retain backup files
}
```