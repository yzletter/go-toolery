package standardx_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/yzletter/go-toolery/standard/standardx"
)

func TestTernary(t *testing.T) {
	num1 := 1
	num2 := 3
	a := standardx.Ternary(num1 == num2, 5, 6)
	fmt.Println(a)
}

func TestHash(t *testing.T) {
	for i := 0; i < 10; i++ {
		n := rand.Intn(100)
		res := standardx.Hash(strconv.Itoa(n))
		fmt.Println(res)
	}
}

// go test -v ./standard/standardx -run=^TestTernary$ -count=1
// go test -v ./standard/standardx -run=^TestHash$ -count=1
