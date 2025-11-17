package randx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/standard/randx"
)

func TestString(t *testing.T) {
	seed := "abcdefghijklmn"
	len := 10

	for i := 0; i < 10; i++ {
		fmt.Println(randx.RandString(seed, len))
	}
}

// go test -v ./standard/randx -run=^TestString$ -count=1
