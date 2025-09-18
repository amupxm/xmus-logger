package xmuslogger

import (
	"fmt"
	"time"
	"unicode/utf8"
)

// Fixed serializer.go

// Properly escape JSON strings
func appendString(dst []byte, key, val string) []byte {
	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"', ':', '"')
	dst = appendEscapedString(dst, val)
	dst = append(dst, '"', ',')
	return dst
}

// Escape special characters in JSON strings
func appendEscapedString(dst []byte, s string) []byte {
	for i, r := range s {
		switch r {
		case '"':
			dst = append(dst, '\\', '"')
		case '\\':
			dst = append(dst, '\\', '\\')
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			if r < 0x20 {
				// Control characters need unicode escape
				dst = append(dst, fmt.Sprintf("\\u%04x", r)...)
			} else if r == utf8.RuneError {
				// Handle invalid UTF-8
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					// Invalid UTF-8, replace with unicode escape
					dst = append(dst, fmt.Sprintf("\\u%04x", s[i])...)
				} else {
					// Valid multi-byte character
					dst = append(dst, string(r)...)
				}
			} else {
				// Normal character
				dst = append(dst, string(r)...)
			}
		}
	}
	return dst
}

func appendBytes(dst []byte, val []byte) []byte {
	dst = append(dst, val...)
	if len(dst) > 0 && dst[len(dst)-1] != ',' {
		dst = append(dst)
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

func appendInt64(dst []byte, key string, val int64) []byte {
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
