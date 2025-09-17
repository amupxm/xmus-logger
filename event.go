package xmuslogger

import (
	"fmt"
	"io"
	"time"
)

type Event struct {
	buf          []byte
	writers      []io.Writer
	remoteWriter RemoteWriter
	level        Level
	done         func(*Event)
	async        bool
}

// Field methods
func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return e
	}
	e.buf = appendString(e.buf, key, val)
	return e
}

func (e *Event) Int(key string, val int) *Event {
	if e == nil {
		return e
	}
	e.buf = appendInt(e.buf, key, val)
	return e
}

func (e *Event) Bool(key string, val bool) *Event {
	if e == nil {
		return e
	}
	e.buf = appendBool(e.buf, key, val)
	return e
}

func (e *Event) Err(err error) *Event {
	if e == nil || err == nil {
		return e
	}
	e.buf = appendString(e.buf, "error", err.Error())
	return e
}

// Message methods
func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}

	e.buf = appendString(e.buf, "message", msg)
	e.buf = appendTime(e.buf, "time", time.Now())
	e.buf = appendString(e.buf, "level", e.level.String())

	finalBuf := wrapJSON(e.buf)

	// Write to local outputs
	for _, w := range e.writers {
		_, err := w.Write(append(finalBuf, '\n'))
		if err != nil {
			// TODO : handle error
		}
	}

	// Write to remote
	if e.remoteWriter != nil {
		if e.async {
			if err := e.remoteWriter.WriteAsync(finalBuf); err != nil {
				// TODO : handle error
			}
		} else {
			if err := e.remoteWriter.Write(finalBuf); err != nil {
				// TODO : handle error
			}
		}
	}

	// Return to pool
	if e.done != nil {
		e.done(e)
	}
}

func (e *Event) Msgf(format string, v ...interface{}) {
	if e == nil {
		return
	}
	e.Msg(fmt.Sprintf(format, v...))
}

func (e *Event) Send() {
	if e == nil {
		return
	}
	e.Msg("")
}

func (l *Logger) Trace() *Event { return l.newEvent(TraceLevel) }
func (l *Logger) Debug() *Event { return l.newEvent(DebugLevel) }
func (l *Logger) Info() *Event  { return l.newEvent(InfoLevel) }
func (l *Logger) Warn() *Event  { return l.newEvent(WarnLevel) }
func (l *Logger) Error() *Event { return l.newEvent(ErrorLevel) }

func (l *Logger) newEvent(level Level) *Event {
	if !l.enabled[level] {
		return nil // Zero cost for disabled levels
	}

	e := getEvent()
	e.level = level
	e.writers = l.writers
	e.remoteWriter = l.remoteWriter
	e.async = l.async
	e.done = putEvent

	// Copy pre-serialized context
	if len(l.context) > 0 {
		e.buf = append(e.buf[:0], l.context...)
	}

	return e
}
