package xmuslogger

type Level int8

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	return [...]string{"trace", "debug", "info", "warn", "error", "fatal"}[l]
}
