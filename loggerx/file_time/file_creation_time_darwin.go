//go:build darwin
// +build darwin

package file_time

import (
	"os"
	"syscall"
	"time"
)

// macOS 下 Stat_t 里有 Birthtimespec 字段
func fileCreationTime(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, syscall.EINVAL
	}

	sec := stat.Birthtimespec.Sec
	nsec := stat.Birthtimespec.Nsec

	return time.Unix(sec, nsec), nil
}
