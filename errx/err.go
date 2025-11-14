package errx

import "errors"

var (
	ErrInvalidParam = errors.New("PKCS7 : UnPadding 传入非法参数")
)

var (
	ErrPriorityQueueEmpty = errors.New("PriorityQueue : 非法访问空堆 ")
)
