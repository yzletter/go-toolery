package loggerx

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yzletter/go-toolery/loggerx/fileutil"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Log struct {
	logger     *log.Logger // 基本库日志
	logFile    string      // 日志输出文件名
	logOut     *os.File    // 日志输出文件句柄
	logLevel   int         // 日志等级
	writeLock  sync.Mutex  // 锁
	day        int         // 创建时间
	rotateLock sync.Mutex
}

func NewLog(logFile string, logLevel int) *Log {
	if logOut, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666); err == nil {
		return &Log{
			logger:   log.New(logOut, "", log.Ldate|log.Lmicroseconds),
			logFile:  logFile,
			logOut:   logOut,
			logLevel: logLevel,
			day:      time.Now().YearDay(),
		}
	} else {
		_, _ = os.Stderr.WriteString("open file failed")
		return nil
	}
}

func (l *Log) print(level int, content string) {
	var prefix string

	switch level {
	case DebugLevel:
		prefix = "[DEBUG] "
	case WarnLevel:
		prefix = "[WARN] "
	case InfoLevel:
		prefix = "[INFO] "
	case ErrorLevel:
		prefix = "[ERROR] "
	}

	msg := prefix + " " + currentStack() + " " + content
	l.logger.Println(msg)
}

func (l *Log) Debug(content string) {
	if l.logLevel <= DebugLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		l.print(DebugLevel, content)
	}

}

func (l *Log) Debugf(format string, v ...any) {
	if l.logLevel <= DebugLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		l.print(DebugLevel, content)
	}
}

func (l *Log) Warn(content string) {
	if l.logLevel <= WarnLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		l.print(WarnLevel, content)
	}
}

func (l *Log) Warnf(format string, v ...any) {
	if l.logLevel <= WarnLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		l.print(WarnLevel, content)
	}
}

func (l *Log) Info(content string) {
	if l.logLevel <= InfoLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		l.print(InfoLevel, content)
	}
}

func (l *Log) Infof(format string, v ...any) {
	if l.logLevel <= InfoLevel {
		l.writeLock.Lock()

		defer l.writeLock.Unlock()
		content := fmt.Sprintf(format, v...)
		l.print(InfoLevel, content)
	}
}

func (l *Log) Error(content string) {
	if l.logLevel <= ErrorLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		l.print(ErrorLevel, content)

		for _, stk := range stackPath() {
			_, _ = l.logOut.WriteString("\t" + stk + "\n")
		}
	}
}

func (l *Log) Errorf(format string, v ...any) {
	if l.logLevel <= ErrorLevel {
		l.writeLock.Lock()
		defer l.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		l.print(ErrorLevel, content)
		for _, stk := range stackPath() {
			_, _ = l.logOut.WriteString("\t" + stk + "\n")
		}
	}
}

// 判断是否进行滚动
func (l *Log) rotate() {
	// 判断是否过了一天
	now := time.Now()
	if now.YearDay() == l.day {
		return
	}

	l.rotateLock.Lock()
	defer l.rotateLock.Unlock()

	// 再判断一次, 解决并发问题
	if now.YearDay() == l.day {
		return
	}

	// 进行滚动

	// 关闭旧文件
	_ = l.logOut.Close()

	// 给老文件加上日期后缀
	createTime, _ := fileutil.FileCreationTime(l.logFile)
	postFix := createTime.Format("20060102")

	if err := os.Rename(l.logFile, l.logFile+"."+postFix); err != nil {
		// 如果logger本身出错，则把错误信息打到标准错误输出里
		_, _ = os.Stderr.WriteString(fmt.Sprintf("append date postfix %s to l logOut %s failed: %v\n", postFix, l.logFile, err))
		return
	}

	// 再打开一个新的文件句柄, 进行替换
	if logOut, err := os.OpenFile(l.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("create l logOut %s failed %v\n", l.logFile, err))
		return
	} else {
		l.logger = log.New(logOut, "", log.Ldate|log.Lmicroseconds)
		l.logOut = logOut
		l.day = now.YearDay() // 更新 day 值为今天
	}
}
