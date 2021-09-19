package logger

import "fmt"

// Log logs a message at log level
func (l *logger) Log(a ...interface{}) LogResult {
	l.doLog(Log, a...)
	return &logResult{
		logger: l,
	}
}

// Logf logs a message at log level with string formater
func (l *logger) Logf(format string, a ...interface{}) {
	l.doLog(Log, fmt.Sprintf(format, a...))

}

// Alert logs a message at log level
func (l *logger) Alert(a ...interface{}) LogResult {
	l.doLog(Alert, a...)
	return &logResult{
		logger: l,
	}
}

// Alertf logs a message at log level with string formater
func (l *logger) Alertf(format string, a ...interface{}) {
	l.doLog(Alert, fmt.Sprintf(format, a...))

}

// Error logs a message at log level
func (l *logger) Error(a ...interface{}) LogResult {
	l.doLog(Error, a...)
	return &logResult{
		logger: l,
	}
}

// Errorf logs a message at log level with string formater
func (l *logger) Errorf(format string, a ...interface{}) {
	l.doLog(Error, fmt.Sprintf(format, a...))

}

// Highlight logs a message at log level
func (l *logger) Highlight(a ...interface{}) LogResult {
	l.doLog(Highlight, a...)
	return &logResult{
		logger: l,
	}
}

// Highlightf logs a message at log level with string formater
func (l *logger) Highlightf(format string, a ...interface{}) {
	l.doLog(Highlight, fmt.Sprintf(format, a...))

}

// Inform logs a message at log level
func (l *logger) Inform(a ...interface{}) LogResult {
	l.doLog(Inform, a...)
	return &logResult{
		logger: l,
	}
}

// Informf logs a message at log level with string formater
func (l *logger) Informf(format string, a ...interface{}) {
	l.doLog(Inform, fmt.Sprintf(format, a...))

}

// Trace logs a message at log level
func (l *logger) Trace(a ...interface{}) LogResult {
	l.doLog(Trace, a...)
	return &logResult{
		logger: l,
	}
}

// Tracef logs a message at log level with string formater
func (l *logger) Debugf(format string, a ...interface{}) {
	l.doLog(Trace, fmt.Sprintf(format, a...))

}

// Warn logs a message at log level
func (l *logger) Warn(a ...interface{}) LogResult {
	l.doLog(Warn, a...)
	return &logResult{
		logger: l,
	}
}

// Warnf logs a message at log level with string formater
func (l *logger) Warnf(format string, a ...interface{}) {
	l.doLog(Warn, fmt.Sprintf(format, a...))

}
