package main

import (
	"github.com/yzletter/go-toolery/loggerx"
)

func main() {
	a()
}

func a() {
	b()
}

func b() {
	log := loggerx.NewLog("./loggerx/my.log", loggerx.DebugLevel)
	log.Debug("1")
	log.Info("1")
	log.Error("1")
}
