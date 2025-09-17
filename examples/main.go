package main

import (
	"fmt"
	"log"

	xmuslogger "github.com/amupxm/xmus-logger"
)

func main() {
	logxmus := xmuslogger.New().Level(xmuslogger.DebugLevel)
	log.SetPrefix("main")
	ctxLogger := logxmus.With().Str("component", "main").Logger()
	ctxLogger.Info().Msg("This is an info message")
	ctxLogger.Debug().Msg("This is a debug message")
	logxmus.Error().Err(fmt.Errorf("test error")).Msg("test")
	ctxLogger.Print("Hello world")
	anyFuncwithDefaultLogger(logxmus)
}

func anyFuncwithDefaultLogger(log *xmuslogger.Logger) {
	log.Println("This is a log message from anyFuncwithDefaultLogger")
}
