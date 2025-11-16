package dequeuex

import (
	"github.com/yzletter/go-toolery/errx"
)

// Dequeue 双端队列
type Dequeue[T any] struct {
	Data   []T
	Length int
}

// NewDequeue 构造函数
func NewDequeue[T any]() *Dequeue[T] {
	return &Dequeue[T]{
		Data:   make([]T, 0),
		Length: 0,
	}
}

// PushBack 插到队尾
func (dq *Dequeue[T]) PushBack(val T) {
	dq.Data = append(dq.Data, val)
	dq.Length++
}

// PushFront 插到队头
func (dq *Dequeue[T]) PushFront(val T) {
	dq.Data = append(append([]T(nil), val), dq.Data...)
	dq.Length++
}

// PopFront 弹出队头, 返回可能的错误
func (dq *Dequeue[T]) PopFront() error {
	if dq.Size() == 0 {
		return errx.ErrDequeueEmpty
	}
	dq.Data = dq.Data[1:]
	dq.Length--
	return nil
}

// PopBack 弹出队尾
func (dq *Dequeue[T]) PopBack() error {
	if dq.Size() == 0 {
		return errx.ErrDequeueEmpty
	}
	dq.Data = dq.Data[:dq.Size()-1]
	dq.Length--
	return nil
}

// Front 取队头
func (dq *Dequeue[T]) Front() (T, error) {
	if dq.Size() == 0 {
		var t T
		return t, errx.ErrDequeueEmpty
	}
	return dq.Data[0], nil
}

// Back 取队尾
func (dq *Dequeue[T]) Back() (T, error) {
	if dq.Size() == 0 {
		var t T
		return t, errx.ErrDequeueEmpty
	}
	return dq.Data[dq.Size()-1], nil
}

// Size 返回队列长度
func (dq *Dequeue[T]) Size() int {
	return dq.Length
}
