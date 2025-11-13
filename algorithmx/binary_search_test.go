package algorithmx_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/yzletter/go-toolery/algorithmx"
)

func TestBoundSearch(t *testing.T) {
	const L = 100 // 测试数组长度
LOOP:
	for cnt := 0; cnt < 1000; cnt++ { // 测试 case 数
		arr := make([]int, L)
		for j := 0; j < L; j++ {
			arr[j] = rand.Intn(100) // 随机生成 0 ~ 99 的整数
		}

		// 排序
		sort.Slice(arr, func(i, j int) bool {
			return arr[i] < arr[j]
		})

		for j := 0; j < L; j++ {
			target := arr[j]
			idx, ok := algorithmx.BinarySearch(arr, target)
			// 未找到或者找到的数据不对
			if ok == false || arr[idx] < arr[j] {
				t.Fail()
				break LOOP
			}
		}
	}
}

// go test -v ./algorithmx -run=^TestBoundSearch$ -count=1
