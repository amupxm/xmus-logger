package xmuslogger

import "time"

type loggerWriter struct {
	parent *Logger
}

func (lw *loggerWriter) Write(p []byte) (n int, err error) {
	message := string(p)

	// Extract message from standard log format
	if len(message) > 20 && message[len(message)-1] == '\n' {
		if idx := findMessageStart(message); idx > 0 {
			message = message[idx:]
		}
		message = message[:len(message)-1] // Remove newline
	}

	// Build JSON
	var buf []byte
	if lw.parent.context != nil {
		buf = appendBytes(buf, lw.parent.context)
	}
	buf = appendString(buf, "message", message)
	buf = appendTime(buf, "time", time.Now())
	buf = appendString(buf, "level", "info")
	buf = appendString(buf, "source", "stdlib")

	jsonBuf := wrapJSON(buf)

	// Write to all outputs
	for _, w := range lw.parent.writers {
		_, err := w.Write(append(jsonBuf, '\n'))
		if err != nil {
			// TODO : handle error
		}
	}

	// Write to remote
	if lw.parent.remoteWriter != nil {
		if lw.parent.async {
			if err := lw.parent.remoteWriter.WriteAsync(jsonBuf); err != nil {
				// TODO : handle error
			}
		} else {
			if err := lw.parent.remoteWriter.Write(jsonBuf); err != nil {
				// TODO : handle error
			}
		}
	}

	return len(p), nil
}

func findMessageStart(logLine string) int {
	count := 0
	for i, c := range logLine {
		if c == ' ' {
			count++
			if count == 2 {
				return i + 1
			}
		}
	}
	return 0
}
