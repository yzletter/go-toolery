package errx

import "errors"

// 非法访问错误
var (
	ErrStackEmpty         = errors.New("go-toolery stackx : 非法访问空栈 ")
	ErrDequeueEmpty       = errors.New("go-toolery queuex : 非法访问空队 ")
	ErrPriorityQueueEmpty = errors.New("go-toolery priorityqueuex : 非法访问空堆 ")
	ErrLinkedListEmpty    = errors.New("go-toolery listx : 非法访问空链表 ")
)

// 其他错误处理
var (
	ErrNilSlice               = errors.New("go-toolery utilx : Jaccard 传入切片为空")
	ErrPKCS7InvalidParam      = errors.New("go-toolery utilx : PKCS7 UnPadding 传入非法参数")
	ErrLinkedListInvalidParam = errors.New("go-toolery utilx : LinkedList 传入非法参数")
)
