package main

import (
	"time"

	logger "github.com/amupxm/xmus-logger"
)

func main() {
	logOptions := logger.Options{
		LogLevel: 6,
		Verbose:  true,
		Std:      true,
	}
	log := logger.CreateLogger(&logOptions)
	c := log.BeginWithPrefix("main")
	c.Errorf("asdasd %s", "asdasd")
	time.Sleep(time.Second * 1)
	c.GetCaller()
	c.Errorf("asdasd %s", "asdasd")
	c2 := log.BeginWithPrefix("main2")
	c2.AddToWhitelist("main2")

	c2.Errorf("asdasd %s", "asdasd")

	log.End()

}
