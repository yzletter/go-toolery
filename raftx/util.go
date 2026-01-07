package raftx

import (
	"math/rand/v2"
	"net/url"
	"path"
	"time"
)

// 返回 [timeout, 2 * timeout] 间的随机值
func randomTimeout(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return nil
	}
	duration := time.Duration(rand.Int64()) % timeout
	return time.After(duration)
}

// 把 thePath 拼接到 connString 后面
func joinUrlPath(connString, thePath string) string {
	// string 转为 URL
	u, err := url.Parse(connString)
	if err != nil {
		panic(err)
	}

	// 拼接
	u.Path = path.Join(u.Path, thePath)
	return u.String()
}
