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
log.Informln(" my cal stack is ").TraceStack()

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
LogF(format string, a ...interface{}) LogResul


// Alert logs a message at log level
Alertln(a ...interface{}) LogResult
// Alertln logs a message at log level to new line
Alert(a ...interface{}) LogResult
// AlertF logs a message at log level with string formater
AlertF(format string, a ...interface{}) LogResul


// Error logs a message at log level
Error(a ...interface{}) LogResult
// Errorln logs a message at log level to new line
Errorln(a ...interface{}) LogResult
// ErrorF logs a message at log level with string formater
ErrorF(format string, a ...interface{}) LogResul


// Highlight logs a message at log level
Highlight(a ...interface{}) LogResult
// Highlightln logs a message at log level to new line
Highlightln(a ...interface{}) LogResult
// HighlightF logs a message at log level with string formater
HighlightF(format string, a ...interface{}) LogResul


// Inform logs a message at log level
Inform(a ...interface{}) LogResult
// Informln logs a message at log level to new line
Informln(a ...interface{}) LogResult
// InformF logs a message at log level with string formater
InformF(format string, a ...interface{}) LogResul


// Trace logs a message at log level
Trace(a ...interface{}) LogResult
// Traceln logs a message at log level to new line
Traceln(a ...interface{}) LogResult
// TraceF logs a message at log level with string formater
TraceF(format string, a ...interface{}) LogResul


// Warn logs a message at log level
Warn(a ...interface{}) LogResult
// Warnln logs a message at log level to new line
Warnln(a ...interface{}) LogResult
// WarnF logs a message at log level with string formater
WarnF(format string, a ...interface{}) LogResult
```


## LogLevels :

|LogLevel|int (logLevelNumber)|limits|
--- | --- | ---
|Nothing|0|ban all logs|
|Alert|1|prints only alert and error|
|Error|1|prints only alert and error|
|Warn|2|prints warn and all in level 1|
|Highlight|3|prints Highlight and all in level 2|
|Inform|4|Inform  and all in level 3|
|Log|5|prints logs and all in level 4|
|Trace|6|prints trace and all in level 5|


## BenchMark your app :

by using END on end of your app you can get kind of benchmark for your fun (Its not exactly benchmark , it only calculating time)

```bash

BEGIN : 


YourLogs...
...

END : 55.649µs

```