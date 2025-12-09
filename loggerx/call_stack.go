package loggerx

import (
	"fmt"
	"runtime"
	"strings"
)

// 返回调用当前函数的堆栈信息
func currentStack() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d", shortFile(file), line)
	} else {
		return ""
	}
}

// 避免路径过长, 只保留最后三节
func shortFile(file string) string {
	arr := strings.Split(file, "/")
	if len(arr) > 3 {
		arr = arr[len(arr)-3:]
	}
	return strings.Join(arr, "/")
}

func stackPath() []string {
	rect := make([]string, 3)
	for i := 3; i < 5; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			rect = append(rect, fmt.Sprintf("%s:%d", shortFile(file), line))
		}
	}
	return rect
}
