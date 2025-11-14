package listx

import "cmp"

// ListNode 链表节点
type ListNode[T cmp.Ordered] struct {
	Val  T            // 节点携带数据
	Prev *ListNode[T] // 前驱节点
	Next *ListNode[T] // 后继节点
}

// NewListNode 根据 Val 生成一个新 ListNode
func NewListNode[T cmp.Ordered](val T) *ListNode[T] {
	return &ListNode[T]{Val: val}
}

// InsertNodeBefore 在 b 节点前插入一个 a 节点
func InsertNodeBefore[T cmp.Ordered](a, b *ListNode[T]) {
	a.Next = b
	a.Prev = b.Prev
	b.Prev.Next = a
	b.Prev = a
}

// DeleteNode 删除节点 a
func DeleteNode[T cmp.Ordered](a *ListNode[T]) {
	a.Prev.Next = a.Next
	a.Next.Prev = a.Prev
	a.Next = nil
	a.Prev = nil
}
