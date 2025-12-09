//go:build windows
// +build windows

package file_time

import (
	"os"
	"syscall"
	"time"
)

// Windows 下通过 Win32FileAttributeData.CreationTime 获取
func fileCreationTime(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	stat, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, syscall.EINVAL
	}

	// CreationTime 是 FILETIME，单位是 100 纳秒，从 1601-01-01 UTC 起算
	ft := stat.CreationTime
	return time.Unix(0, ft.Nanoseconds()), nil
}
