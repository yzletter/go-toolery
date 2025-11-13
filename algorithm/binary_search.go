package algorithm

import (
	"cmp"

	"github.com/yzletter/go-toolery/errs"
)

// 返回 arr 中第一个大于等于 target 的元素的下标, 若存在返回 下标, nil, 若不存在返回 -1 和 error
func BinarySearch[T cmp.Ordered](arr []T, target T) (int, error) {
	l, r := 0, len(arr)-1
	for l < r {
		mid := (l + r) >> 1
		if arr[mid] < target {
			l = mid + 1
		} else {
			r = mid
		}
	}

	if arr[l] >= target {
		return l, nil
	} else {
		return -1, errs.ErrNotFound
	}
}
