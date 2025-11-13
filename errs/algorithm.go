package errs

import "errors"

var (
	ErrNotFound = errors.New("algorithm : 未找到目标元素")
)
