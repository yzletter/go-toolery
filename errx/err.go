package errx

import "errors"

// Otherx 错误处理
var (
	ErrInvalidParam = errors.New("PKCS7 : UnPadding 传入非法参数")
)

// PriorityQueuex 错误处理
var (
	ErrPriorityQueueEmpty = errors.New("PriorityQueue : 非法访问空堆")
)

// Stackx 错误处理
var (
	ErrStackEmpty = errors.New("Stack : 非法访问空栈 ")
)

// Queuex 错误处理
var (
	ErrDequeueEmpty = errors.New("Dequeue : 非法访问空队 ")
)
