package utilx_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yzletter/go-toolery/utilx"
)

func TestMinimumConcurrencyBalancer(t *testing.T) {
	// 构造 Balancer
	endPoints := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5"}
	concurrency := make([]int32, len(endPoints))
	balancer := utilx.NewMinimumConcurrencyBalancer(endPoints, concurrency)

	useCount := make([]int32, len(endPoints)) // 统计每个 endPoints 的负载

	const P = 100

	wg := sync.WaitGroup{}
	wg.Add(P)

	for t := 1; t <= P; t++ {
		go func() {
			defer wg.Done()
			for i := 0; i < P; i++ {
				idx, _ := balancer.Take() // 取 EndPoint

				atomic.AddInt32(&useCount[idx], 1) // 记录负载
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

				_ = balancer.Return(idx) // 还 EndPoint

			}
		}()
	}

	wg.Wait()

	// 查看每个 EndPoint 的负载
	fmt.Println(useCount)
}

// go test -v ./utilx -run=^TestMinimumConcurrencyBalancer$ -count=1

/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test -v ./utilx -run=^TestMinimumConcurrencyBalancer$ -count=1

=== RUN   TestMinimumConcurrencyBalancer
[2007 1987 1991 2008 2007]
--- PASS: TestMinimumConcurrencyBalancer (2.89s)
PASS
ok      github.com/yzletter/go-toolery/utilx    3.447s
*/
