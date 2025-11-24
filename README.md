# go-toolery

## 简介 (Introduction)

个人练手用的现代化的 Go 工具库，聚焦于数据结构、算法实现、工程辅助函数与常用工具组件，在实际工程项目中快速解决常见问题，减少重复造轮子。

## 目录结构

```
./
├── README.md
├── algorithmx
│   ├── binary_search.go
│   ├── binary_search_test.go
│   ├── quick_sort.go
│   └── quick_sort_test.go
├── dequeuex
│   ├── queue.go
│   └── queue_test.go
├── errx
│   └── err.go
├── go.mod
├── listx
│   ├── list.go
│   ├── list_test.go
│   └── node.go
├── mathx
│   ├── math.go
│   └── math_test.go
├── priority_queuex
│   ├── priority_queue.go
│   └── priority_queue_test.go
├── randx
│   ├── string.go
│   └── string_test.go
├── setx
│   ├── set.go
│   └── set_test.go
├── slicex
│   ├── slice.go
│   └── slice_test.go
├── stackx
│   ├── stack.go
│   └── stack_test.go
├── standardx
│   ├── function.go
│   └── function_test.go
├── treex
│   ├── binary_tree.go
│   ├── binary_tree_test.go
│   └── node.go
└── utilx
    ├── jaccard.go
    ├── jaccard_test.go
    ├── jwt.go
    ├── jwt_test.go
    ├── pkcs7.go
    └── pkcs7_test.go
```

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

// FindNode 寻找从 0 开始下标为 idx 的节点, 返回节点和可能的错误
func (list *LinkedList[T]) FindNode(idx int) (node *ListNode[T], err error) {}

// Values 将链表值转为切片
func (list *LinkedList[T]) Values() []T {}

// LastNode 返回最后一个节点和可能的错误
func (list *LinkedList[T]) LastNode() (*ListNode[T], error) {}

// FirstNode 返回第一个节点和可能的错误
func (list *LinkedList[T]) FirstNode() (*ListNode[T], error) {}
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

```go
// Dequeue 双端队列
type Dequeue[T any] struct {
	Data   []T
	Length int
}

// NewDequeue 构造函数
func NewDequeue[T any]() *Dequeue[T] {}

// PushBack 插到队尾
func (dq *Dequeue[T]) PushBack(val T) {}

// PushFront 插到队头
func (dq *Dequeue[T]) PushFront(val T) {}

// PopFront 弹出队头, 返回可能的错误
func (dq *Dequeue[T]) PopFront() error {}

// PopBack 弹出队尾
func (dq *Dequeue[T]) PopBack() error {}

// Front 取队头
func (dq *Dequeue[T]) Front() (T, error) {}

// Back 取队尾
func (dq *Dequeue[T]) Back() (T, error) {}

// Size 返回队列长度
func (dq *Dequeue[T]) Size() int {}
```

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

### 手写二叉树 (Treex)

```go
// BinaryTree 二叉树
type BinaryTree struct {
	Root *BNode
}

// NewBinaryTree 根据 root 构造一颗二叉树
func NewBinaryTree(root *BNode) *BinaryTree {}

// PreOrder 二叉树先序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) PreOrder(operate BNodeOperationFunc) {}

// MiddleOrder 二叉树中序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) MiddleOrder(operate BNodeOperationFunc) {}

// PostOrder 二叉树后序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) PostOrder(operate BNodeOperationFunc) {}

// LevelOrder 二叉树层序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) LevelOrder(operate BNodeOperationFunc) {}
```

## 标准库辅助

### Slice 辅助 (Slicex)

```go
// Unique 借助 set 去重, 返回无序的去重切片
func Unique[T comparable](target []T) []T {}
```

### Math 辅助 (Mathx)

```go
// QMI 快速幂求 a ^ k % p
func QMI(a, k, p int) int {}
```

### Rand 辅助 (Randx)

```go
// RandString 根据传入的种子, 生成长度为 len 的随机字符串
func RandString(seed string, length int) string {}
```

### Function 辅助 (Standardx)

```go
// Ternary 三目运算符, 传入 bool 和可能返回的两个变量
func Ternary[T any](condition bool, a, b T) T {}

// Hash 返回字符串 MD5 哈希后 32 位的十六进制编码结果
func Hash(password string) string {}
```

## 其他

### PKCS7 数据填充

```go
// Padding 将数据填充, 填充至 blockSize 的整数倍, 返回填充好的切片
func Padding(src []byte, blockSize int) []byte {}

// UnPadding 将数据还原, 返回填充好的切片和可能存在的错误
func UnPadding(src []byte, blockSize int) ([]byte, error) {}
```

### Jaccard 相似度

```go
// Jaccard 计算相似度 = 交集 / 并集
func Jaccard[T comparable](collection1, collection2 []T) (float64, error) {}

// JaccardForSorted 计算有序集合的相似度
func JaccardForSorted[T cmp.Ordered](collection1, collection2 []T) (float64, error) {}
```

### LRU 