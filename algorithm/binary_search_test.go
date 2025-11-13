package algorithm_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/yzletter/go-toolery/algorithm"
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
			if ans, err := algorithm.BinarySearch(arr, target); err != nil {
				t.Fail()
				fmt.Println(1)
				break LOOP
			} else if arr[ans] != arr[j] {
				t.Fail()
				fmt.Println(arr, ans, j)
				break LOOP
			}
		}
	}
}
