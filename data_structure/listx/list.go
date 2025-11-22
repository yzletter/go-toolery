package listx

import (
	"cmp"

	"github.com/yzletter/go-toolery/errx"
)

// LinkedList 带头节点的双向循环链表
type LinkedList[T cmp.Ordered] struct {
	Head   *ListNode[T] // 头节点
	Length int          // 链表长度
}

// NewLinkedList 构造函数
func NewLinkedList[T cmp.Ordered]() *LinkedList[T] {
	// 初始化头节点
	head := &ListNode[T]{}
	head.Next = head
	head.Prev = head

	return &LinkedList[T]{
		Head:   head,
		Length: 0,
	}
}

// NewLinkedListFromSlice 将 slice 转为链表
func NewLinkedListFromSlice[T cmp.Ordered](src []T) *LinkedList[T] {
	head := &ListNode[T]{}
	head.Next = head
	head.Prev = head
	for i := 0; i < len(src); i++ {
		node := NewListNode(src[i])
		InsertNodeBefore(node, head) // 将该节点插到头节点前面
	}
	return &LinkedList[T]{
		Head:   head,
		Length: len(src),
	}
}

// Traverse 正序遍历整个 LinkedList, 传入一个 operate 函数对节点进行操作
func (list *LinkedList[T]) Traverse(operate func(listNode *ListNode[T])) {
	nowNode := list.Head.Next
	for nowNode != list.Head {
		operate(nowNode) // 操作当前遍历到的节点
		nowNode = nowNode.Next
	}
}

// ReverseTraverse 逆序遍历整个 LinkedList, 传入一个 operate 函数对节点进行操作
func (list *LinkedList[T]) ReverseTraverse(operate func(listNode *ListNode[T])) {
	nowNode := list.Head.Prev
	for nowNode != list.Head {
		operate(nowNode) // 操作当前遍历到的节点
		nowNode = nowNode.Prev
	}
}

// InsertToHead 头插法
func (list *LinkedList[T]) InsertToHead(val T) {
	node := NewListNode(val)               // 初始化新节点
	InsertNodeBefore(node, list.Head.Next) // 将该节点插到头节点后面
	list.Length += 1                       // 修改链表长度
}

// InsertToTail 尾插法
func (list *LinkedList[T]) InsertToTail(val T) {
	node := NewListNode(val)          // 初始化新节点
	InsertNodeBefore(node, list.Head) // 将该节点插到头节点前面
	list.Length += 1                  // 修改链表长度
}

// InsertBefore 在 node 前面添加一个数据为 val 的新节点
func (list *LinkedList[T]) InsertBefore(val T, node *ListNode[T]) {
	newNode := NewListNode(val)
	InsertNodeBefore(newNode, node)
	list.Length += 1
}

// InsertAfter 在 node 后面添加一个数据为 val 的新节点
func (list *LinkedList[T]) InsertAfter(val T, node *ListNode[T]) {
	newNode := NewListNode(val)
	InsertNodeBefore(newNode, node.Next)
	list.Length += 1
}

// FindNode 寻找从 0 开始下标为 idx 的节点, 返回节点和可能的错误
func (list *LinkedList[T]) FindNode(idx int) (node *ListNode[T], err error) {
	if idx < 0 || idx >= list.Length { // 判断下标是否合法
		return nil, errx.ErrLinkedListInvalidParam
	}
	// 遍历节点
	nowNode := list.Head.Next
	for i := 0; i < idx; i++ {
		nowNode = nowNode.Next
	}
	return nowNode, nil
}

// Values 将链表值转为切片
func (list *LinkedList[T]) Values() []T {
	arr := make([]T, list.Length)

	nowNode := list.Head.Next
	for i := 0; i < list.Length; i++ {
		arr[i] = nowNode.Val
		nowNode = nowNode.Next
	}

	return arr
}

// LastNode 返回最后一个节点和可能的错误
func (list *LinkedList[T]) LastNode() (*ListNode[T], error) {
	if list.Length == 0 {
		return nil, errx.ErrLinkedListEmpty
	}
	return list.Head.Prev, nil
}

// FirstNode 返回第一个节点和可能的错误
func (list *LinkedList[T]) FirstNode() (*ListNode[T], error) {
	if list.Length == 0 {
		return nil, errx.ErrLinkedListEmpty
	}
	return list.Head.Next, nil
}

// todo
// 翻转链表（非原地）
// 翻转链表 (原地）
// 有序链表去重
// 无序链表去重
// 链表排序
// 删除指定节点
// 删除下标为 idx 的节点
// 将节点插入到 idx 的位置
