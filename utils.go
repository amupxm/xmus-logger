package logger

func ArrayContains[T comparable](v []T, s T) bool {
	for _, vv := range v {
		if vv == s {
			return true
		}
	}
	return false
}
