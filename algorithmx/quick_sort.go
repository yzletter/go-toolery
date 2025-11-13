package algorithmx

import (
	"cmp"
)

// 原地快速排序
func QuickSort[T cmp.Ordered](arr []T) {
	var internal_quick_sort func(arr []T, l, r int)
	internal_quick_sort =
		func(arr []T, l, r int) {
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

			internal_quick_sort(arr, l, j)
			internal_quick_sort(arr, j+1, r)
		}
	internal_quick_sort(arr, 0, len(arr)-1)
}
