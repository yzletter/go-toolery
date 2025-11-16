# go-toolery

## 简介 (Introduction)

个人练手用的现代化的 Go 工具库，聚焦于数据结构、算法实现、工程辅助函数与常用工具组件，在实际工程项目中快速解决常见问题，减少重复造轮子。

## 算法 (Algorithm)

### 手写二分查找 (BinarySearch)
```go
// 返回 arr 中第一个大于等于 target 的元素的下标, 返回 下标位置, 是否找到
func BinarySearch[T cmp.Ordered](arr []T, target T) (idx int, ok bool) {}

```

### 手写快排 (QuickSort)

```go
// 原地快速排序
func QuickSort[T cmp.Ordered](arr []T) {}

```

## 数据结构 (DataStructure)

### 手写带头结点的双向循环链表 (Listx)

```go
// DoubleList 带头节点的双向循环链表
type DoubleList[T cmp.Ordered] struct {
    Head   *ListNode[T] // 头节点
    Length int          // 链表长度
}

// NewDoubleList 构造函数
func NewDoubleList[T cmp.Ordered]() *DoubleList[T] {}

// NewDoubleListFromSlice 将 slice 转为链表
func NewDoubleListFromSlice[T cmp.Ordered](src []T) *DoubleList[T] {}

// Traverse 正序遍历整个 DoubleList, 传入一个 operate 函数对节点进行操作
func (list *DoubleList[T]) Traverse(operate func(listNode *ListNode[T])) {}

// ReverseTraverse 逆序遍历整个 DoubleList, 传入一个 operate 函数对节点进行操作
func (list *DoubleList[T]) ReverseTraverse(operate func(listNode *ListNode[T])) {}

// InsertToHead 头插法
func (list *DoubleList[T]) InsertToHead(val T) {}

// InsertToTail 尾插法
func (list *DoubleList[T]) InsertToTail(val T) {}

// InsertBefore 在 node 前面添加一个数据为 val 的新节点
func (list *DoubleList[T]) InsertBefore(val T, node *ListNode[T]) {}

// InsertAfter 在 node 后面添加一个数据为 val 的新节点
func (list *DoubleList[T]) InsertAfter(val T, node *ListNode[T]) {}

// FindNode 寻找从 0 开始下标为 idx 的节点, 若不存在则返回 nil
func (list *DoubleList[T]) FindNode(idx int) (node *ListNode[T]) {}

// Values 将链表值转为切片
func (list *DoubleList[T]) Values() []T {}

// LastNode 返回最后一个节点
func (list *DoubleList[T]) LastNode() *ListNode[T] {}

```

### 手写集合 (Setx)

```go
// EmptyStruct 空结构体
type EmptyStruct struct{}

// Set 集合
type Set[T comparable] map[T]EmptyStruct

// NewSet 构造函数
func NewSet[T comparable]() Set[T] {}

// Insert 插入一个新元素
func (hash Set[T]) Insert(val T) {}

// Delete 删除一个元素
func (hash Set[T]) Delete(val T) {}

// Exist 查询一个元素
func (hash Set[T]) Exist(val T) bool {}

// Size 返回集合中元素个数
func (hash Set[T]) Size() int {}

// Values 返回集合中所有元素的切片
func (hash Set[T]) Values() []T {}

```

### 手写栈 (Stackx)

```go
// Stack 栈
type Stack[T any] struct {
	Data []T
}

// NewStack 构造函数
func NewStack[T any]() *Stack[T] {}

// Top 取栈顶元素, 返回栈顶元素和可能的错误
func (stk *Stack[T]) Top() (T, error) {}

// Pop 弹出栈顶, 返回可能的错误
func (stk *Stack[T]) Pop() error {}

// Push 新元素入栈
func (stk *Stack[T]) Push(val T) {}

// Size 返回栈的长度
func (stk *Stack[T]) Size() int {}

```

### 手写双端队列 (DeQueuex)

### 手写堆 (PriorityQueuex)

```go
// PriorityQueue 堆
type PriorityQueue[T cmp.Ordered] struct {
    Data        []T
    compareFunc func(a, b T) bool
}

// NewPriorityQueue 构造一个空堆, 传入比较函数
func NewPriorityQueue[T cmp.Ordered](compareFunc func(a, b T) bool) *PriorityQueue[T] {}

// Push 新元素入堆
func (heap *PriorityQueue[T]) Push(val T) {}

// Pop 弹出堆顶, 返回可能的错误
func (heap *PriorityQueue[T]) Pop() error {}

// Top 返回堆顶元素和可能的错误
func (heap *PriorityQueue[T]) Top() (T, error) {}

// Size 返回堆的大小
func (heap *PriorityQueue[T]) Size() int {}

// 将堆底元素向上更新
func (heap *PriorityQueue[T]) pushUp() {}

// 将堆顶元素向下更新
func (heap *PriorityQueue[T]) pushDown() {}

```

### 手写并发安全 map (Mapx)

### 手写二叉树 (Treex)

```go
// BinaryTree 二叉树
type BinaryTree struct {
	Root *BNode
}

// NewBinaryTree 根据 root 构造一颗二叉树
func NewBinaryTree(root *BNode) *BinaryTree {}

// PreOrder 二叉树先序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) PreOrder(operate func(node *BNode)) {}

// MiddleOrder 二叉树中序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) MiddleOrder(operate func(node *BNode)) {}

// PostOrder 二叉树后序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) PostOrder(operate func(node *BNode)) {}

```

## 标准库辅助

### Slice 辅助 (Slicex)

### Math 辅助 (Mathx)

## 其他

### PKCS7 数据填充

```go
// Padding 将数据填充, 填充至 blockSize 的整数倍, 返回填充好的切片
func Padding(src []byte, blockSize int) []byte {}

// UnPadding 将数据还原, 返回填充好的切片和可能存在的错误
func UnPadding(src []byte, blockSize int) ([]byte, error) {}

```

### LRU 

### Jaccard 相似度