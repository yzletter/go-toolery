package slicex

import (
	"github.com/yzletter/go-toolery/setx"
)

// Unique 借助 set 去重, 返回无序的去重切片
func Unique[T comparable](target []T) []T {
	hash := setx.NewSet[T]()
	for _, v := range target {
		hash.Insert(v)
	}
	return hash.Values()
}
