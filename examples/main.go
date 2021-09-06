package main

import logger "github.com/amupxm/xmus-logger"

func main() {
	logOptions := logger.LoggerOptions{
		LogLevel: 6,
		Verbose:  false,
		Std:      true,
	}
	log := logger.CreateLogger(&logOptions)
	c := log.BeginWithPrefix("main")
	c.Prefix("main").Alert(1)
	c.Error(23)
	log.End()

}
