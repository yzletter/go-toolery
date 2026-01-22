package concurrent_mapx_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/yzletter/go-toolery/datastructurex/concurrent_mapx"
)

func TestConcurrentMap(t *testing.T) {
	MyMap := concurrent_mapx.NewConcurrentMap(8, 1000)
	for i := 0; i < 10; i++ {
		MyMap.Set("key_"+strconv.Itoa(i), i)
	}
	for i := 0; i < 10; i++ {
		res, ok := MyMap.Get("key_" + strconv.Itoa(i))
		if !ok {
			t.Fail()
			return
		}

		fmt.Println(res)
	}

}

func TestConcurrentHashMapIterator(t *testing.T) {
	MyMap := concurrent_mapx.NewConcurrentMap(8, 1000)
	for i := 0; i < 10; i++ {
		MyMap.Set(strconv.Itoa(i), i)
	}
	iterator := MyMap.NewConcurrentMapIterator()
	entry := iterator.Next()
	for entry != nil {
		fmt.Println(entry.Key, entry.Value)
		entry = iterator.Next()
	}
}

// go test -v ./datastructurex/concurrent_mapx -run=^TestConcurrentMap$ -count=1
// go test -v ./datastructurex/concurrent_mapx -run=^TestConcurrentHashMapIterator$ -count=1

const P = 300
const T = 10000
const C = 20000

var myMap = concurrent_mapx.NewConcurrentMap(C, P*T)
var syncMap = sync.Map{}

func writeMyMap() {
	for i := 0; i < T; i++ {
		key := strconv.Itoa(rand.Int())
		myMap.Set(key, 1)
	}
}

func readMyMap() {
	for i := 0; i < T; i++ {
		key := strconv.Itoa(rand.Int())
		myMap.Get(key)
	}
}

func writeSyncMap() {
	for i := 0; i < T; i++ {
		key := strconv.Itoa(rand.Int())
		syncMap.Store(key, 1)
	}
}

func readSyncMap() {
	for i := 0; i < T; i++ {
		key := strconv.Itoa(rand.Int())
		syncMap.Load(key)
	}
}

func BenchmarkMyMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(P * 2)

		// P 个写程
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				writeMyMap()
				time.Sleep(100 * time.Millisecond)
			}()
		}

		// P 个读程
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				readMyMap()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSyncMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const P = 300
		wg := sync.WaitGroup{}
		wg.Add(P * 2)

		// P 个写程
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				writeSyncMap()
				time.Sleep(100 * time.Millisecond)
			}()
		}

		// P 个读程
		for i := 0; i < P; i++ {
			go func() {
				defer wg.Done()
				readSyncMap()
			}()
		}
		wg.Wait()
	}
}

// go test ./datastructurex/concurrent_mapx -run=^$ -bench=^Benchmark -benchtime=3s -count=1 -benchmem

/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./datastructurex/concurrent_mapx -run=^$ -bench=^Benchmark -benchtime=3s -count=1 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/datastructurex/concurrent_mapx
cpu: Apple M1 Pro
BenchmarkMyMap-10             10         479212142 ns/op        429347919 B/op   6019169 allocs/op
BenchmarkSyncMap-10            8         735929958 ns/op        523149185 B/op  13170668 allocs/op
PASS
ok      github.com/yzletter/go-toolery/datastructurex/concurrent_mapx   13.237s
*/
