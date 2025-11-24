package stackx_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/yzletter/go-toolery/stackx"
)

func TestStack(t *testing.T) {
	stk := stackx.NewStack[int]()

	for i := 0; i < 20; i++ {
		t := rand.Intn(20)
		stk.Push(t)
		fmt.Printf("%v ", t)
	}

	fmt.Println()
	for stk.Size() > 0 {
		val, err := stk.Top()
		if err != nil {
			t.Fail()
		}
		fmt.Printf("%v ", val)
		if err := stk.Pop(); err != nil {
			t.Fail()
		}
	}
}

// go test -v ./data_structure/stackx -run=^TestStack$ -count=1
