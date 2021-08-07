package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	logger "github.com/amupxm/xmus-logger/srv"
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
				LogLevel:    logger.Trace, // max log level
				Verbose:     false,
				File:        true,
				FilePath:    "string",
				Std:         true,
				UseCollores: true,
				Stdout:      &b,
			},
		)
		wordToTest := fmt.Sprint(v)
		logger.Log(v)
		logger.End()
		if wordToTest != b.String() {
			t.Errorf("ERROR :: Expected  : " + wordToTest + " GOT : " + b.String())
		}
	}
}

func TestLogLn(t *testing.T) {
	for _, v := range ValuesTotest {
		var b bytes.Buffer
		logger := logger.CreateLogger(
			&logger.LoggerOptions{
				LogLevel:    logger.Trace,
				Verbose:     false,
				File:        true,
				FilePath:    "string",
				Std:         true,
				UseCollores: true,
				Stdout:      &b,
			},
		)
		wordToTest := fmt.Sprint(v) + "\n"
		logger.Logln(v)
		logger.End()

		if wordToTest != b.String() {
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
				LogLevel:    logger.Trace,
				Verbose:     false,
				File:        true,
				FilePath:    "string",
				Std:         true,
				UseCollores: true,
				Stdout:      &b,
			},
		)
		wordToTest := fmt.Sprint(v[1:]...)
		logger.LogF(v[0].(string), v[1:]...)
		logger.End()
		if wordToTest != b.String() {
			t.Errorf("ERROR :: Expected  : " + wordToTest + " GOT : " + b.String())
		}
	}
}

func TestLogLevels(t *testing.T) {
	guideMapOfLevels := []logger.LogLevel{
		logger.Alert,
		logger.Error,
		logger.Warn,
		logger.Highlight,
		logger.Inform,
		logger.Log,
		logger.Trace,
	}
	//	li <= len cuz we need to test all levels + 1
	for li := 0; li <= len(guideMapOfLevels); li++ {
		currentLevel := guideMapOfLevels[li]
		strToTest := "1"
		for _, f := range guideMapOfLevels {
			var b bytes.Buffer
			log := logger.CreateLogger(
				&logger.LoggerOptions{
					LogLevel:    currentLevel,
					Verbose:     false,
					File:        true,
					FilePath:    "string",
					Std:         true,
					UseCollores: true,
					Stdout:      &b,
				},
			)
			mapOfFuncs := map[logger.LogLevel]func(a ...interface{}) logger.LogResult{
				logger.Alert:     log.Alert,
				logger.Error:     log.Error,
				logger.Warn:      log.Warn,
				logger.Highlight: log.Highlight,
				logger.Inform:    log.Inform,
				logger.Log:       log.Log,
				logger.Trace:     log.Trace,
			}
			mapOfFuncLns := map[logger.LogLevel]func(a ...interface{}) logger.LogResult{
				logger.Alert:     log.Alertln,
				logger.Error:     log.Errorln,
				logger.Warn:      log.Warnln,
				logger.Highlight: log.Highlightln,
				logger.Inform:    log.Informln,
				logger.Log:       log.Logln,
				logger.Trace:     log.Traceln,
			}
			mapOfFuncFs := map[logger.LogLevel]func(format string, a ...interface{}) logger.LogResult{
				logger.Alert:     log.AlertF,
				logger.Error:     log.ErrorF,
				logger.Warn:      log.WarnF,
				logger.Highlight: log.HighlightF,
				logger.Inform:    log.InformF,
				logger.Log:       log.LogF,
				logger.Trace:     log.TraceF,
			}

			mapOfFuncs[f](strToTest)
			mapOfFuncLns[f](strToTest)
			mapOfFuncFs[f]("%s", strToTest)

			var ac logger.LogLevel
			if f > currentLevel {
				ac = currentLevel
			} else {
				ac = f
			}
			if len(b.String()) != int(ac)*len(strToTest)+int(ac)*len(strToTest) {
				t.Errorf("ERROR :: Expected  : " + fmt.Sprint(int(ac)*len(strToTest)+int(ac)*len(strToTest)) + " GOT : " + b.String() + "in loglevel " + fmt.Sprint(currentLevel))
			}
		}
	}

	//Trace is bigger one
	// for c := 0; c < int(logger.Trace); c++ {
	// 	i := c
	// 	if i >= len(guideMapOfLevels) {
	// 		i = len(guideMapOfLevels) - 1
	// 	}
	// 	currentLevel := guideMapOfLevels[c]
	// 	var b bytes.Buffer
	// 	log := logger.CreateLogger(
	// 		&logger.LoggerOptions{
	// 			LogLevel:    currentLevel,
	// 			Verbose:     true,
	// 			File:        true,
	// 			FilePath:    "string",
	// 			Std:         true,
	// 			UseCollores: true,
	// 			Stdout:      &b,
	// 		},
	// 	)
	// 	mapOfFuncs := map[logger.LogLevel]func(a ...interface{}) logger.LogResult{
	// 		logger.Alert:     log.Alert,
	// 		logger.Error:     log.Error,
	// 		logger.Warn:      log.Warn,
	// 		logger.Highlight: log.Highlight,
	// 		logger.Inform:    log.Inform,
	// 		logger.Log:       log.Log,
	// 		logger.Trace:     log.Trace,
	// 	}
	// 	mapOfFuncLns := map[logger.LogLevel]func(a ...interface{}) logger.LogResult{
	// 		logger.Alert:     log.Alertln,
	// 		logger.Error:     log.Errorln,
	// 		logger.Warn:      log.Warnln,
	// 		logger.Highlight: log.Highlightln,
	// 		logger.Inform:    log.Informln,
	// 		logger.Log:       log.Logln,
	// 		logger.Trace:     log.Traceln,
	// 	}
	// 	mapOfFuncFs := map[logger.LogLevel]func(format string, a ...interface{}) logger.LogResult{
	// 		logger.Alert:     log.AlertF,
	// 		logger.Error:     log.ErrorF,
	// 		logger.Warn:      log.WarnF,
	// 		logger.Highlight: log.HighlightF,
	// 		logger.Inform:    log.InformF,
	// 		logger.Log:       log.LogF,
	// 		logger.Trace:     log.TraceF,
	// 	}

	// 	testStr := "testStr"
	// 	for _, level := range guideMapOfLevels {
	// 		mapOfFuncs[level](testStr)
	// 		mapOfFuncLns[level](testStr)
	// 		mapOfFuncFs[level]("%s", testStr)
	// 		if level > currentLevel {
	// 			if b.String() != "" {
	// 				t.Errorf("ERROR :: Expected  : \"\" GOT : " + b.String())
	// 			}
	// 		}
	// 	}
	// }
}
