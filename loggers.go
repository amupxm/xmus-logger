package logger

import "fmt"

// Log logs a message at log level
func (l *logger) Log(a ...interface{}) {
	l.doLog(Log, a...)

}

// Logf logs a message at log level with string formater
func (l *logger) Logf(format string, a ...interface{}) {
	l.doLog(Log, fmt.Sprintf(format, a...))

}

// Alert logs a message at log level
func (l *logger) Alert(a ...interface{}) {
	l.doLog(Alert, a...)
}

// Alertf logs a message at log level with string formater
func (l *logger) Alertf(format string, a ...interface{}) {
	l.doLog(Alert, fmt.Sprintf(format, a...))

}

// Error logs a message at log level
func (l *logger) Error(a ...interface{}) {
	l.doLog(Error, a...)
}

// Errorf logs a message at log level with string formater
func (l *logger) Errorf(format string, a ...interface{}) {
	l.doLog(Error, fmt.Sprintf(format, a...))

}

// Highlight logs a message at log level
func (l *logger) Highlight(a ...interface{}) {
	l.doLog(Highlight, a...)
}

// Highlightf logs a message at log level with string formater
func (l *logger) Highlightf(format string, a ...interface{}) {
	l.doLog(Highlight, fmt.Sprintf(format, a...))

}

// Info logs a message at log level
func (l *logger) Info(a ...interface{}) {
	l.doLog(Info, a...)
}

// Infof logs a message at log level with string formater
func (l *logger) Infof(format string, a ...interface{}) {
	l.doLog(Info, fmt.Sprintf(format, a...))

}

// Trace logs a message at log level
func (l *logger) Trace(a ...interface{}) {
	l.doLog(Trace, a...)
}

// Tracef logs a message at log level with string formater
func (l *logger) Debugf(format string, a ...interface{}) {
	l.doLog(Trace, fmt.Sprintf(format, a...))

}

// Warn logs a message at log level
func (l *logger) Warn(a ...interface{}) {
	l.doLog(Warn, a...)
}

// Warnf logs a message at log level with string formater
func (l *logger) Warnf(format string, a ...interface{}) {
	l.doLog(Warn, fmt.Sprintf(format, a...))

}
