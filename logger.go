package xmuslogger

import (
	"io"
	"log"
	"os"
	"sync"
)

type XmusLogger interface {
	// Standard Go logger compatibility
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	SetOutput(w io.Writer)
	SetFlags(flag int)
	SetPrefix(prefix string)
	Writer() io.Writer
	Flags() int
	Prefix() string

	// Zerolog-style methods
	Trace() *Event
	Debug() *Event
	Info() *Event
	Warn() *Event
	Error() *Event

	// Configuration
	Level(level Level) *Logger
	Output(w io.Writer) *Logger
	Remote(w RemoteWriter) *Logger
	RemoteHTTP(endpoint string, options ...HTTPOption) *Logger
	With() *Context

	// Lifecycle
	Flush() error
	Close() error
}

type RemoteWriter interface {
	Write(data []byte) error
	WriteAsync(data []byte) error
	Flush() error
	Close() error
}

type Logger struct {
	*log.Logger               // Embedded for compatibility
	level        Level        // Current log level
	writers      []io.Writer  // Local outputs
	remoteWriter RemoteWriter // Remote output
	context      []byte       // Pre-serialized context
	async        bool         // Async remote sending
	mu           sync.RWMutex // Thread safety
	enabled      [8]bool      // Level cache
}

// Constructor
func New() *Logger {
	l := &Logger{
		level:   InfoLevel,
		writers: []io.Writer{os.Stdout},
		context: []byte{},
	}

	// Route standard logger through our JSON formatter
	l.Logger = log.New(&loggerWriter{parent: l}, "", log.LstdFlags)
	l.updateEnabledLevels()

	return l
}

func NewWithOutput(w io.Writer) *Logger {
	l := &Logger{
		level:   InfoLevel,
		writers: []io.Writer{w},
		context: []byte{},
	}

	l.Logger = log.New(&loggerWriter{parent: l}, "", log.LstdFlags)
	l.updateEnabledLevels()

	return l
}

// GLOBAL FUNCTIONS
var std = New()

func Print(v ...interface{})                 { std.Print(v...) }
func Printf(format string, v ...interface{}) { std.Printf(format, v...) }
func Println(v ...interface{})               { std.Println(v...) }
func Fatal(v ...interface{})                 { std.Fatal(v...) }
func Fatalf(format string, v ...interface{}) { std.Fatalf(format, v...) }
func Fatalln(v ...interface{})               { std.Fatalln(v...) }
func Panic(v ...interface{})                 { std.Panic(v...) }
func Panicf(format string, v ...interface{}) { std.Panicf(format, v...) }
func Panicln(v ...interface{})               { std.Panicln(v...) }

func SetOutput(w io.Writer)   { std.SetOutput(w) }
func SetFlags(flag int)       { std.SetFlags(flag) }
func SetPrefix(prefix string) { std.SetPrefix(prefix) }

func SetDefault(l *Logger) { std = l }
func Default() *Logger     { return std }

// ============================================================================
// USAGE EXAMPLES
// ============================================================================

/*
// Standard Go logger (produces JSON)
log := New()
log.Printf("User %s logged in", "john")
// Output: {"message":"User john logged in","time":"...","level":"info","source":"stdlib"}

// Zerolog style (produces JSON)
log.Info().
    Str("user", "john").
    Int("age", 30).
    Msg("Login successful")
// Output: {"user":"john","age":30,"message":"Login successful","time":"...","level":"info"}

// With remote logging
log := New().RemoteHTTP("https://logs.company.com/api/v1/logs")
log.Printf("This goes to HTTP endpoint too")

// With context
ctxLog := log.With().Str("service", "auth").Logger()
ctxLog.Info().Msg("Service started")
*/
