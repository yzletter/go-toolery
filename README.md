# go-toolery

## 简介

## 算法 (Algorithm)

#### 二分查找 (BinarySearch)
```go
// 返回 arr 中第一个大于等于 target 的元素的下标, 若存在返回 下标, nil, 若不存在返回 -1 和 error
func BinarySearch[T cmp.Ordered](arr []T, target T) (int, error) 
```

