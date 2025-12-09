package loggerx

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yzletter/go-toolery/loggerx/file_time"
	udp_logger "github.com/yzletter/go-toolery/loggerx/udp_logger/producer"
)

const (
	DebugLevel = iota
	WarnLevel
	InfoLevel
	ErrorLevel
)

type Log struct {
	logger        *log.Logger // 基本库日志
	logFile       string      // 日志输出文件名
	logOut        *os.File    // 日志输出文件句柄
	logLevel      int         // 日志等级
	writeLock     sync.Mutex  // 锁
	addStackTrace bool        // 是否需要打印堆栈, 默认为 false
	udpProducer   *udp_logger.LogProducer
}

func NewLog(logFile string, logLevel int) *Log {
	var l *Log
	l = &Log{
		logFile:       logFile,
		logLevel:      logLevel,
		addStackTrace: false,
	}

	l.rotate()                     // 开启定时前, 先执行一次滚动
	go Schedule(l.rotate, 0, 0, 0) // 开启定时任务
	return l
}

// AddStackTrace 日志会打印三层调用堆栈
func (l *Log) AddStackTrace() {
	l.addStackTrace = true
}

// SetUDPProducer 设置 UDP Producer
func (l *Log) SetUDPProducer(collectorAddr string) {
	var err error
	l.udpProducer, err = udp_logger.NewLogProducer(collectorAddr, 4096)
	if err != nil {
		l.Errorf("SetUDPProducer failed %s", err)
	}
}

func (l *Log) Close() {
	if l.logOut != nil {
		l.logOut.Close()
	}

	if l.udpProducer != nil {
		l.udpProducer.Close()
	}
}

// 滚动
func (l *Log) rotate() {
	// 如果是初始化
	if l.logOut == nil {

		// 使用终端输出
		if l.logFile == "" {
			l.logOut = os.Stdout
			l.logger = log.New(l.logOut, "", log.Ldate|log.Lmicroseconds)
			return
		}

		// 使用指定文件
		logOut, err := os.OpenFile(l.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("create l logOut %s failed %v\n", l.logFile, err))
			return
		}
		l.logOut = logOut
		l.logger = log.New(l.logOut, "", log.Ldate|log.Lmicroseconds)
		return
	}

	// 当前时间
	now := time.Now()

	if l.logFile == "" {
		l.logOut = os.Stdout
	} else {
		// 获取旧文件创建时间
		createTime, _ := file_time.FileCreationTime(l.logFile)
		if createTime.Year() != now.Year() || createTime.YearDay() != now.YearDay() {
			postFix := createTime.Format("20060102")

			if err := os.Rename(l.logFile, l.logFile+"."+postFix); err != nil {
				// 如果logger本身出错，则把错误信息打到标准错误输出里
				_, _ = os.Stderr.WriteString(fmt.Sprintf("append date postfix %s to l logOut %s failed: %v\n", postFix, l.logFile, err))
				return
			}

			// 打开一个新的文件句柄, 进行替换
			logOut, err := os.OpenFile(l.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
			if err != nil {
				_, _ = os.Stderr.WriteString(fmt.Sprintf("create l logOut %s failed %v\n", l.logFile, err))
				return
			}

			l.logOut = logOut
		}
	}

	l.logger = log.New(l.logOut, "", log.Ldate|log.Lmicroseconds)
	return
}

func (l *Log) print(level int, content string) {
	var prefix string
	switch level {
	case DebugLevel:
		if len(l.logFile) == 0 {
			prefix = Magenta.Print("[DEBUG]") + " "
		} else {
			prefix = "[DEBUG]"
		}
	case InfoLevel:
		if len(l.logFile) == 0 {
			prefix = Blue.Print("[INFO]") + " "
		} else {
			prefix = "[INFO]"
		}
	case WarnLevel:
		if len(l.logFile) == 0 {
			prefix = Yellow.Print("[WARN]") + " "
		} else {
			prefix = "[WARN]"
		}
	case ErrorLevel:
		if len(l.logFile) == 0 {
			prefix = Red.Print("[ERROR]") + " "
		} else {
			prefix = "[ERROR]"
		}
	}

	msg := prefix + " " + currentStack() + " " + content
	l.logger.Println(msg)

	if l.udpProducer != nil {
		l.udpProducer.Send(msg)
	}
}

func (l *Log) Debug(content string) {
	if l.logLevel <= DebugLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}
		l.print(DebugLevel, content)
	}

}

func (l *Log) Debugf(format string, v ...any) {
	if l.logLevel <= DebugLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}

		content := fmt.Sprintf(format, v...)
		l.print(DebugLevel, content)
	}
}

func (l *Log) Warn(content string) {
	if l.logLevel <= WarnLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}
		l.print(WarnLevel, content)
	}
}

func (l *Log) Warnf(format string, v ...any) {
	if l.logLevel <= WarnLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}

		content := fmt.Sprintf(format, v...)
		l.print(WarnLevel, content)
	}
}

func (l *Log) Info(content string) {
	if l.logLevel <= InfoLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}

		l.print(InfoLevel, content)
	}
}

func (l *Log) Infof(format string, v ...any) {
	if l.logLevel <= InfoLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}
		content := fmt.Sprintf(format, v...)
		l.print(InfoLevel, content)
	}
}

func (l *Log) Error(content string) {
	if l.logLevel <= ErrorLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}

		l.print(ErrorLevel, content)

		if l.addStackTrace {
			for _, stk := range stackPath() {
				_, _ = l.logOut.WriteString("\t" + stk + "\n")
			}
		}
	}
}

func (l *Log) Errorf(format string, v ...any) {
	if l.logLevel <= ErrorLevel {
		if l.addStackTrace {
			l.writeLock.Lock()
			defer l.writeLock.Unlock()
		}

		content := fmt.Sprintf(format, v...)
		l.print(ErrorLevel, content)

		if l.addStackTrace {
			for _, stk := range stackPath() {
				_, _ = l.logOut.WriteString("\t" + stk + "\n")
			}
		}
	}
}
