package utilx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/utilx"
)

func TestSlug(t *testing.T) {
	s := "Golang学习"
	fmt.Println(utilx.Slugify(s))

	s = "Golang*学习"
	fmt.Println(utilx.Slugify(s))

	s = "go*Lang学习"
	fmt.Println(utilx.Slugify(s))

	s = "golang*学习"
	fmt.Println(utilx.Slugify(s))
}

// go test -v ./utilx -run=^TestSlug$ -count=1
