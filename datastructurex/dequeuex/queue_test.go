package dequeuex_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/datastructurex/dequeuex"
)

func TestDequeue(t *testing.T) {
	q := dequeuex.NewDequeue[int]()
	q.PushFront(1) // 1
	q.PushBack(2)  // 1 2
	q.PushFront(3) // 3 1 2
	q.PushBack(4)  // 3 1 2 4
	q.PushFront(9) // 9 3 1 2 4

	for q.Size() > 0 {
		val, err := q.Front()
		if err != nil {
			t.Fail()
		}
		fmt.Printf("%v ", val)

		err = q.PopFront()
		if err != nil {
			t.Fail()
		}
	}
}

// go test -v ./data_structure/dequeuex -run=^TestDequeue$ -count=1
