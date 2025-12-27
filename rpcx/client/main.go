package main

import (
	"log"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/yzletter/go-toolery/rpcx"
	"github.com/yzletter/go-toolery/rpcx/serializer"
)

func main() {

	s := serializer.MySerializer{}

	// 开 10 个协程并发调用
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)

	for i := 1; i <= P; i++ {
		go func() {
			defer wg.Done()
			client := NewClient("127.0.0.1", 5678, s)
			for j := 1; j <= P; j++ {
				req := &rpcx.RpcxData{
					A:  rand.Int(),
					B:  0,
					C:  false,
					D:  0,
					E:  "",
					Id: uuid.NewString(),
				}

				resp := client.Call(req)
				if resp == nil {
					log.Println("调用失败")
				} else if resp.Id != req.Id {
					log.Println("调用结果不符合预期")
				} else {
					log.Printf("协程 %d 的第 %d 次调用成功", i, j)
				}
			}
		}()
	}
	wg.Wait()
}
