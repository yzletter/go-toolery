package main

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/yzletter/go-toolery/rpcx"
	"github.com/yzletter/go-toolery/rpcx/serializer"
)

type EchoClient struct {
	conn          net.Conn
	s             serializer.Serializer
	requestBuffer sync.Map
}

func NewClient(serverIP string, port int, s serializer.Serializer) *EchoClient {
	conn, err := net.DialTimeout("udp", serverIP+":"+strconv.Itoa(port), time.Minute*30)
	if err != nil {
		log.Println(err)
		return nil
	}

	// 构造 client
	client := &EchoClient{
		conn:          conn,
		s:             s,
		requestBuffer: sync.Map{},
	}

	go client.receive() // 开启协程接收信息

	return client
}

func (client *EchoClient) receive() {

	for {
		// 读取数据
		resp := make([]byte, 4096)
		n, err := client.conn.Read(resp)
		if err != nil {
			log.Println(err)
		} else {
			// 反序列化
			var data rpcx.RpcxData
			err = client.s.Unmarshal(resp[:n], &data)
			if err != nil {
				log.Println(err)
				return
			}

			// 取出 map 里的 v 并断言成 channel
			if v, exist := client.requestBuffer.Load(data.Id); exist {
				// 删除 map 里的 key
				client.requestBuffer.Delete(data.Id)

				ch, ok := v.(chan rpcx.RpcxData)
				if !ok {
					log.Println(err)
					return
				}

				ch <- data // 将 data 写入 ch 解除 Call 的阻塞
			} else {
				log.Println("map 里没有 id")
			}
		}
	}
}

// Call 发起调用
func (client *EchoClient) Call(req *rpcx.RpcxData) *rpcx.RpcxData { // 需要传指针，因为要修改它的 map
	// 序列化
	stream, err := client.s.Marshal(*req)
	if err != nil {
		log.Println(err)
		return nil
	}

	// 发送请求
	_, err = client.conn.Write(stream)
	if err != nil {
		log.Println(err)
		return nil
	}

	// 将 <id, ch> 存入
	ch := make(chan rpcx.RpcxData, 1)
	client.requestBuffer.Store(req.Id, ch)

	rpcData := <-ch // 等待 receive 后往管道塞入 rpcData 解除阻塞
	return &rpcData
}
