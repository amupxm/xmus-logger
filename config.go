package xmuslogger

import "io"

func (l *Logger) Level(level Level) *Logger {
	newLogger := l.clone()
	newLogger.level = level
	newLogger.updateEnabledLevels()
	return newLogger
}

func (l *Logger) Output(w io.Writer) *Logger {
	newLogger := l.clone()
	newLogger.SetOutput(w)
	return newLogger
}

func (l *Logger) Remote(w RemoteWriter) *Logger {
	newLogger := l.clone()
	newLogger.remoteWriter = w
	return newLogger
}

func (l *Logger) RemoteHTTP(endpoint string, options ...HTTPOption) *Logger {
	return l.Remote(NewHTTPRemoteWriter(endpoint, options...))
}

func (l *Logger) With() *Context {
	return &Context{logger: l}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.SetOutput(&loggerWriter{parent: l})
	l.writers = []io.Writer{w}
}

func (l *Logger) Flush() error {
	if l.remoteWriter != nil {
		return l.remoteWriter.Flush()
	}
	return nil
}

func (l *Logger) Close() error {
	if l.remoteWriter != nil {
		return l.remoteWriter.Close()
	}
	return nil
}
