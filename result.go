package logger

func GetResult(l *logger) *logResult {
	return &logResult{
		logger: l,
	}
}
