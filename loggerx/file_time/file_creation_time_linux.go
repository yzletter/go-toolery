//go:build linux
// +build linux

package file_time

import (
	"time"
)

// 说明：Linux 上多数文件系统没有“创建时间”的统一暴露接口，
// Stat_t 里的 Ctim 是“状态变更时间”，不是创建时间。
// 勉强用 ctime 代替
func fileCreationTime(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, nil
	}
	stat := fi.Sys().(*syscall.Stat_t)
	sec := stat.Ctim.Sec
	nsec := stat.Ctim.Nsec
	return time.Unix(sec, nsec), nil

}
