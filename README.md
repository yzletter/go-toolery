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