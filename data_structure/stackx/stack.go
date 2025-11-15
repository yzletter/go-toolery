package stackx

import "github.com/yzletter/go-toolery/errx"

// todo 栈

// Stack 栈
type Stack[T any] struct {
	Data []T
}

// NewStack 构造函数
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		Data: make([]T, 0),
	}
}

// Top 取栈顶元素, 返回栈顶元素和可能的错误
func (stk *Stack[T]) Top() (T, error) {
	if stk.Size() == 0 {
		var t T
		return t, errx.ErrStackEmpty
	}

	return stk.Data[stk.Size()-1], nil
}

// Pop 弹出栈顶, 返回可能的错误
func (stk *Stack[T]) Pop() error {
	if stk.Size() == 0 {
		return errx.ErrStackEmpty
	}
	stk.Data = stk.Data[:stk.Size()-1]
	return nil
}

// Push 新元素入栈
func (stk *Stack[T]) Push(val T) {
	stk.Data = append(stk.Data, val)
}

// Size 返回栈的长度
func (stk *Stack[T]) Size() int {
	return len(stk.Data)
}
