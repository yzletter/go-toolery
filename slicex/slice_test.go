package slicex_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/slicex"
)

func TestSlice(t *testing.T) {
	arr := []int{0, 1, 0, 1, 5, 6, 4, 2, 3, 5}
	fmt.Println(slicex.Unique(arr))
}

// go test -v ./standard/slicex -run=^TestSlice$ -count=1
