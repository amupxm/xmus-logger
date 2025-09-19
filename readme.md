# XmusLogger

[![Go Report Card](https://goreportcard.com/badge/github.com/amupxm/xmus-logger)](https://goreportcard.com/report/github.com/amupxm/xmus-logger)
[![codecov](https://codecov.io/gh/amupxm/xmus-logger/branch/main/graph/badge.svg)](https://codecov.io/gh/amupxm/xmus-logger)
[![Go Reference](https://pkg.go.dev/badge/github.com/amupxm/xmus-logger.svg)](https://pkg.go.dev/github.com/amupxm/xmus-logger)
[![GitHub release](https://img.shields.io/github/release/amupxm/xmus-logger.svg)](https://github.com/amupxm/xmus-logger/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Test Coverage](https://github.com/amupxm/xmus-logger/actions/workflows/ci.yml/badge.svg)](https://github.com/amupxm/xmus-logger/actions/workflows/ci.yml)

A high-performance, structured logging library for Go that combines the best of both worlds: **zerolog-style chaining** for structured logging and **standard library compatibility** for easy migration.

## ✨ Features

- 🚀 **Zero Allocation** - Optimized for performance with object pooling
- 🔗 **Fluent API** - Zerolog-style method chaining for structured logging
- 🔄 **Standard Library Compatible** - Drop-in replacement for Go's `log` package
- 🌐 **Remote Logging** - Built-in HTTP remote writer support with async capabilities
- 🎯 **Context Support** - Pre-serialized context fields for efficient logging
- 📝 **JSON Output** - Structured JSON logging with proper escaping
- 🔒 **Thread Safe** - Concurrent logging without data races
- 📊 **Level Filtering** - Runtime log level control with zero-cost disabled levels
- 🎨 **Flexible Configuration** - Multiple output writers and configuration options

## 📦 Installation

```bash
go get github.com/amupxm/xmus-logger
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    xmuslogger "github.com/amupxm/xmus-logger"
)

func main() {
    // Create a new logger
    logger := xmuslogger.New()
    
    // Simple message
    logger.Info().Msg("Application started")
    
    // Structured logging with fields
    logger.Info().
        Str("service", "api").
        Int("port", 8080).
        Bool("tls", true).
        Msg("Server starting")
    
    // Error logging
    err := fmt.Errorf("database connection failed")
    logger.Error().Err(err).Str("host", "localhost").Msg("Connection error")
}
```

**Output:**
```json
{"message":"Application started","time":"2023-12-07T10:30:00Z","level":"info"}
{"service":"api","port":8080,"tls":true,"message":"Server starting","time":"2023-12-07T10:30:00Z","level":"info"}
{"error":"database connection failed","host":"localhost","message":"Connection error","time":"2023-12-07T10:30:00Z","level":"error"}
```

### Standard Library Compatibility

```go
package main

import (
    "log"
    xmuslogger "github.com/amupxm/xmus-logger"
)

func main() {
    // Replace standard logger
    logger := xmuslogger.New()
    
    // Use like standard log package
    logger.Print("This works like log.Print")
    logger.Printf("User %s logged in", "john")
    logger.Println("Application ready")
    
    // Or use it with existing code
    anyLibraryFunction(logger) // Pass to functions expecting *log.Logger
}

func anyLibraryFunction(l *xmuslogger.Logger) {
    l.Println("This library doesn't know it's using structured logging!")
}
```

**Output:**
```json
{"message":"This works like log.Print","time":"2023-12-07T10:30:00Z","level":"info","source":"stdlib"}
{"message":"User john logged in","time":"2023-12-07T10:30:00Z","level":"info","source":"stdlib"}
{"message":"Application ready","time":"2023-12-07T10:30:00Z","level":"info","source":"stdlib"}
```

## 📋 Usage Examples

### Log Levels

```go
logger := xmuslogger.New().Level(xmuslogger.DebugLevel)

logger.Trace().Msg("Very detailed information")
logger.Debug().Msg("Debug information")
logger.Info().Msg("General information")
logger.Warn().Msg("Warning message")
logger.Error().Msg("Error occurred")
// logger.Fatal().Msg("Fatal error") // Calls os.Exit(1)
```

### Context Logging

```go
// Create a context logger with pre-defined fields
ctxLogger := logger.With().
    Str("service", "user-service").
    Str("version", "v1.2.3").
    Int("instance", 1).
    Logger()

// All logs from this logger will include the context
ctxLogger.Info().Str("user_id", "123").Msg("User login")
ctxLogger.Error().Str("user_id", "123").Msg("Login failed")

// Create nested context
requestLogger := ctxLogger.With().
    Str("request_id", "req-456").
    Str("ip", "192.168.1.1").
    Logger()

requestLogger.Info().Msg("Processing request")
```

**Output:**
```json
{"service":"user-service","version":"v1.2.3","instance":1,"user_id":"123","message":"User login","time":"2023-12-07T10:30:00Z","level":"info"}
{"service":"user-service","version":"v1.2.3","instance":1,"request_id":"req-456","ip":"192.168.1.1","message":"Processing request","time":"2023-12-07T10:30:00Z","level":"info"}
```

### Different Field Types

```go
logger.Info().
    Str("string_field", "value").
    Int("int_field", 42).
    Int64("int64_field", 1234567890).
    Bool("bool_field", true).
    Err(fmt.Errorf("example error")).
    Msg("Demonstrating field types")
```

### Custom Output

```go
// Write to file
file, _ := os.Create("app.log")
defer file.Close()
logger := xmuslogger.NewWithOutput(file)

// Write to buffer
var buf bytes.Buffer
logger = xmuslogger.New().Output(&buf)

// Multiple outputs (manual setup)
logger = xmuslogger.New()
logger.SetOutput(io.MultiWriter(os.Stdout, file))
```

### Remote Logging

```go
// HTTP remote logging
logger := xmuslogger.New().RemoteHTTP(
    "https://logs.example.com/api/v1/logs",
    xmuslogger.WithHTTPAuth("your-api-token"),
    xmuslogger.WithHTTPHeaders(map[string]string{
        "Content-Type": "application/json",
        "X-Service":    "my-app",
    }),
)

// Logs will be sent to both stdout and the remote endpoint
logger.Info().Msg("This goes to both local and remote")
```

### Custom Remote Writer

```go
type CustomRemoteWriter struct {
    // Your implementation
}

func (w *CustomRemoteWriter) Write(data []byte) error {
    // Send to your logging service
    return nil
}

func (w *CustomRemoteWriter) WriteAsync(data []byte) error {
    // Async send to your logging service
    return nil
}

func (w *CustomRemoteWriter) Flush() error { return nil }
func (w *CustomRemoteWriter) Close() error { return nil }

// Use custom remote writer
logger := xmuslogger.New().Remote(&CustomRemoteWriter{})
```

### Global Logger Usage

```go
// Set a global logger
xmuslogger.SetDefault(
    xmuslogger.New().
        Level(xmuslogger.InfoLevel).
        RemoteHTTP("https://logs.example.com/api/v1/logs"),
)

// Use global functions anywhere in your app
xmuslogger.Print("Using global logger")
xmuslogger.Printf("User %s performed action", "alice")

// Get the default logger
logger := xmuslogger.Default()
logger.Info().Msg("Using default logger instance")
```

### Performance Optimizations

```go
// Disabled log levels have zero cost
logger := xmuslogger.New().Level(xmuslogger.WarnLevel)

// These calls have zero allocation and return immediately
logger.Debug().Str("expensive", "operation").Msg("Debug info") // No-op
logger.Info().Str("another", "field").Msg("Info message")     // No-op

// Only these will be processed
logger.Warn().Msg("Warning message")  // Processed
logger.Error().Msg("Error message")   // Processed
```

## 🎛️ Configuration

### Log Levels

```go
const (
    TraceLevel xmuslogger.Level = iota
    DebugLevel
    InfoLevel   // Default
    WarnLevel
    ErrorLevel
    FatalLevel
)
```

### Configuration Methods

```go
logger := xmuslogger.New().
    Level(xmuslogger.DebugLevel).                    // Set log level
    Output(file).                                   // Set output writer
    Remote(customRemoteWriter).                     // Set remote writer
    RemoteHTTP("https://api.example.com/logs")      // Set HTTP remote
```

## 🔧 Advanced Features

### Context Isolation

```go
baseLogger := xmuslogger.New()

// Each context logger is independent
serviceA := baseLogger.With().Str("service", "A").Logger()
serviceB := baseLogger.With().Str("service", "B").Logger()

serviceA.Info().Msg("Service A message") // Only has service: "A"
serviceB.Info().Msg("Service B message") // Only has service: "B"
baseLogger.Info().Msg("Base message")    // No service field
```

### Error Handling

```go
// Errors in remote writing don't affect local logging
logger := xmuslogger.New().RemoteHTTP("https://unreachable.example.com")

// This will log locally even if remote fails
logger.Error().Msg("This message is guaranteed to be logged locally")
```

### Thread Safety

```go
// XmusLogger is fully thread-safe
logger := xmuslogger.New()

// Safe to use across goroutines
go func() {
    logger.Info().Int("goroutine", 1).Msg("Concurrent logging")
}()

go func() {
    logger.Info().Int("goroutine", 2).Msg("Concurrent logging")
}()
```

## 📊 Performance

XmusLogger is designed for high performance:

- **Zero allocation** for disabled log levels
- **Object pooling** for event reuse
- **Pre-serialized context** for efficient structured logging
- **Lock-free** hot path for enabled levels
- **Optimized JSON serialization** with proper escaping

### Benchmarks

```
BenchmarkSimpleLogging-8        	 5000000	  245 ns/op	      48 B/op	       1 allocs/op
BenchmarkStructuredLogging-8    	 3000000	  456 ns/op	     128 B/op	       2 allocs/op
BenchmarkDisabledLogging-8      	100000000	   12.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkContextLogging-8       	 4000000	  298 ns/op	      64 B/op	       1 allocs/op
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development

```bash
# Clone the repository
git clone https://github.com/amupxm/xmus-logger.git
cd xmus-logger

# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [zerolog](https://github.com/rs/zerolog) for the fluent API design
- Compatible with Go's standard `log` package for easy migration
- Built for high-performance applications requiring structured logging

## 📚 Documentation

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/amupxm/xmus-logger).

---

**Made with ❤️ for the Go community**