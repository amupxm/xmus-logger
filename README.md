# XMUS-LOGGER


pure golang logger compatible with golang io standards.

## USAGE :

```go
	logOptions := logger.LoggerOptions{
		LogLevel: 6,    // read more about log levels in logLevel section
		Verbose:  true, // if true, print more detailed log messages and benchmark
		Std:      true, // if true, print log messages to stdout
	}
	log := logger.CreateLogger(&logOptions)
	log.LogF("📑 %s \n", "your first log")

	// Trace function caller
	log.GetCaller().Alertln(" called me ")

	// Or Trace call stack
	log.Alertln(" called me ").TraceStack()

	// Or Use prefix :
	prefixLogger := log.Prefix("Prefix", "log")
	prefixLogger.AlertF("📑%s\n", "your first log with prefix")

	prefixLogger.End()
	log.End()
```

## AVAILABLE METHODS :


```go

		// Log logs a message at log level
		Logln(a ...interface{}) LogResult
		// Logln logs a message at log level to new line
		Log(a ...interface{}) LogResult
		// LogF logs a message at log level with string formater
		LogF(format string, a ...interface{}) LogResult

		// Alert logs a message at log level
		Alertln(a ...interface{}) LogResult
		// Alertln logs a message at log level to new line
		Alert(a ...interface{}) LogResult
		// AlertF logs a message at log level with string formater
		AlertF(format string, a ...interface{}) LogResult

		// Error logs a message at log level
		Error(a ...interface{}) LogResult
		// Errorln logs a message at log level to new line
		Errorln(a ...interface{}) LogResult
		// ErrorF logs a message at log level with string formater
		ErrorF(format string, a ...interface{}) LogResult

		// Highlight logs a message at log level
		Highlight(a ...interface{}) LogResult
		// Highlightln logs a message at log level to new line
		Highlightln(a ...interface{}) LogResult
		// HighlightF logs a message at log level with string formater
		HighlightF(format string, a ...interface{}) LogResult

		// Inform logs a message at log level
		Inform(a ...interface{}) LogResult
		// Informln logs a message at log level to new line
		Informln(a ...interface{}) LogResult
		// InformF logs a message at log level with string formater
		InformF(format string, a ...interface{}) LogResult

		// Trace logs a message at log level
		Trace(a ...interface{}) LogResult
		// Traceln logs a message at log level to new line
		Traceln(a ...interface{}) LogResult
		// TraceF logs a message at log level with string formater
		TraceF(format string, a ...interface{}) LogResult

		// Warn logs a message at log level
		Warn(a ...interface{}) LogResult
		// Warnln logs a message at log level to new line
		Warnln(a ...interface{}) LogResult
		// WarnF logs a message at log level with string formater
		WarnF(format string, a ...interface{}) LogResult
```