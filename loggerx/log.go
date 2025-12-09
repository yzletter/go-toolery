package loggerx

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Log struct {
	logger    *log.Logger // 基本库日志
	logFile   string      // 日志输出文件
	logOut    *os.File    // 日志输出文件句柄
	logLevel  int         // 日志等级
	writeLock sync.Mutex  // 锁
}

func NewLog(logFile string, logLevel int) *Log {
	if file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666); err == nil {
		return &Log{
			logger:   log.New(file, "", log.Ldate|log.Lmicroseconds),
			logOut:   file,
			logLevel: logLevel,
		}
	} else {
		return nil
	}
}

func (log *Log) print(level int, content string) {
	var prefix string

	switch level {
	case DebugLevel:
		prefix = "DEBUG "
	case WarnLevel:
		prefix = "WARN "
	case InfoLevel:
		prefix = "INFO "
	case ErrorLevel:
		prefix = "ERROR "
	}
	
	msg := prefix + " " + currentStack() + " " + content
	log.logger.Println(msg)
}

func (log *Log) Debug(content string) {
	if log.logLevel <= DebugLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		log.print(DebugLevel, content)
	}

}

func (log *Log) Debugf(format string, v ...any) {
	if log.logLevel <= DebugLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		log.print(DebugLevel, content)
	}
}

func (log *Log) Warn(content string) {
	if log.logLevel <= WarnLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		log.print(WarnLevel, content)
	}
}

func (log *Log) Warnf(format string, v ...any) {
	if log.logLevel <= WarnLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		log.print(WarnLevel, content)
	}
}

func (log *Log) Info(content string) {
	if log.logLevel <= InfoLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		log.print(InfoLevel, content)
	}
}

func (log *Log) Infof(format string, v ...any) {
	if log.logLevel <= InfoLevel {
		log.writeLock.Lock()

		defer log.writeLock.Unlock()
		content := fmt.Sprintf(format, v...)
		log.print(InfoLevel, content)
	}
}

func (log *Log) Error(content string) {
	if log.logLevel <= ErrorLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		log.print(ErrorLevel, content)
		for _, stk := range stackPath() {
			_, _ = log.logOut.WriteString("\t" + stk + "\n")
		}
	}
}

func (log *Log) Errorf(format string, v ...any) {
	if log.logLevel <= ErrorLevel {
		log.writeLock.Lock()
		defer log.writeLock.Unlock()

		content := fmt.Sprintf(format, v...)
		log.print(ErrorLevel, content)
		for _, stk := range stackPath() {
			_, _ = log.logOut.WriteString("\t" + stk + "\n")
		}
	}
}
