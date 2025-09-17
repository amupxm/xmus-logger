package xmuslogger

import (
	"io"
	"log"
)

func (l *Logger) clone() *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	newLogger := &Logger{
		level:        l.level,
		writers:      make([]io.Writer, len(l.writers)),
		remoteWriter: l.remoteWriter,
		context:      make([]byte, len(l.context)),
		async:        l.async,
	}

	copy(newLogger.writers, l.writers)
	copy(newLogger.context, l.context)

	newLogger.Logger = log.New(&loggerWriter{parent: newLogger}, l.Prefix(), l.Flags())
	newLogger.updateEnabledLevels()

	return newLogger
}

func (l *Logger) updateEnabledLevels() {
	for i := range l.enabled {
		l.enabled[i] = Level(i) >= l.level
	}
}
