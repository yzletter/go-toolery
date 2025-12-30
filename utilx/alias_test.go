package utilx_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/yzletter/go-toolery/utilx"
)

func TestAliasSample(t *testing.T) {
	endPoints := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5"}
	weight := []float64{1, 2, 3, 4, 5} // 预期 1 : 2 : 3 : 4 : 5

	sampler, _ := utilx.NewAliasSampler(weight)

	useCount := make([]int32, len(endPoints)) // 统计每个 endPoints 的负载

	const P = 100

	wg := sync.WaitGroup{}
	wg.Add(P)

	for t := 1; t <= P; t++ {
		go func() {
			defer wg.Done()
			for i := 0; i < P; i++ {
				idx := sampler.Sample()            // 采样
				atomic.AddInt32(&useCount[idx], 1) // 记录负载
			}
		}()
	}

	wg.Wait()

	// 查看每个 EndPoint 的负载
	fmt.Println(useCount)
}

// go test -v ./utilx -run=^TestAliasSample$ -count=1

/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test -v ./utilx -run=^TestAliasSample$ -count=1
=== RUN   TestAliasSample
[622 1343 2071 2661 3303]
--- PASS: TestAliasSample (0.00s)
PASS
ok      github.com/yzletter/go-toolery/utilx    0.539s
*/
