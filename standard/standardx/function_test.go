package standardx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/standard/standardx"
)

func TestTernary(t *testing.T) {
	num1 := 1
	num2 := 3
	a := standardx.Ternary(num1 == num2, 5, 6)
	fmt.Println(a)
}

// go test -v ./standard/standardx -run=^TestTernary$ -count=1
