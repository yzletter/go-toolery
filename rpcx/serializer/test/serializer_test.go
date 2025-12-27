package serializer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yzletter/go-toolery/rpcx"
	"github.com/yzletter/go-toolery/rpcx/serializer"
)

var myData = rpcx.RpcxData{
	A:  123,
	B:  3.1415,
	C:  false,
	D:  3.14159265333333,
	E:  "yzletter",
	Id: "yzletter",
}

func TestMySerializer(t *testing.T) {
	s := serializer.MySerializer{}
	bs, err := s.Marshal(myData)
	if err != nil {
		fmt.Println("MySerializer 序列化失败")
		t.Fail()
	} else {
		fmt.Println("MySerializer 序列化成功")
		var target rpcx.RpcxData
		err = s.Unmarshal(bs, &target)
		if err != nil {
			fmt.Println("MySerializer 反序列化失败")
			t.Fail()
		} else {
			fmt.Println("MySerializer 反序列化成功")
			fmt.Printf("反序列化结果%v\n", target)
			// 校验结果
			assert.Equal(t, myData, target, "MySerializer 反序列化后跟原始值不同")
		}
	}
}

// 使用 Bytedance/sonic
func BenchmarkBytedance(b *testing.B) {
	s := serializer.JsonByBytedanceSonic{}
	var target rpcx.RpcxData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(myData)
		s.Unmarshal(stream, &target)
	}
}

// 使用 encoding/gob
func BenchmarkGob(b *testing.B) {
	s := serializer.Gob{}
	var target rpcx.RpcxData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(myData)
		s.Unmarshal(stream, &target)
	}
}

// 使用 MySerializer
func BenchmarkMySerializer(b *testing.B) {
	s := serializer.MySerializer{}

	var target rpcx.RpcxData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(myData)
		s.Unmarshal(stream, &target)
	}
}

/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./rpcx/serializer/test -bench=^Benchmark -run=^$ -count=1 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/rpcx/serializer/test
cpu: Apple M1 Pro
BenchmarkBytedance-10            1680181               701.3 ns/op           523 B/op          6 allocs/op
BenchmarkGob-10                   115826              9968 ns/op            9048 B/op        187 allocs/op
BenchmarkMySerializer-10          671557              1749 ns/op            1816 B/op         60 allocs/op
PASS
ok      github.com/yzletter/go-toolery/rpcx/serializer/test     4.970s
*/
