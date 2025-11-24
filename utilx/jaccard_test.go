package utilx_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/yzletter/go-toolery/utilx"
)

func TestJaccard(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器"}
	if res, err := utilx.Jaccard(l1, l2); err != nil {
		t.Fail()
	} else {
		fmt.Println(res)
	}
}

func TestJaccardForSorted(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器", "Dota2"}
	slices.Sort(l1)
	slices.Sort(l2)

	if res, err := utilx.JaccardForSorted(l1, l2); err != nil {
		t.Fail()
	} else {
		fmt.Println(res)
	}
}

// go test -v ./utilx -run=^TestJaccard$ -count=1
// go test -v ./utilx -run=^TestJaccardForSorted$ -count=1
