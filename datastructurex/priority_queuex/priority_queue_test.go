package priority_queuex_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/yzletter/go-toolery/datastructurex/priority_queuex"
)

func TestPriorityQueue(t *testing.T) {
	// 创建小根堆
	heap1 := priority_queuex.NewPriorityQueue(func(a, b int) bool {
		return a < b
	})
	// 创建大根堆
	heap2 := priority_queuex.NewPriorityQueue(func(a, b int) bool {
		return a > b
	})

	// 随机插入 20 个数
	for i := 0; i < 20; i++ {
		heap1.Push(rand.Intn(100))
		heap2.Push(rand.Intn(100))
	}

	for heap1.Size() > 0 {
		if ele, err := heap1.Top(); err != nil {
			t.Fail()
		} else {
			fmt.Printf("%v ", ele)
		}
		if err := heap1.Pop(); err != nil {
			t.Fail()
		}
	}

	fmt.Println()
	for heap2.Size() > 0 {
		if ele, err := heap2.Top(); err != nil {
			t.Fail()
		} else {
			fmt.Printf("%v ", ele)
		}
		if err := heap2.Pop(); err != nil {
			t.Fail()
		}
	}

}

// go test -v ./data_structure/priority_queuex -run=^TestPriorityQueue$ -count=5
