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
		logger.Info("add suffix to file failed", zap.String("suffix", ".png"), zap.String("file", "path/to/file"), zap.String("error", "forbidden"))
	}
}

func BenchmarkZapSugar(b *testing.B) {
	file, err := os.OpenFile("../loggerx/log/zap_sugar.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
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
		logger.Sugar().Infof("add suffix %s to file %s failed: %s", ".png", "path/to/file", "forbidden")
	}
}

/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./loggerx -bench=^Benchmark -run=^$ -count=1 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/loggerx
cpu: Apple M1 Pro
BenchmarkMyLogger-10              453636              2552 ns/op             616 B/op          9 allocs/op
BenchmarkZap-10                   422684              2758 ns/op             529 B/op          7 allocs/op
BenchmarkZapSugar-10              408788              2956 ns/op             537 B/op          9 allocs/op
PASS
ok      github.com/yzletter/go-toolery/loggerx  4.062s
*/
