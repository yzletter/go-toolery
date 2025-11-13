package algorithmx_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/yzletter/go-toolery/algorithmx"
)

func TestQuickSort(t *testing.T) {
	const L = 200
LOOP:
	for i := 0; i < 1000; i++ {
		arr1, arr2 := make([]int, L), make([]int, L)

		for j := 0; j < L; j++ {
			arr1[j] = rand.Intn(L)
		}
		copy(arr2, arr1)

		sort.Slice(arr1, func(i, j int) bool {
			return arr1[i] < arr1[j]
		})

		algorithmx.QuickSort(arr2)

		for j := 0; j < L; j++ {
			if arr1[j] != arr2[j] {
				t.Fail()
				break LOOP
			}
		}
	}
}

// go test -v ./algorithmx -run=^TestQuickSort$ -count=1
