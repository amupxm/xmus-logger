package logger

import "fmt"

// Log logs a message at log level
func (l *logger) Log(a ...interface{}) LogResult {
	l.doLog(Log, a...)
	return &logResult{
		logger: l,
	}
}

// Logln logs a message at log level to new line
func (l *logger) Logln(a ...interface{}) LogResult {
	l.doLog(Log, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// LogF logs a message at log level with string formater
func (l *logger) LogF(format string, a ...interface{}) LogResult {
	l.doLog(Log, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Alert logs a message at log level
func (l *logger) Alert(a ...interface{}) LogResult {
	l.doLog(Alert, a...)
	return &logResult{
		logger: l,
	}
}

// Alertln logs a message at log level to new line
func (l *logger) Alertln(a ...interface{}) LogResult {
	l.doLog(Alert, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// AlertF logs a message at log level with string formater
func (l *logger) AlertF(format string, a ...interface{}) LogResult {
	l.doLog(Alert, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Error logs a message at log level
func (l *logger) Error(a ...interface{}) LogResult {
	l.doLog(Error, a...)
	return &logResult{
		logger: l,
	}
}

// Errorln logs a message at log level to new line
func (l *logger) Errorln(a ...interface{}) LogResult {
	l.doLog(Error, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// ErrorF logs a message at log level with string formater
func (l *logger) ErrorF(format string, a ...interface{}) LogResult {
	l.doLog(Error, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Highlight logs a message at log level
func (l *logger) Highlight(a ...interface{}) LogResult {
	l.doLog(Highlight, a...)
	return &logResult{
		logger: l,
	}
}

// Highlightln logs a message at log level to new line
func (l *logger) Highlightln(a ...interface{}) LogResult {
	l.doLog(Highlight, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// HighlightF logs a message at log level with string formater
func (l *logger) HighlightF(format string, a ...interface{}) LogResult {
	l.doLog(Highlight, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Inform logs a message at log level
func (l *logger) Inform(a ...interface{}) LogResult {
	l.doLog(Inform, a...)
	return &logResult{
		logger: l,
	}
}

// Informln logs a message at log level to new line
func (l *logger) Informln(a ...interface{}) LogResult {
	l.doLog(Inform, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// InformF logs a message at log level with string formater
func (l *logger) InformF(format string, a ...interface{}) LogResult {
	l.doLog(Inform, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Trace logs a message at log level
func (l *logger) Trace(a ...interface{}) LogResult {
	l.doLog(Trace, a...)
	return &logResult{
		logger: l,
	}
}

// Traceln logs a message at log level to new line
func (l *logger) Traceln(a ...interface{}) LogResult {
	l.doLog(Trace, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// TraceF logs a message at log level with string formater
func (l *logger) TraceF(format string, a ...interface{}) LogResult {
	l.doLog(Trace, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}

// Warn logs a message at log level
func (l *logger) Warn(a ...interface{}) LogResult {
	l.doLog(Warn, a...)
	return &logResult{
		logger: l,
	}
}

// Warnln logs a message at log level to new line
func (l *logger) Warnln(a ...interface{}) LogResult {
	l.doLog(Warn, append(a, "\n")...)
	return &logResult{
		logger: l,
	}
}

// WarnF logs a message at log level with string formater
func (l *logger) WarnF(format string, a ...interface{}) LogResult {
	l.doLog(Warn, fmt.Sprintf(format, a...))
	return &logResult{
		logger: l,
	}
}
