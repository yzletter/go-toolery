package utilx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/utilx"
)

func TestPKCS7(t *testing.T) {
	// 测试用例
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(utilx.Padding(src, 8))
	fmt.Println(utilx.UnPadding(utilx.Padding(src, 8), 8))
	src = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Println(utilx.Padding(src, 8))
	fmt.Println(utilx.UnPadding(utilx.Padding(src, 8), 8))
}

// go test -v ./utilx -run=^TestPKCS7$ -count=1
