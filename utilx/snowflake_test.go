package utilx_test

import (
	"sync"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/yzletter/go-toolery/utilx"
)

func TestSnowflake(t *testing.T) {
	const P = 10
	const LOOP = 100000
	idChannel := make(chan int64, P*LOOP)
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < LOOP; j++ {
				idChannel <- utilx.NewSnowflake().GenerateID()
			}
		}()
	}
	wg.Wait()
	close(idChannel)

	duplicate := make(map[int64]struct{}, P*LOOP) // id排重
	for id := range idChannel {
		duplicate[id] = struct{}{}
	}

	// 判断是否生成了 P * LOOP 个不同 ID
	if len(duplicate) != P*LOOP {
		t.Errorf("共生成%d个ID", len(duplicate))
	}
}

func BenchmarkSnowflake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utilx.NewSnowflake().GenerateID()
	}
}

func BenchmarkSnowflakeByBwmarrin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		node, _ := snowflake.NewNode(int64(0))
		node.Generate()
	}
}

// go test -v ./utilx -run=^TestSnowflake$ -count=1
// go test ./utilx -bench=^BenchmarkSnow -run=^$ -benchmem -count=1
/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./utilx -bench=^BenchmarkSnow -run=^$ -benchmem -count=1
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/utilx
cpu: Apple M1 Pro
BenchmarkSnowflake-10                   21216313                53.28 ns/op            0 B/op          0 allocs/op
BenchmarkSnowflakeByBwmarrin-10         10079784               119.8 ns/op            96 B/op          1 allocs/op
PASS
ok      github.com/yzletter/go-toolery/utilx    2.901s
*/
