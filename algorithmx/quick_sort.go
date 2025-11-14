package algorithmx

import (
	"cmp"
)

// QuickSort 原地快速排序
func QuickSort[T cmp.Ordered](arr []T) {
	internalQuickSort(arr, 0, len(arr)-1)
}

func internalQuickSort[T cmp.Ordered](arr []T, l, r int) {
	if l >= r {
		return
	}
	x, i, j := arr[(l+r)>>1], l-1, r+1
	for i < j {
		for {
			i++
			if arr[i] >= x {
				break
			}
		}

		for {
			j--
			if arr[j] <= x {
				break
			}
		}

		if i < j {
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	internalQuickSort[T](arr, l, j)
	internalQuickSort[T](arr, j+1, r)
}
