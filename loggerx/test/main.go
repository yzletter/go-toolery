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
	logger := loggerx.NewLog("./loggerx/test/my.log", loggerx.DebugLevel)
	//log := loggerx.NewLog("", loggerx.DebugLevel)
	defer logger.Close()
	logger.AddStackTrace()
	logger.SetUDPProducer("127.0.0.1:4321")

	logger.Debugf("这是一条%s日志", "debug")
	logger.Infof("这是一条%s日志", "info")
	logger.Warnf("这是一条%s日志", "warn")
	logger.Errorf("这是一条%s日志", "error")
}
