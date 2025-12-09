package loggerx_test

import (
	"testing"

	"github.com/yzletter/go-toolery/loggerx"
)

func TestLog(T *testing.T) {
	log := loggerx.NewLog("../loggerx/my.log", loggerx.DebugLevel)
	log.Debug("1")
	log.Info("1")
	log.Error("1")
}
