package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	logger "github.com/amupxm/xmus-logger"
)

var ValuesTotest = []interface{}{
	"ali",
	"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
	"",
	"#@!$%&^*)(*!@#^(!@",
	"0000000000000000000000000001",
	1234123,
	func() string { return "hello" }(),
	0.00001,
	true,
}

func TestLog(t *testing.T) {
	for _, v := range ValuesTotest {
		var b bytes.Buffer
		logger := logger.CreateLogger(
			&logger.LoggerOptions{
				LogLevel: logger.Trace, // max log level
				Verbose:  false,
				FilePath: "string",
				Std:      true,
				Stdout:   &b,
			},
		)
		wordToTest := fmt.Sprint(v)
		logger.Log(v)
		logger.End()
		if fmt.Sprintf("%s\n", wordToTest) != b.String() {
			t.Errorf("ERROR :: Expected  : " + wordToTest + " GOT : " + b.String())
		}
	}
}

func TestLogf(t *testing.T) {
	testCases := [][]interface{}{
		{"%s", "ali"},
		{"%s", "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."},
		{"%s", ""},
		{"%s", "#@!$%&^*)(*!@#^(!@"},
		{"%s", "0000000000000000000000000001"},
		{"%d", 1234123},
		{"%s", "hello"},
		// Bug : KNOWN BUG:
		// {"%f", 0.00001},
		// error in 1e-05%!(EXTRA string= Expected  : 1e-05 Got : 0.000010)
		{"%t", true},
	}
	for _, v := range testCases {
		var b bytes.Buffer
		logger := logger.CreateLogger(
			&logger.LoggerOptions{
				LogLevel: logger.Trace,
				Verbose:  false,
				FilePath: "string",
				Std:      true,
				Stdout:   &b,
			},
		)
		wordToTest := fmt.Sprint(v[1:]...)
		logger.Logf(v[0].(string), v[1:]...)
		logger.End()

		if fmt.Sprintf("%s\n", wordToTest) != b.String() {
			t.Errorf("ERROR :: Expected  : " + wordToTest + " GOT : " + b.String())
		}
	}
}
