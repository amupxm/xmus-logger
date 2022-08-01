package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/reiver/go-xim"
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
		traceCode    string
	}
	LoggerOptions struct {
		LogLevel LogLevel
		Verbose  bool
		FilePath string
		Std      bool
		Stdout   io.Writer
	}
	Logger interface {
		BeginWithPrefix(format ...string) Logger
		Begin() Logger
		// use for add custom output
		SetCustomOut(outPutt io.Writer)
		// doLog send log to stdout or file
		doLog(level LogLevel, v ...interface{})
		// End send finished signal to log
		End()
		// Prefix the log with a string
		Prefix(format ...string) *logger
		// GetCaller return the caller of the log
		GetCaller() *logger

		// Log logs a message at log level
		Log(v ...interface{}) LogResult
		// Logf logs a message at log level with string formater
		Logf(format string, v ...interface{})

		// Alert logs a message at log level
		Alert(v ...interface{}) LogResult
		// Alertf logs a message at log level with string formater
		Alertf(format string, v ...interface{})

		// Error logs a message at log level
		Error(v ...interface{}) LogResult
		// Errorf logs a message at log level with string formater
		Errorf(format string, v ...interface{})

		// Highlight logs a message at log level
		Highlight(v ...interface{}) LogResult
		// Highlightf logs a message at log level with string formater
		Highlightf(format string, v ...interface{})

		// Info logs a message at log level
		Info(v ...interface{}) LogResult
		// Infof logs a message at log level with string formater
		Infof(format string, v ...interface{})

		// Trace logs a message at log level
		Trace(v ...interface{}) LogResult
		// Tracef logs a message at log level with string formater
		Debugf(format string, v ...interface{})

		// Warn logs a message at log level
		Warn(v ...interface{}) LogResult
		// Warnf logs a message at log level with string formater
		Warnf(format string, v ...interface{})
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
	Info                      //5
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

func (l logger) Begin() Logger {
	return l.BeginWithPrefix("")
}

func (l *logger) Level(level uint8) Logger {
	l.LogLevel = LogLevel(level)
	return l
}

// use for add custom output
func (l *logger) SetCustomOut(outPutt io.Writer) {
	l.stdout = outPutt
}

func createUID() string {
	var x = xim.Generate()
	return x.String()
}

// Prefix the log with a string
func (l logger) BeginWithPrefix(format ...string) Logger {
	colorReset := "\033[0m"

	colorRed := "\033[31m"
	// colorYellow := "\033[33m"
	// colorBlue := "\033[34m"
	// colorPurple := "\033[35m"
	// colorCyan := "\033[36m"
	// colorWhite := "\033[37m"

	clone := logger{
		LogLevel:  l.LogLevel,
		verbose:   l.verbose,
		filePath:  l.filePath,
		std:       l.std,
		stdout:    l.stdout,
		traceCode: createUID(),
	}
	var tmpF []string
	for _, v := range format {
		v = strings.ToUpper(v)
		sv := strings.Split(v, "/")
		tmpF = append(tmpF, sv...)
	}
	clone.prefixString = strings.Join(tmpF, ":")
	clone.prefixString = fmt.Sprintf("%s%s%s", colorRed, clone.prefixString, colorReset)
	t := time.Now()
	clone.time = &t

	if clone.verbose {
		l.doLog(Alert, "BEGIN : ")
	}
	return &clone
}

// Prefix the log with a string
func (l logger) Prefix(format ...string) *logger {
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
func (l logger) End() {

	if l.time != nil {
		d := time.Since(
			*l.time,
		)
		l.duration = &d
		if l.verbose {
			l.Logf("END : %s", l.duration.String())
		}
		l.prefixString = ""
	}

}

func (lr *logResult) TraceStack() {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	lr.logger.Logf("%s", stackSlice[0:s])
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
	l.Logf("%s :: ", l.getCaller())
	return l
}

// doLog send log to stdout or file
func (l logger) doLog(level LogLevel, v ...interface{}) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	// Check log level permission :
	// permission Nothing is not allowed to log
	if l.LogLevel <= Nothing || level > l.LogLevel {
		return
	}
	var msg string
	// to write to file
	if l.prefixString != "" {
		msg = fmt.Sprintf("[%s]", l.prefixString)
	}
	if l.traceCode != "" {
		msg = fmt.Sprintf("%s[%s]", msg, fmt.Sprintf("%s%s%s", colorGreen, l.traceCode, colorReset))
	}
	for _, v := range v {
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
