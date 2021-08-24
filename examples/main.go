package main

import logger "github.com/amupxm/xmus-logger"

func main() {
	// logOptions := logger.LoggerOptions{
	// 	LogLevel: 6,
	// 	Verbose:  true,
	// 	Std:      true,
	// }
	// log := logger.CreateLogger(&logOptions)
	// log.Logln("log")
	// pref := log.Prefix("Prefix", "text")
	// pref.Logln("from prefixlogger").TraceStack()
	log := logger.Begin()
	log.Level(uint8(logger.Log))
	log.Log("log")
	log.Alert("alert")
	log.End()

}
