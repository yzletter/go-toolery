package mathx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/standard/mathx"
)

func TestQMI(t *testing.T) {
	fmt.Println(mathx.QMI(2, 10, 1e9+7))
}

// go test -v ./standard/mathx -run=^TestQMI$ -count=1
