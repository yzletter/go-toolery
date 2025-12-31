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

const T = 1000

func initArr(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = rand.Int()

	}
	return res
}

func BenchmarkQuickSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for t := 0; t <= T; t++ {
			arr := initArr(10000)
			algorithmx.QuickSort(arr)
		}
	}
}

func BenchmarkSliceSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for t := 0; t <= T; t++ {
			arr := initArr(10000)
			sort.Slice(arr, func(i, j int) bool {
				return arr[i] < arr[j]
			})
		}
	}
}

// go test ./algorithmx -bench=Benchmark -run=^$ -benchmem
/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./algorithmx -bench=Benchmark -run=^$ -benchmem
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/algorithmx
cpu: Apple M1 Pro
BenchmarkQuickSort-10                  2         560753667 ns/op        82002064 B/op       1002 allocs/op
BenchmarkSliceSort-10                  2         819695230 ns/op        82058080 B/op       3004 allocs/op
PASS
ok      github.com/yzletter/go-toolery/algorithmx       4.738s
*/
