package setx_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/yzletter/go-toolery/setx"
)

func TestSet(t *testing.T) {
	hash := setx.NewSet[int]()
	for i := 0; i < 100; i++ {
		hash.Insert(rand.Intn(10))
	}

	vals := hash.Values()
	fmt.Println(vals)

	idx := 0
	for hash.Size() > 0 {
		hash.Delete(vals[idx])
		fmt.Println(hash.Values())
		idx++
	}
}

// go test -v ./data_structure/setx -run=^TestSet$ -count=1
