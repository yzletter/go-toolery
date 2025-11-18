package otherx

import (
	"cmp"

	"github.com/yzletter/go-toolery/data_structure/setx"
	"github.com/yzletter/go-toolery/errx"
)

// Jaccard 计算相似度 = 交集 / 并集
func Jaccard[T comparable](collection1, collection2 []T) (float64, error) {
	if len(collection1) <= 0 || len(collection2) <= 0 {
		return 0.0, errx.ErrNilSlice
	}
	hash1 := setx.NewSet[T]()
	for _, v := range collection1 {
		hash1.Insert(v)
	}
	intersectionCnt := 0 // 交集个数
	for _, v := range collection2 {
		if hash1.Exist(v) {
			intersectionCnt++
		}
	}
	return float64(intersectionCnt) / float64(len(collection1)+len(collection2)-intersectionCnt), nil
}

// JaccardForSorted 计算有序集合的相似度
func JaccardForSorted[T cmp.Ordered](collection1, collection2 []T) (float64, error) {
	if len(collection1) <= 0 || len(collection2) <= 0 {
		return 0.0, errx.ErrNilSlice
	}
	i, j := 0, 0
	intersectionCnt := 0 // 交集个数
	for i < len(collection1) || j < len(collection2) {
		if i < len(collection1) && j < len(collection2) {
			if collection1[i] == collection2[j] {
				intersectionCnt++
				i++
				j++
			} else if collection1[i] < collection2[j] {
				i++
			} else {
				j++
			}
			continue
		}
		break
	}

	return float64(intersectionCnt) / float64(len(collection1)+len(collection2)-intersectionCnt), nil
}
