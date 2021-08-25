package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type (
	logger struct {
		time         *time.Time
		LogLevel     LogLevel
		verbose      bool
		std          bool
		filePath     string
		duration     *time.Duration
		prefixString string
		stdout       io.Writer
	}
	LoggerOptions struct {
		LogLevel LogLevel
		Verbose  bool
		FilePath string
		Std      bool
		Stdout   io.Writer
	}
	Logger interface {
		Begin() Logger
		// use for add custom output
		SetCustomOut(outPutt io.Writer)
		// doLog send log to stdout or file
		doLog(level LogLevel, a ...interface{})
		// End send finished signal to log
		End()
		// Prefix the log with a string
		Prefix(format ...string) *logger
		// GetCaller return the caller of the log
		GetCaller() *logger

		// Log logs a message at log level
		Log(a ...interface{}) LogResult
		// LogF logs a message at log level with string formater
		LogF(format string, a ...interface{}) LogResult

		// Alert logs a message at log level
		Alert(a ...interface{}) LogResult
		// AlertF logs a message at log level with string formater
		AlertF(format string, a ...interface{}) LogResult

		// Error logs a message at log level
		Error(a ...interface{}) LogResult
		// ErrorF logs a message at log level with string formater
		ErrorF(format string, a ...interface{}) LogResult

		// Highlight logs a message at log level
		Highlight(a ...interface{}) LogResult
		// HighlightF logs a message at log level with string formater
		HighlightF(format string, a ...interface{}) LogResult

		// Inform logs a message at log level
		Inform(a ...interface{}) LogResult
		// InformF logs a message at log level with string formater
		InformF(format string, a ...interface{}) LogResult

		// Trace logs a message at log level
		Trace(a ...interface{}) LogResult
		// TraceF logs a message at log level with string formater
		TraceF(format string, a ...interface{}) LogResult

		// Warn logs a message at log level
		Warn(a ...interface{}) LogResult
		// WarnF logs a message at log level with string formater
		WarnF(format string, a ...interface{}) LogResult
		// Set LogLevel
		Level(level uint8) Logger
	}
	logResult struct {
		logger *logger
	}
	LogResult interface {
		// TraceStack trace the stack of the log caller
		TraceStack()
	}
	LogLevel int
)

const (
	Nothing   LogLevel = iota //0
	Alert                     //1
	Error                     //2
	Warn                      //3
	Highlight                 //4
	Inform                    //5
	Log                       //6
	Trace                     //7
)

func CreateLogger(LoggerOpts *LoggerOptions) Logger {

	if LoggerOpts.LogLevel > Trace {
		LoggerOpts.LogLevel = Trace
	}

	// Cuz alert and error are in 1 level
	if LoggerOpts.LogLevel >= Alert {
		LoggerOpts.LogLevel += 1
	}

	if LoggerOpts.Stdout == nil {
		LoggerOpts.Stdout = os.Stdout
	}

	l := &logger{
		LogLevel:     LoggerOpts.LogLevel,
		verbose:      LoggerOpts.Verbose,
		filePath:     LoggerOpts.FilePath,
		std:          LoggerOpts.Std,
		stdout:       LoggerOpts.Stdout,
		prefixString: "",
	}
	t := time.Now()
	l.time = &t

	if l.verbose {
		l.doLog(Alert, "BEGIN : ")
	}
	return l
}

func (l *logger) Begin() Logger {

	t := time.Now()
	l.time = &t

	if l.verbose {
		l.doLog(Alert, "BEGIN : ")
	}
	return l
}

func (l *logger) Level(level uint8) Logger {
	l.LogLevel = LogLevel(level)
	return l
}

// use for add custom output
func (l *logger) SetCustomOut(outPutt io.Writer) {
	l.stdout = outPutt
}

// Prefix the log with a string
func (l *logger) Prefix(format ...string) *logger {
	clone := logger{
		LogLevel: l.LogLevel,
		verbose:  l.verbose,
		filePath: l.filePath,
		std:      l.std,
		stdout:   l.stdout,
	}
	clone.prefixString = strings.Join(format, ":")
	return &clone
}

// End send finished signal to log
func (l *logger) End() {
	d := time.Since(
		*l.time,
	)
	l.duration = &d
	if l.verbose {
		l.LogF("END : %s", l.duration.String())
	}
	l.prefixString = ""
}

func (lr *logResult) TraceStack() {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	lr.logger.LogF("%s", stackSlice[0:s])
}

func (l *logger) getCaller() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return ""
	}
	caller := runtime.FuncForPC(fpcs[0] - 2)
	return fmt.Sprintf("%s() ", caller.Name())
}

// GetCaller return the caller of the log
func (l *logger) GetCaller() *logger {
	l.LogF("%s :: ", l.getCaller())
	return l
}

// doLog send log to stdout or file
func (l *logger) doLog(level LogLevel, a ...interface{}) {

	// Check log level permission :
	// permission Nothing is not allowed to log
	if l.LogLevel <= Nothing || level > l.LogLevel {
		return
	}
	var msg string
	// to write to file

	if l.prefixString != "" {
		msg = fmt.Sprintf("%s =>", l.prefixString)
	}
	for _, v := range a {
		msg = fmt.Sprintf("%s%s", msg, fmt.Sprint(v))

	}
	l.wStd(fmt.Sprintf("%s\n", msg))
	// TODO make write to file

}

func (l *logger) wStd(msg interface{}) {
	fmt.Fprint(
		l.stdout,
		msg,
	)
}
