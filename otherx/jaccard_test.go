package otherx_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/yzletter/go-toolery/otherx"
)

func TestJaccard(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器"}
	fmt.Println(otherx.Jaccard(l1, l2))
}

func TestJaccardForSorted(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器", "Dota2"}
	slices.Sort(l1)
	slices.Sort(l2)
	fmt.Println(otherx.JaccardForSorted(l1, l2))
}

// go test -v ./otherx -run=^TestJaccard$ -count=1
// go test -v ./otherx -run=^TestJaccardForSorted$ -count=1
