# go-toolery

## 简介 (Introduction)

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

### 手写栈 (Stackx)

### 手写双端队列 (DeQueuex)

### 手写堆 (PriorityQueuex)

### 手写并发安全 map (Mapx)

## 标准库辅助

### Slice 辅助 (Slicex)

### Math 辅助 (Mathx)

## 其他

### PKCS7 数据填充

### LRU 