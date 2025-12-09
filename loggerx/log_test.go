package loggerx_test

import (
	"os"
	"testing"

	"github.com/yzletter/go-toolery/loggerx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkMyLogger(b *testing.B) {
	logger := loggerx.NewLog("../loggerx/log/my.log", loggerx.InfoLevel)
	for b.Loop() {
		logger.Infof("add suffix %s to file %s failed: %s", ".png", "path/to/file", "forbidden")
	}
}

func BenchmarkZap(b *testing.B) {
	file, err := os.OpenFile("../loggerx/log/zap.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //指定时间格式
	encoderConfig.TimeKey = "time"                                                    //默认是ts
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder  // level会带颜色
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), //日志为console格式（Field还是json格式）
		zapcore.AddSync(file),                    //指定输出到文件
		zapcore.InfoLevel,                        //设置最低级别
	)
	logger := zap.New(
		core,
		zap.AddCaller(), //上报文件名和行号
		// zap.AddStacktrace(zapcore.ErrorLevel), //error级别及其以上的日志打印调用堆栈
	)
	for b.Loop() {
		//logger.Info("add suffix to file failed", zap.String("suffix", ".png"), zap.String("file", "path/to/file"), zap.String("error", "forbidden"))
		logger.Sugar().Infof("add suffix %s to file %s failed: %s", ".png", "path/to/file", "forbidden")
	}
}

// go test ./loggerx -bench=^BenchmarkMyLogger$ -run=^$ -timeout=10s
// go test ./loggerx -bench=^BenchmarkZap$ -run=^$ -timeout=10s

// MyLogger-10      469512              2534 ns/op
// Zap-10           432817              2752 ns/op
// Zap-10           389336              3351 ns/op	// Sugar
