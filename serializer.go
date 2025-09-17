package xmuslogger

import (
	"fmt"
	"time"
)

func appendString(dst []byte, key, val string) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"', ':', '"')
	dst = append(dst, val...)
	dst = append(dst, '"', ',')
	return dst
}

func appendBytes(dst []byte, val []byte) []byte {
	dst = append(dst, val...)
	if len(dst) > 0 && dst[len(dst)-1] != ',' {
		dst = append(dst, ',')
	}
	return dst
}

func appendInt(dst []byte, key string, val int) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"', ':')
	dst = append(dst, fmt.Sprintf("%d", val)...)
	dst = append(dst, ',')
	return dst
}

func appendBool(dst []byte, key string, val bool) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"', ':')
	if val {
		dst = append(dst, "true"...)
	} else {
		dst = append(dst, "false"...)
	}
	dst = append(dst, ',')
	return dst
}

func appendTime(dst []byte, key string, val time.Time) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"', ':', '"')
	dst = val.AppendFormat(dst, time.RFC3339)
	dst = append(dst, '"', ',')
	return dst
}

func wrapJSON(buf []byte) []byte {
	if len(buf) > 0 && buf[len(buf)-1] == ',' {
		buf = buf[:len(buf)-1] // Remove trailing comma
	}

	result := make([]byte, 0, len(buf)+2)
	result = append(result, '{')
	result = append(result, buf...)
	result = append(result, '}')

	return result
}
